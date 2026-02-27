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
	"golang.org/x/sync/errgroup"
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

// Evmc is the main client for interacting with EVM-compatible blockchain nodes.
// It holds an underlying rpc.Client and exposes namespaces via getter methods
// such as [Evmc.Eth], [Evmc.Debug], and [Evmc.Kaia].
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
	// Clone the default transport to avoid mutating the global http.DefaultTransport.
	base := http.DefaultTransport.(*http.Transport).Clone()
	base.MaxIdleConns = o.connPool
	base.MaxIdleConnsPerHost = o.connPool
	base.MaxConnsPerHost = o.connPool
	base.IdleConnTimeout = o.idleConnTimeout
	base.DisableKeepAlives = false
	return &http.Client{Transport: base, Timeout: o.reqTimeout}
}

// New creates a new Evmc client connected to the given HTTP/HTTPS RPC endpoint.
func New(httpURL string, opts ...Options) (*Evmc, error) {
	return NewWithContext(context.Background(), httpURL, opts...)
}

// NewWithContext creates a new Evmc client connected to the given HTTP/HTTPS
// RPC endpoint using the provided context for the dial operation.
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

// NewWebsocket creates a new Evmc client connected to the given WS/WSS
// RPC endpoint for real-time subscriptions.
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

// Close shuts down the underlying RPC client connection.
func (e *Evmc) Close() {
	e.c.Close()
}

// IsWebsocket reports whether this client is connected via WebSocket.
func (e *Evmc) IsWebsocket() bool {
	return e.isWebsocket
}

// ChainID returns the chain ID of the connected network.
func (e *Evmc) ChainID() (uint64, error) {
	return e.eth.ChainID()
}

// NodeClient returns the name and version of the connected node client
// by parsing the web3_clientVersion response.
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

// Web3 returns the web3 namespace for utility RPC methods.
func (e *Evmc) Web3() *web3Namespace {
	return e.web3
}

// Eth returns the eth namespace for standard Ethereum RPC methods
// including blocks, transactions, receipts, and logs.
func (e *Evmc) Eth() *ethNamespace {
	return e.eth
}

// Kaia returns the kaia namespace for Kaia blockchain-specific RPC methods.
func (e *Evmc) Kaia() *kaiaNamespace {
	return e.kaia
}

// Debug returns the debug namespace for trace and debugging RPC methods.
func (e *Evmc) Debug() *debugNamespace {
	return e.debug
}

// Contract returns the contract namespace for raw smart contract calls.
func (e *Evmc) Contract() *contract {
	return e.contract
}

// BatchCallWithContext splits elements into chunks of maxBatchItems and
// sends them in parallel using up to workers goroutines.
// If workers is less than 1, it defaults to batchCallWorkers.
func (e *Evmc) BatchCallWithContext(ctx context.Context, elements []rpc.BatchElem, workers int) error {
	if workers < 1 {
		workers = e.batchCallWorkers
	}
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(workers)
	for i := 0; i < len(elements); i += e.maxBatchItems {
		chunk := elements[i:min(i+e.maxBatchItems, len(elements))]
		g.Go(func() error {
			return e.batchCall(ctx, chunk)
		})
	}
	return g.Wait()
}

// BatchCall is a batch call with workers.
// If workers is less than 1, it will be set to 1.
func (e *Evmc) BatchCall(elements []rpc.BatchElem, workers int) error {
	return e.BatchCallWithContext(context.Background(), elements, workers)
}

// ERC20 returns the ERC-20 token contract namespace for standard token operations.
func (e *Evmc) ERC20() *erc20Contract {
	return e.erc20
}

// ERC721 returns the ERC-721 NFT contract namespace (not fully implemented).
func (e *Evmc) ERC721() *erc721Contract {
	return e.erc721
}

// ERC1155 returns the ERC-1155 multi-token contract namespace (not fully implemented).
func (e *Evmc) ERC1155() *erc1155Contract {
	return e.erc1155
}

func (e *Evmc) call(
	ctx context.Context,
	result interface{},
	method Procedure,
	params ...interface{},
) error {
	return e.c.CallContext(ctx, result, method.String(), params...)
}

func (e *Evmc) batchCall(ctx context.Context, elements []rpc.BatchElem) error {
	return e.c.BatchCallContext(ctx, elements)
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
