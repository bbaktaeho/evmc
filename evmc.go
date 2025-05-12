package evmc

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
)

// TODO: websocket RPC
// TODO: backoff retry

type clientInfo interface {
	ChainID() (uint64, error)
	NodeClient() (name, version string, err error)
	IsWebsocket() bool
}

type caller interface {
	BatchCallWithContext(ctx context.Context, elements []rpc.BatchElem, workers int) error
	call(ctx context.Context, result interface{}, method Procedure, params ...interface{}) error
	batchCall(ctx context.Context, elements []rpc.BatchElem) error
}

type subscriber interface {
	subscribe(ctx context.Context, namespace string, ch interface{}, args ...interface{}) (evmctypes.Subscription, error)
}

type nodeSetter interface {
	setNode(clientVersion string)
}

type transactionSender interface {
	sendRawTransaction(ctx context.Context, rawTx string) (string, error)
}

type Evmc struct {
	c           *rpc.Client
	isWebsocket bool

	maxBatchItems    int
	batchCallWorkers int

	eth   *ethNamespace
	web3  *web3Namespace
	debug *debugNamespace
	// trace *traceNamespace
	// ots   *otsNamespace
	kaia *kaiaNamespace

	contract *contract
	erc20    *erc20Contract
	erc721   *erc721Contract
	erc1155  *erc1155Contract

	abiCache *lru.Cache[string, interface{}]
}

func httpClient(o *options) *http.Client {
	transport := http.DefaultTransport.(*http.Transport)
	transport.MaxIdleConns = o.connPool
	transport.MaxIdleConnsPerHost = o.connPool
	transport.MaxConnsPerHost = o.connPool
	transport.IdleConnTimeout = o.idleConnTimeout
	transport.DisableKeepAlives = false
	return &http.Client{Transport: transport, Timeout: o.reqTimeout}
}

func New(httpURL string, opts ...Options) (*Evmc, error) {
	return NewWithContext(context.Background(), httpURL, opts...)
}

func NewWithContext(ctx context.Context, httpURL string, opts ...Options) (*Evmc, error) {
	u, err := url.Parse(httpURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("invalid http scheme")
	}
	return newClient(ctx, httpURL, false, opts...)
}

func NewWebsocket(ctx context.Context, wsURL string, opts ...Options) (*Evmc, error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return nil, errors.New("invalid websocket scheme")
	}
	return newClient(ctx, wsURL, true, opts...)
}

func newClient(ctx context.Context, url string, isWs bool, opts ...Options) (*Evmc, error) {
	o := newOps()
	for _, opt := range opts {
		opt.apply(o)
	}

	rpcClient, err := rpc.DialOptions(
		ctx,
		url,
		rpc.WithHTTPClient(httpClient(o)),
		rpc.WithBatchItemLimit(o.maxBatchItems),
		rpc.WithBatchResponseSizeLimit(o.maxBatchSize),
		rpc.WithWebsocketDialer(websocket.Dialer{
			ReadBufferSize:  o.wsReadBufferSize,
			WriteBufferSize: o.wsWriteBufferSize,
		}),
		rpc.WithWebsocketMessageSizeLimit(int64(o.wsMessageSizeLimit)),
	)
	if err != nil {
		return nil, err
	}

	evmc := &Evmc{
		c:                rpcClient,
		isWebsocket:      isWs,
		abiCache:         lru.NewCache[string, interface{}](10),
		maxBatchItems:    o.maxBatchItems,
		batchCallWorkers: o.batchCallWorkers,
	}
	evmc.eth = &ethNamespace{info: evmc, c: evmc, s: evmc, ts: evmc}
	evmc.web3 = &web3Namespace{c: evmc}
	evmc.debug = &debugNamespace{c: evmc}
	evmc.kaia = &kaiaNamespace{c: evmc}
	evmc.contract = &contract{c: evmc}
	evmc.erc20 = &erc20Contract{info: evmc, c: evmc, ts: evmc}

	return evmc, nil
}

func (e *Evmc) Close() {
	e.c.Close()
}

func (e *Evmc) IsWebsocket() bool {
	return e.isWebsocket
}

func (e *Evmc) ChainID() (uint64, error) {
	return e.eth.ChainID()
}

func (e *Evmc) NodeClient() (name, version string, err error) {
	cv, err := e.web3.ClientVersion()
	if err != nil {
		return "", "", err
	}
	cvarr := strings.Split(cv, "/")
	if len(cvarr) < 2 {
		return
	}
	name = cvarr[0]
	version = cvarr[1]
	return
}

func (e *Evmc) Web3() *web3Namespace {
	return e.web3
}

func (e *Evmc) Eth() *ethNamespace {
	return e.eth
}

func (e *Evmc) Kaia() *kaiaNamespace {
	return e.kaia
}

func (e *Evmc) Debug() *debugNamespace {
	return e.debug
}

func (e *Evmc) Contract() *contract {
	return e.contract
}

// BatchCallWithContext is a batch call with context and workers.
// If workers is less than 1, it will be set to batchCallWorkers.
func (e *Evmc) BatchCallWithContext(ctx context.Context, elements []rpc.BatchElem, workers int) error {
	if workers < 1 {
		workers = e.batchCallWorkers
	}
	var (
		elementsCh = make(chan []rpc.BatchElem, len(elements)/e.maxBatchItems+1)
		finishCh   = make(chan struct{}, workers)
		errs       = make([]error, workers)
	)
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer func() {
				finishCh <- struct{}{}
			}()
			for es := range elementsCh {
				if err := e.batchCall(ctx, es); err != nil {
					errs[workerID] = err
					return
				}
			}
		}(i)
	}
	for i := 0; i < len(elements); i += e.maxBatchItems {
		j := min(i+e.maxBatchItems, len(elements))
		elementsCh <- elements[i:j]
	}
	close(elementsCh)
	for i := 0; i < workers; i++ {
		<-finishCh
	}
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchCall is a batch call with workers.
// If workers is less than 1, it will be set to 1.
func (e *Evmc) BatchCall(elements []rpc.BatchElem, workers int) error {
	return e.BatchCallWithContext(context.Background(), elements, workers)
}

// func (e *Evmc) Trace() *traceNamespace {
// 	return e.trace
// }

// func (e *Evmc) Ots() *otsNamespace {
// 	return e.ots
// }

func (e *Evmc) ERC20() *erc20Contract {
	return e.erc20
}

func (e *Evmc) ERC721() *erc721Contract {
	return e.erc721
}

func (e *Evmc) ERC1155() *erc1155Contract {
	return e.erc1155
}

// func (e *Evmc) Ots() {}

func (e *Evmc) call(
	ctx context.Context,
	result interface{},
	method Procedure,
	params ...interface{},
) error {
	if err := e.c.CallContext(ctx, result, method.String(), params...); err != nil {
		return err
	}
	return nil
}

func (e *Evmc) batchCall(ctx context.Context, elements []rpc.BatchElem) error {
	return e.c.BatchCallContext(ctx, elements)
	// var (
	// 	size      = len(elements)
	// 	loopCount = size / e.maxBatchItems
	// 	wg        = new(sync.WaitGroup)
	// 	errCh     = make(chan error)
	// 	errs      = make([]error, 0, 1)
	// )
	// if size%e.maxBatchItems != 0 {
	// 	loopCount++
	// }
	// go func() {
	// 	for e := range errCh {
	// 		errs = append(errs, e)
	// 	}
	// }()
	// for i := 0; i < loopCount; i++ {
	// 	low := e.maxBatchItems * i
	// 	high := e.maxBatchItems * (i + 1)
	// 	if high > size {
	// 		high = size
	// 	}
	// 	div := elements[low:high]

	// 	wg.Add(1)
	// 	go func(subElements []rpc.BatchElem) {
	// 		defer wg.Done()

	// 		if err := e.c.BatchCallContext(ctx, subElements); err != nil {
	// 			errCh <- err
	// 		}
	// 	}(div)
	// }
	// wg.Wait()
	// close(errCh)

	// if len(errs) > 0 {
	// 	return errs[0]
	// }
	// return nil
}

func (e *Evmc) subscribe(
	ctx context.Context,
	namespace string,
	ch interface{},
	args ...interface{},
) (evmctypes.Subscription, error) {
	subscription, err := e.c.Subscribe(ctx, namespace, ch, args...)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (e *Evmc) sendRawTransaction(ctx context.Context, rawTx string) (string, error) {
	result := new(string)
	if err := e.call(ctx, result, EthSendRawTransaction, rawTx); err != nil {
		return "", err
	}
	return *result, nil
}
