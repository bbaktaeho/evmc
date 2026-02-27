package evmc

import (
	"context"
	"encoding/json"
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

type PrestateTracerConfig struct {
	// If true, this tracer will return state modifications
	DiffMode bool `json:"diffMode"`
	// If true, this tracer will not return the contract code
	DisableCode bool `json:"disableCode"`
	// If true, this tracer will not return the contract storage
	DisableStorage bool `json:"disableStorage"`
	// If true, this tracer will return empty state objects
	IncludeEmpty bool `json:"includeEmpty"`
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

// newTracerConfig builds a TraceConfig for a specific tracer with optional timeout, reexec, and tracer-specific config.
func newTracerConfig(tracer Tracer, timeout time.Duration, reexec *uint64, cfg interface{}) *TraceConfig {
	tc := &TraceConfig{Tracer: tracer, Timeout: timeout.String(), Reexec: reexec}
	if cfg != nil {
		tc.TracerConfig = cfg
	}
	return tc
}

// newCustomTracerConfig builds a TraceConfig for a custom JavaScript tracer with optional timeout and reexec.
func newCustomTracerConfig(jsTracer string, timeout time.Duration, reexec *uint64) *TraceConfig {
	return &TraceConfig{
		Tracer:  Tracer(jsTracer),
		Timeout: timeout.String(),
		Reexec:  reexec,
	}
}

// ─── TraceBlockByNumber ──────────────────────────────────────────────────────

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
	var callTracers = []*evmctypes.CallTracer{}
	traceCfg := newTracerConfig(CallTracer, timeout, reexec, cfg)
	if err := d.traceBlockByNumber(ctx, blockNumber, traceCfg, &callTracers); err != nil {
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
	var flatCallTracers = []*evmctypes.FlatCallTracer{}
	traceCfg := newTracerConfig(FlatCallTracer, timeout, reexec, cfg)
	if err := d.traceBlockByNumber(ctx, blockNumber, traceCfg, &flatCallTracers); err != nil {
		return nil, err
	}
	assignIndexFlatCalls(flatCallTracers)
	return flatCallTracers, nil
}

func (d *debugNamespace) TraceBlockByNumberWithContext_prestateTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	[]*evmctypes.PrestateTracer,
	error,
) {
	return d.traceBlockByNumber_prestateTracer(ctx, blockNumber, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceBlockByNumber_prestateTracer(
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	[]*evmctypes.PrestateTracer,
	error,
) {
	return d.traceBlockByNumber_prestateTracer(context.Background(), blockNumber, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByNumber_prestateTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	[]*evmctypes.PrestateTracer,
	error,
) {
	var prestateTracers = []*evmctypes.PrestateTracer{}
	traceCfg := newTracerConfig(PrestateTracer, timeout, reexec, cfg)
	if err := d.traceBlockByNumber(ctx, blockNumber, traceCfg, &prestateTracers); err != nil {
		return nil, err
	}
	return prestateTracers, nil
}

func (d *debugNamespace) TraceBlockByNumber_customTracer(
	blockNumber uint64,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (
	[]*evmctypes.CustomTraceResult,
	error,
) {
	return d.TraceBlockByNumberWithContext_customTracer(context.Background(), blockNumber, jsTracer, timeout, reexec)
}

func (d *debugNamespace) TraceBlockByNumberWithContext_customTracer(
	ctx context.Context,
	blockNumber uint64,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (
	[]*evmctypes.CustomTraceResult,
	error,
) {
	var result = []*evmctypes.CustomTraceResult{}
	traceCfg := newCustomTracerConfig(jsTracer, timeout, reexec)
	if err := d.traceBlockByNumber(ctx, blockNumber, traceCfg, &result); err != nil {
		return nil, err
	}
	return result, nil
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
	return d.c.call(ctx, result, DebugTraceBlockByNumber, params...)
}

// ─── TraceBlockByHash ────────────────────────────────────────────────────────

func (d *debugNamespace) TraceBlockByHashWithContext(
	ctx context.Context,
	hash string,
	cfg *TraceConfig,
) (
	[]*evmctypes.TraceResult,
	error,
) {
	var result = []*evmctypes.TraceResult{}
	if err := d.traceBlockByHash(ctx, hash, cfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *debugNamespace) TraceBlockByHash(hash string, cfg *TraceConfig) ([]*evmctypes.TraceResult, error) {
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
	var callTracers = []*evmctypes.CallTracer{}
	traceCfg := newTracerConfig(CallTracer, timeout, reexec, cfg)
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
	var flatCallTracers = []*evmctypes.FlatCallTracer{}
	traceCfg := newTracerConfig(FlatCallTracer, timeout, reexec, cfg)
	if err := d.traceBlockByHash(ctx, hash, traceCfg, &flatCallTracers); err != nil {
		return nil, err
	}
	assignIndexFlatCalls(flatCallTracers)
	return flatCallTracers, nil
}

func (d *debugNamespace) TraceBlockByHashWithContext_prestateTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	[]*evmctypes.PrestateTracer,
	error,
) {
	return d.traceBlockByHash_prestateTracer(ctx, hash, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceBlockByHash_prestateTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	[]*evmctypes.PrestateTracer,
	error,
) {
	return d.traceBlockByHash_prestateTracer(context.Background(), hash, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByHash_prestateTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	[]*evmctypes.PrestateTracer,
	error,
) {
	var prestateTracers = []*evmctypes.PrestateTracer{}
	traceCfg := newTracerConfig(PrestateTracer, timeout, reexec, cfg)
	if err := d.traceBlockByHash(ctx, hash, traceCfg, &prestateTracers); err != nil {
		return nil, err
	}
	return prestateTracers, nil
}

func (d *debugNamespace) TraceBlockByHash_customTracer(
	hash string,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (
	[]*evmctypes.CustomTraceResult,
	error,
) {
	return d.TraceBlockByHashWithContext_customTracer(context.Background(), hash, jsTracer, timeout, reexec)
}

func (d *debugNamespace) TraceBlockByHashWithContext_customTracer(
	ctx context.Context,
	hash string,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (
	[]*evmctypes.CustomTraceResult,
	error,
) {
	var result = []*evmctypes.CustomTraceResult{}
	traceCfg := newCustomTracerConfig(jsTracer, timeout, reexec)
	if err := d.traceBlockByHash(ctx, hash, traceCfg, &result); err != nil {
		return nil, err
	}
	return result, nil
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
	return d.c.call(ctx, result, DebugTraceBlockByHash, params...)
}

// ─── TraceTransaction ────────────────────────────────────────────────────────

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
	callFrame := &evmctypes.CallFrame{}
	traceCfg := newTracerConfig(CallTracer, timeout, reexec, cfg)
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
	[]*evmctypes.FlatCallFrame,
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
	[]*evmctypes.FlatCallFrame,
	error,
) {
	var flatCallFrames = []*evmctypes.FlatCallFrame{}
	traceCfg := newTracerConfig(FlatCallTracer, timeout, reexec, cfg)
	if err := d.traceTransaction(ctx, hash, traceCfg, &flatCallFrames); err != nil {
		return nil, err
	}
	return flatCallFrames, nil
}

func (d *debugNamespace) TraceTransaction_prestateTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	*evmctypes.PrestateResult,
	error,
) {
	return d.TraceTransactionWithContext_prestateTracer(context.Background(), hash, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceTransactionWithContext_prestateTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	*evmctypes.PrestateResult,
	error,
) {
	prestateResult := &evmctypes.PrestateResult{}
	traceCfg := newTracerConfig(PrestateTracer, timeout, reexec, cfg)
	if err := d.traceTransaction(ctx, hash, traceCfg, prestateResult); err != nil {
		return nil, err
	}
	return prestateResult, nil
}

func (d *debugNamespace) TraceTransaction_customTracer(
	hash string,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (
	json.RawMessage,
	error,
) {
	return d.TraceTransactionWithContext_customTracer(context.Background(), hash, jsTracer, timeout, reexec)
}

func (d *debugNamespace) TraceTransactionWithContext_customTracer(
	ctx context.Context,
	hash string,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (
	json.RawMessage,
	error,
) {
	var result json.RawMessage
	traceCfg := newCustomTracerConfig(jsTracer, timeout, reexec)
	if err := d.traceTransaction(ctx, hash, traceCfg, &result); err != nil {
		return nil, err
	}
	return result, nil
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
	return d.c.call(ctx, result, DebugTraceTransaction, params...)
}

// ─── TraceCall ───────────────────────────────────────────────────────────────

func (d *debugNamespace) TraceCall(tx *Tx, blockAndTag evmctypes.BlockAndTag, cfg *TraceConfig) (interface{}, error) {
	return d.TraceCallWithContext(context.Background(), tx, blockAndTag, cfg)
}

func (d *debugNamespace) TraceCallWithContext(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	cfg *TraceConfig,
) (
	interface{},
	error,
) {
	var result interface{}
	if err := d.traceCall(ctx, tx, blockAndTag, cfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *debugNamespace) TraceCall_callTracer(
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	*evmctypes.CallFrame,
	error,
) {
	return d.TraceCallWithContext_callTracer(context.Background(), tx, blockAndTag, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceCallWithContext_callTracer(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (
	*evmctypes.CallFrame,
	error,
) {
	callFrame := &evmctypes.CallFrame{}
	traceCfg := newTracerConfig(CallTracer, timeout, reexec, cfg)
	if err := d.traceCall(ctx, tx, blockAndTag, traceCfg, callFrame); err != nil {
		return nil, err
	}
	assignIndexCalls(callFrame)
	return callFrame, nil
}

func (d *debugNamespace) TraceCall_flatCallTracer(
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	[]*evmctypes.FlatCallFrame,
	error,
) {
	return d.TraceCallWithContext_flatCallTracer(context.Background(), tx, blockAndTag, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceCallWithContext_flatCallTracer(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) (
	[]*evmctypes.FlatCallFrame,
	error,
) {
	var flatCallFrames = []*evmctypes.FlatCallFrame{}
	traceCfg := newTracerConfig(FlatCallTracer, timeout, reexec, cfg)
	if err := d.traceCall(ctx, tx, blockAndTag, traceCfg, &flatCallFrames); err != nil {
		return nil, err
	}
	return flatCallFrames, nil
}

func (d *debugNamespace) TraceCall_prestateTracer(
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	*evmctypes.PrestateResult,
	error,
) {
	return d.TraceCallWithContext_prestateTracer(context.Background(), tx, blockAndTag, timeout, reexec, cfg)
}

func (d *debugNamespace) TraceCallWithContext_prestateTracer(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (
	*evmctypes.PrestateResult,
	error,
) {
	prestateResult := &evmctypes.PrestateResult{}
	traceCfg := newTracerConfig(PrestateTracer, timeout, reexec, cfg)
	if err := d.traceCall(ctx, tx, blockAndTag, traceCfg, prestateResult); err != nil {
		return nil, err
	}
	return prestateResult, nil
}

func (d *debugNamespace) TraceCall_customTracer(
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (
	json.RawMessage,
	error,
) {
	return d.TraceCallWithContext_customTracer(context.Background(), tx, blockAndTag, jsTracer, timeout, reexec)
}

func (d *debugNamespace) TraceCallWithContext_customTracer(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (
	json.RawMessage,
	error,
) {
	var result json.RawMessage
	traceCfg := newCustomTracerConfig(jsTracer, timeout, reexec)
	if err := d.traceCall(ctx, tx, blockAndTag, traceCfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *debugNamespace) traceCall(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	traceCfg *TraceConfig,
	result interface{},
) error {
	msg, err := tx.parseCallMsg()
	if err != nil {
		return err
	}
	params := []interface{}{msg, blockAndTag.String()}
	if traceCfg != nil {
		params = append(params, *traceCfg)
	}
	return d.c.call(ctx, result, DebugTraceCall, params...)
}

// ─── Raw / Getter Methods ───────────────────────────────────────────────────

func (d *debugNamespace) GetRawHeader(blockAndTag evmctypes.BlockAndTag) (string, error) {
	return d.GetRawHeaderWithContext(context.Background(), blockAndTag)
}

func (d *debugNamespace) GetRawHeaderWithContext(ctx context.Context, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return d.getRawHeader(ctx, blockAndTag)
}

func (d *debugNamespace) getRawHeader(ctx context.Context, blockAndTag evmctypes.BlockAndTag) (string, error) {
	result := new(string)
	if err := d.c.call(ctx, result, DebugGetRawHeader, blockAndTag.String()); err != nil {
		return "", err
	}
	return *result, nil
}

func (d *debugNamespace) GetRawBlock(blockAndTag evmctypes.BlockAndTag) (string, error) {
	return d.GetRawBlockWithContext(context.Background(), blockAndTag)
}

func (d *debugNamespace) GetRawBlockWithContext(ctx context.Context, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return d.getRawBlock(ctx, blockAndTag)
}

func (d *debugNamespace) getRawBlock(ctx context.Context, blockAndTag evmctypes.BlockAndTag) (string, error) {
	result := new(string)
	if err := d.c.call(ctx, result, DebugGetRawBlock, blockAndTag.String()); err != nil {
		return "", err
	}
	return *result, nil
}

func (d *debugNamespace) GetRawTransaction(hash string) (string, error) {
	return d.GetRawTransactionWithContext(context.Background(), hash)
}

func (d *debugNamespace) GetRawTransactionWithContext(ctx context.Context, hash string) (string, error) {
	return d.getRawTransaction(ctx, hash)
}

func (d *debugNamespace) getRawTransaction(ctx context.Context, hash string) (string, error) {
	result := new(string)
	if err := d.c.call(ctx, result, DebugGetRawTransaction, hash); err != nil {
		return "", err
	}
	return *result, nil
}

func (d *debugNamespace) GetRawReceipts(blockAndTag evmctypes.BlockAndTag) ([]string, error) {
	return d.GetRawReceiptsWithContext(context.Background(), blockAndTag)
}

func (d *debugNamespace) GetRawReceiptsWithContext(ctx context.Context, blockAndTag evmctypes.BlockAndTag) ([]string, error) {
	return d.getRawReceipts(ctx, blockAndTag)
}

func (d *debugNamespace) getRawReceipts(ctx context.Context, blockAndTag evmctypes.BlockAndTag) ([]string, error) {
	result := new([]string)
	if err := d.c.call(ctx, result, DebugGetRawReceipts, blockAndTag.String()); err != nil {
		return nil, err
	}
	return *result, nil
}

func (d *debugNamespace) GetBadBlocks() ([]*evmctypes.BadBlock, error) {
	return d.GetBadBlocksWithContext(context.Background())
}

func (d *debugNamespace) GetBadBlocksWithContext(ctx context.Context) ([]*evmctypes.BadBlock, error) {
	return d.getBadBlocks(ctx)
}

func (d *debugNamespace) getBadBlocks(ctx context.Context) ([]*evmctypes.BadBlock, error) {
	result := new([]*evmctypes.BadBlock)
	if err := d.c.call(ctx, result, DebugGetBadBlocks); err != nil {
		return nil, err
	}
	return *result, nil
}
