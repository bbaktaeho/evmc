package evmc

// TODO
// - debug_traceBlock
// - debug_traceBlockFromFile
// - debug_traceBadBlock
// - debug_traceCall

import (
	"context"
	"time"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type debugNamespace struct {
	c caller
}

type Tracer string

// native tracer written in Go
const (
	CallTracer Tracer = "callTracer"
	// only go-ethereum client has this tracer
	FlatCallTracer Tracer = "flatCallTracer"
	FourByteTracer Tracer = "4byteTracer"
	MuxTracer      Tracer = "muxTracer"
	PrestateTracer Tracer = "prestateTracer"
)

type DefaultTraceConfig struct {
	EnableMemory     bool `json:"enableMemory"`
	DisableStack     bool `json:"disableStack"`
	DisableStorage   bool `json:"disableStorage"`
	EnableReturnData bool `json:"enableReturnData"`
	Debug            bool `json:"debug"`
	Limit            int  `json:"limit"`
}

type CallTracerConfig struct {
	// If true, call tracer won't collect any subcalls
	OnlyTopCall bool `json:"onlyTopCall"`
	// If true, call tracer will collect event logs
	WithLog bool `json:"withLog"`
}

type FlatCallTracerConfig struct {
	// If true, call tracer converts errors to parity format
	ConvertParityErrors bool `json:"convertParityErrors"`
	// If true, call tracer includes calls to precompiled contracts
	IncludePrecompiles bool `json:"includePrecompiles"`
}

type TraceConfig struct {
	*DefaultTraceConfig

	Tracer Tracer `json:"tracer,omitempty"`
	// Timeout is the maximum time the tracer is allowed to run
	Timeout string `json:"timeout,omitempty"`
	// Reexec is the number of blocks the tracer is willing to go back
	// and reexecute to produce missing historical state necessary to run a specific
	// trace.
	Reexec *uint64 `json:"reexec,omitempty"`

	// TraceConfig is a generic field that can be used to pass tracer-specific
	// configuration. The actual type of the field depends on the value of the Tracer field.
	//
	// For Example, if Tracer is CallTracer, then TracerConfig should be of type CallTracerConfig.
	// If Tracer is FlatCallTracer, then TracerConfig should be of type FlatCallTracerConfig.
	TracerConfig interface{} `json:"tracerConfig,omitempty"`
}

func assignIndexCalls(callFrame *evmctypes.CallFrame) {
	if callFrame == nil {
		return
	}
	var index uint64
	callFrame.Index = index // default 0
	assignIndexCall(&index, callFrame.Calls)
}

func assignIndexCall(index *uint64, calls []*evmctypes.CallFrame) {
	for _, call := range calls {
		*index++
		call.Index = *index
		if len(call.Calls) > 0 {
			assignIndexCall(index, call.Calls)
		}
	}
}

func assignIndexFlatCalls(flatCalls []*evmctypes.FlatCallTracer) {
	for _, flatCall := range flatCalls {
		if len(flatCall.Result) == 0 {
			continue
		}
		for j, call := range flatCall.Result {
			call.Index = uint64(j)
		}
	}
}

func (d *debugNamespace) TraceBlockByNumberWithContext(
	ctx context.Context,
	blockNumber uint64,
	cfg *TraceConfig,
) (
	[]*evmctypes.TraceResult,
	error,
) {
	var result = []*evmctypes.TraceResult{}
	if err := d.traceBlockByNumber(ctx, blockNumber, cfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *debugNamespace) TraceBlockByNumber(blockNumber uint64, cfg *TraceConfig) ([]*evmctypes.TraceResult, error) {
	return d.TraceBlockByNumberWithContext(context.Background(), blockNumber, cfg)
}

func (d *debugNamespace) TraceBlockByNumberWithContext_callTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	[]*evmctypes.CallTracer,
	error,
) {
	return d.traceBlockByNumber_callTracer(ctx, blockNumber, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceBlockByNumber_callTracer(
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	[]*evmctypes.CallTracer,
	error,
) {
	return d.traceBlockByNumber_callTracer(context.Background(), blockNumber, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByNumber_callTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	[]*evmctypes.CallTracer,
	error,
) {
	var (
		callTracers = []*evmctypes.CallTracer{}
		traceCfg    = TraceConfig{Tracer: CallTracer, Timeout: timeout.String(), Reexec: reexec}
	)
	if cfg != nil {
		traceCfg.TracerConfig = *cfg
	}
	if err := d.traceBlockByNumber(ctx, blockNumber, &traceCfg, &callTracers); err != nil {
		return nil, err
	}
	for _, callTracer := range callTracers {
		assignIndexCalls(callTracer.Result)
	}
	return callTracers, nil
}

// only go-ethereum client has this method
func (d *debugNamespace) TraceBlockByNumber_flatCallTracer(
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	[]*evmctypes.FlatCallTracer,
	error,
) {
	return d.TraceBlockByNumberWithContext_flatCallTracer(context.Background(), blockNumber, timeout, reexec, cfg)
}

// only go-ethereum client has this method
func (d *debugNamespace) TraceBlockByNumberWithContext_flatCallTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	[]*evmctypes.FlatCallTracer,
	error,
) {
	return d.traceBlockByNumber_flatCallTracer(ctx, blockNumber, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByNumber_flatCallTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	[]*evmctypes.FlatCallTracer,
	error,
) {
	var (
		flatCallTracers = []*evmctypes.FlatCallTracer{}
		traceCfg        = TraceConfig{Tracer: FlatCallTracer, Timeout: timeout.String(), Reexec: reexec}
	)
	if cfg != nil {
		traceCfg.TracerConfig = *cfg
	}
	if err := d.traceBlockByNumber(ctx, blockNumber, &traceCfg, &flatCallTracers); err != nil {
		return nil, err
	}
	assignIndexFlatCalls(flatCallTracers)
	return flatCallTracers, nil
}

func (d *debugNamespace) traceBlockByNumber(
	ctx context.Context,
	blockNumber uint64,
	traceCfg *TraceConfig,
	result interface{},
) error {
	params := []interface{}{hexutil.EncodeUint64(blockNumber)}
	if traceCfg != nil {
		params = append(params, *traceCfg)
	}
	if err := d.c.call(ctx, result, DebugTraceBlockByNumber, params...); err != nil {
		return err
	}
	return nil
}

func (d *debugNamespace) TraceBlockByHashWithContext(
	ctx context.Context,
	hash string,
	cfg *TraceConfig,
) (
	*evmctypes.TraceResult,
	error,
) {
	var result = &evmctypes.TraceResult{}
	if err := d.traceBlockByHash(ctx, hash, cfg, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *debugNamespace) TraceBlockByHash(hash string, cfg *TraceConfig) (*evmctypes.TraceResult, error) {
	return d.TraceBlockByHashWithContext(context.Background(), hash, cfg)
}

func (d *debugNamespace) TraceBlockByHashWithContext_callTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	[]*evmctypes.CallTracer,
	error,
) {
	return d.traceBlockByHash_callTracer(ctx, hash, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceBlockByHash_callTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	[]*evmctypes.CallTracer,
	error,
) {
	return d.traceBlockByHash_callTracer(context.Background(), hash, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByHash_callTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	[]*evmctypes.CallTracer,
	error,
) {
	var (
		callTracers = []*evmctypes.CallTracer{}
		traceCfg    = &TraceConfig{Tracer: CallTracer, Timeout: timeout.String(), Reexec: reexec}
	)
	if cfg != nil {
		traceCfg.TracerConfig = *cfg
	}
	if err := d.traceBlockByHash(ctx, hash, traceCfg, &callTracers); err != nil {
		return nil, err
	}
	for _, callTracer := range callTracers {
		assignIndexCalls(callTracer.Result)
	}
	return callTracers, nil
}

func (d *debugNamespace) TraceBlockByHash_flatCallTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	[]*evmctypes.FlatCallTracer,
	error,
) {
	return d.TraceBlockByHashWithContext_flatCallTracer(context.Background(), hash, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceBlockByHashWithContext_flatCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	[]*evmctypes.FlatCallTracer,
	error,
) {
	return d.traceBlockByHash_flatCallTracer(ctx, hash, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByHash_flatCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	[]*evmctypes.FlatCallTracer,
	error,
) {
	var (
		flatCallTracers = []*evmctypes.FlatCallTracer{}
		traceCfg        = TraceConfig{Tracer: FlatCallTracer, Timeout: timeout.String(), Reexec: reexec}
	)
	if cfg != nil {
		traceCfg.TracerConfig = *cfg
	}
	if err := d.traceBlockByHash(ctx, hash, &traceCfg, &flatCallTracers); err != nil {
		return nil, err
	}
	assignIndexFlatCalls(flatCallTracers)
	return flatCallTracers, nil
}

func (d *debugNamespace) traceBlockByHash(
	ctx context.Context,
	hash string,
	traceCfg *TraceConfig,
	result interface{},
) error {
	params := []interface{}{hash}
	if traceCfg != nil {
		params = append(params, *traceCfg)
	}
	if err := d.c.call(ctx, result, DebugTraceBlockByHash, params...); err != nil {
		return err
	}
	return nil
}

func (d *debugNamespace) TraceTransaction(hash string, cfg *TraceConfig) (interface{}, error) {
	return d.TraceTransactionWithContext(context.Background(), hash, cfg)
}

func (d *debugNamespace) TraceTransactionWithContext(
	ctx context.Context,
	hash string,
	cfg *TraceConfig,
) (
	interface{},
	error,
) {
	var result interface{}
	if err := d.traceTransaction(ctx, hash, cfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *debugNamespace) TraceTransaction_callTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	*evmctypes.CallFrame,
	error,
) {
	return d.TraceTransactionWithContext_callTracer(context.Background(), hash, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceTransactionWithContext_callTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	*evmctypes.CallFrame,
	error,
) {
	var (
		callFrame = &evmctypes.CallFrame{}
		traceCfg  = &TraceConfig{Tracer: CallTracer, Timeout: timeout.String(), Reexec: reexec}
	)
	if cfg != nil {
		traceCfg.TracerConfig = *cfg
	}
	if err := d.traceTransaction(ctx, hash, traceCfg, callFrame); err != nil {
		return nil, err
	}
	assignIndexCalls(callFrame)
	return callFrame, nil
}

func (d *debugNamespace) TraceTransaction_flatCallTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	*evmctypes.FlatCallTracer,
	error,
) {
	return d.TraceTransactionWithContext_flatCallTracer(context.Background(), hash, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceTransactionWithContext_flatCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	*evmctypes.FlatCallTracer,
	error,
) {
	var (
		flatCallTracer = &evmctypes.FlatCallTracer{}
		traceCfg       = TraceConfig{Tracer: FlatCallTracer, Timeout: timeout.String(), Reexec: reexec}
	)
	if cfg != nil {
		traceCfg.TracerConfig = *cfg
	}
	if err := d.traceTransaction(ctx, hash, &traceCfg, flatCallTracer); err != nil {
		return nil, err
	}
	assignIndexFlatCalls([]*evmctypes.FlatCallTracer{flatCallTracer})
	return flatCallTracer, nil
}

func (d *debugNamespace) traceTransaction(
	ctx context.Context,
	hash string,
	traceCfg *TraceConfig,
	result interface{},
) error {
	params := []interface{}{hash}
	if traceCfg != nil {
		params = append(params, *traceCfg)
	}
	if err := d.c.call(ctx, result, DebugTraceTransaction, params...); err != nil {
		return err
	}
	return nil
}
