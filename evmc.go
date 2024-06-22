package evmc

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
)

type caller interface {
	call(ctx context.Context, result interface{}, procedure string, params ...interface{}) error
	// batchCall(ctx context.Context, elements []rpc.BatchElem) error
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
}

// TODO: http client settings
func New(url string) (*Evmc, error) {
	ctx := context.Background()
	// rpc.WithHTTPClient(&http.Client{Transport: transport})
	rpcClient, err := rpc.DialOptions(ctx, url)
	if err != nil {
		return nil, err
	}

	evmc := &Evmc{c: rpcClient}
	evmc.eth = &ethNamespace{c: evmc}
	evmc.web3 = &web3Namespace{c: evmc}
	evmc.debug = &debugNamespace{c: evmc}

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
	evmc.nodeName = (cvarr[0])
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

// func (e *Evmc) Ots() {}

func (e *Evmc) call(
	ctx context.Context,
	result interface{},
	procedure string,
	params ...interface{},
) error {
	if err := e.c.CallContext(ctx, result, procedure, params...); err != nil {
		return err
	}
	return nil
}
