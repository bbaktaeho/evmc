package evmc

import (
	"context"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
)

type caller interface {
	call(ctx context.Context, result interface{}, method procedure, params ...interface{}) error
	// batchCall(ctx context.Context, elements []rpc.BatchElem) error
}

type contractCaller interface {
	contractCall(ctx context.Context, result interface{}, contract string, data string, parsedNumOrTag string) error
}

type Evmc struct {
	c *rpc.Client

	chainID     uint64
	nodeName    string
	nodeVersion string

	eth   *ethNamespace
	web3  *web3Namespace
	debug *debugNamespace
	// trace *traceNamespace
	// ots   *otsNamespace

	erc20 *erc20Contract
}

func httpTransport(o *options) *http.Transport {
	transport := http.DefaultTransport.(*http.Transport)
	transport.MaxIdleConns = o.connPool
	transport.MaxIdleConnsPerHost = o.connPool
	transport.MaxConnsPerHost = o.connPool
	transport.DisableKeepAlives = false
	return transport
}

func New(url string, opts ...Options) (*Evmc, error) {
	ctx := context.Background()
	o := newOps()
	for _, opt := range opts {
		opt.apply(o)
	}

	rpcClient, err := rpc.DialOptions(
		ctx,
		url,
		rpc.WithHTTPClient(&http.Client{Transport: httpTransport(o)}),
		rpc.WithBatchItemLimit(o.maxBatchItems),
		rpc.WithBatchResponseSizeLimit(o.maxBatchSize),
	)
	if err != nil {
		return nil, err
	}

	evmc := &Evmc{c: rpcClient}
	evmc.eth = &ethNamespace{c: evmc}
	evmc.web3 = &web3Namespace{c: evmc}
	evmc.debug = &debugNamespace{c: evmc}
	evmc.erc20 = &erc20Contract{c: evmc}

	chainID, err := evmc.eth.GetChainID()
	if err != nil {
		return nil, err
	}
	evmc.chainID = chainID

	cv, err := evmc.web3.ClientVersion()
	if err != nil {
		return nil, err
	}
	cvarr := strings.Split(cv, "/")
	evmc.nodeName = cvarr[0]
	evmc.nodeVersion = cvarr[1]

	return evmc, nil
}

func (e *Evmc) ChainID() uint64 {
	return e.chainID
}

func (e *Evmc) NodeClient() (name, version string) {
	return e.nodeName, e.nodeVersion
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

// func (e *Evmc) Ots() {}

func (e *Evmc) contractCall(
	ctx context.Context,
	result interface{},
	contract string,
	data string,
	parsedNumOrTag string,
) error {
	params := []interface{}{
		ContractCallParams{
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
