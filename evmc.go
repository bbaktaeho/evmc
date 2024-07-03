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
	ChainID() uint64
	NodeClient() (name, version string)
	IsWebsocket() bool
}

type caller interface {
	call(ctx context.Context, result interface{}, method procedure, params ...interface{}) error
	// batchCall(ctx context.Context, elements []rpc.BatchElem) error
}

type contractCaller interface {
	contractCall(ctx context.Context, result interface{}, contract, data, parsedNumOrTag string) error
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

	chainID     uint64
	nodeName    ClientName
	nodeVersion string

	eth   *ethNamespace
	web3  *web3Namespace
	debug *debugNamespace
	// trace *traceNamespace
	// ots   *otsNamespace

	erc20   *erc20Contract
	erc721  *erc721Contract
	erc1155 *erc1155Contract

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
	u, err := url.Parse(httpURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("invalid http scheme")
	}
	return newClient(context.Background(), httpURL, false, opts...)
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

	evmc := &Evmc{c: rpcClient, isWebsocket: isWs, abiCache: lru.NewCache[string, interface{}](10)}
	evmc.eth = &ethNamespace{info: evmc, c: evmc, s: evmc, ts: evmc}
	evmc.web3 = &web3Namespace{c: evmc, n: evmc}
	evmc.debug = &debugNamespace{c: evmc}
	evmc.erc20 = &erc20Contract{info: evmc, c: evmc, ts: evmc}

	chainID, err := evmc.eth.ChainID()
	if err != nil {
		return nil, err
	}
	evmc.chainID = chainID

	cv, err := evmc.web3.ClientVersion()
	if err != nil {
		return nil, err
	}
	cvarr := strings.Split(cv, "/")
	evmc.nodeName = ClientName(cvarr[0])
	evmc.nodeVersion = cvarr[1]

	return evmc, nil
}

func (e *Evmc) Close() {
	e.c.Close()
}

func (e *Evmc) IsWebsocket() bool {
	return e.isWebsocket
}

func (e *Evmc) ChainID() uint64 {
	return e.chainID
}

func (e *Evmc) NodeClient() (name, version string) {
	return e.nodeName.String(), e.nodeVersion
}

func (e *Evmc) Web3() *web3Namespace {
	return e.web3
}

func (e *Evmc) Eth() *ethNamespace {
	return e.eth
}
func (e *Evmc) Debug() *debugNamespace {
	return e.debug
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

func (e *Evmc) contractCall(
	ctx context.Context,
	result interface{},
	contract string,
	data string,
	parsedNumOrTag string,
) error {
	params := []interface{}{
		evmctypes.ContractCallParams{
			To:   contract,
			Data: data,
		},
		parsedNumOrTag,
	}
	if err := e.call(ctx, result, ethCall, params...); err != nil {
		return err
	}
	return nil
}

func (e *Evmc) call(
	ctx context.Context,
	result interface{},
	method procedure,
	params ...interface{},
) error {
	if err := e.c.CallContext(ctx, result, method.String(), params...); err != nil {
		return err
	}
	return nil
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
	if err := e.call(ctx, result, ethSendRawTransaction, rawTx); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *Evmc) setNode(cv string) {
	cvarr := strings.Split(cv, "/")
	e.nodeName = ClientName(cvarr[0])
	e.nodeVersion = cvarr[1]
}
