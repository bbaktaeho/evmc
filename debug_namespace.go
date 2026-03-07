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

// Tracer identifies a named tracer built into go-ethereum.
type Tracer string

// Native Go tracers supported by go-ethereum.
const (
	CallTracer     Tracer = "callTracer"
	FlatCallTracer Tracer = "flatCallTracer" // go-ethereum only
	FourByteTracer Tracer = "4byteTracer"
	MuxTracer      Tracer = "muxTracer"
	PrestateTracer Tracer = "prestateTracer"
)

// DefaultTraceConfig configures the built-in struct-logger tracer.
type DefaultTraceConfig struct {
	EnableMemory     bool `json:"enableMemory"`
	DisableStack     bool `json:"disableStack"`
	DisableStorage   bool `json:"disableStorage"`
	EnableReturnData bool `json:"enableReturnData"`
	Debug            bool `json:"debug"`
	Limit            int  `json:"limit"`
}

// CallTracerConfig configures the callTracer.
type CallTracerConfig struct {
	// OnlyTopCall disables subcall collection when true.
	OnlyTopCall bool `json:"onlyTopCall"`
	// WithLog includes event logs in the trace when true.
	WithLog bool `json:"withLog"`
}

// FlatCallTracerConfig configures the flatCallTracer (go-ethereum only).
type FlatCallTracerConfig struct {
	// ConvertParityErrors converts errors to parity format when true.
	ConvertParityErrors bool `json:"convertParityErrors"`
	// IncludePrecompiles includes calls to precompiled contracts when true.
	IncludePrecompiles bool `json:"includePrecompiles"`
}

// PrestateTracerConfig configures the prestateTracer.
type PrestateTracerConfig struct {
	// DiffMode returns state modifications instead of pre-state when true.
	DiffMode bool `json:"diffMode"`
	// DisableCode omits contract bytecode from the result when true.
	DisableCode bool `json:"disableCode"`
	// DisableStorage omits contract storage from the result when true.
	DisableStorage bool `json:"disableStorage"`
	// IncludeEmpty includes empty state objects in the result when true.
	IncludeEmpty bool `json:"includeEmpty"`
}

// TraceConfig is the generic trace configuration accepted by all debug_trace* methods.
type TraceConfig struct {
	*DefaultTraceConfig

	Tracer Tracer `json:"tracer,omitempty"`
	// Timeout is the maximum duration the tracer is allowed to run.
	Timeout string `json:"timeout,omitempty"`
	// Reexec is the maximum number of blocks the tracer will re-execute
	// to reconstruct missing historical state.
	Reexec *uint64 `json:"reexec,omitempty"`
	// TracerConfig holds tracer-specific options whose concrete type depends
	// on the value of Tracer (e.g. [CallTracerConfig] for [CallTracer]).
	TracerConfig any `json:"tracerConfig,omitempty"`
}

func assignIndexCalls(callFrame *evmctypes.CallFrame) {
	if callFrame == nil {
		return
	}
	var index uint64
	callFrame.Index = index // root is always 0
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

// newTracerConfig builds a [TraceConfig] for a named tracer with optional
// timeout, reexec, and tracer-specific config.
func newTracerConfig(tracer Tracer, timeout time.Duration, reexec *uint64, cfg any) *TraceConfig {
	tc := &TraceConfig{Tracer: tracer, Timeout: timeout.String(), Reexec: reexec}
	if cfg != nil {
		tc.TracerConfig = cfg
	}
	return tc
}

// newCustomTracerConfig builds a [TraceConfig] for a custom JavaScript tracer.
func newCustomTracerConfig(jsTracer string, timeout time.Duration, reexec *uint64) *TraceConfig {
	return &TraceConfig{
		Tracer:  Tracer(jsTracer),
		Timeout: timeout.String(),
		Reexec:  reexec,
	}
}

// ─── TraceBlockByNumber ──────────────────────────────────────────────────────

// TraceBlockByNumber replays a block identified by number and returns raw trace
// results using cfg (nil uses the default struct-logger tracer).
func (d *debugNamespace) TraceBlockByNumber(blockNumber uint64, cfg *TraceConfig) ([]*evmctypes.TraceResult, error) {
	return d.TraceBlockByNumberWithContext(context.Background(), blockNumber, cfg)
}

// TraceBlockByNumberWithContext is the context-aware variant of [debugNamespace.TraceBlockByNumber].
func (d *debugNamespace) TraceBlockByNumberWithContext(
	ctx context.Context,
	blockNumber uint64,
	cfg *TraceConfig,
) ([]*evmctypes.TraceResult, error) {
	result := []*evmctypes.TraceResult{}
	if err := d.traceBlockByNumber(ctx, blockNumber, cfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TraceBlockByNumberCallTracer replays a block by number using the callTracer.
func (d *debugNamespace) TraceBlockByNumberCallTracer(
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) ([]*evmctypes.CallTracer, error) {
	return d.TraceBlockByNumberWithContextCallTracer(context.Background(), blockNumber, timeout, reexec, cfg)
}

// TraceBlockByNumberWithContextCallTracer is the context-aware variant of [debugNamespace.TraceBlockByNumberCallTracer].
func (d *debugNamespace) TraceBlockByNumberWithContextCallTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) ([]*evmctypes.CallTracer, error) {
	return d.traceBlockByNumberCallTracer(ctx, blockNumber, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByNumberCallTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) ([]*evmctypes.CallTracer, error) {
	callTracers := []*evmctypes.CallTracer{}
	traceCfg := newTracerConfig(CallTracer, timeout, reexec, cfg)
	if err := d.traceBlockByNumber(ctx, blockNumber, traceCfg, &callTracers); err != nil {
		return nil, err
	}
	for _, callTracer := range callTracers {
		assignIndexCalls(callTracer.Result)
	}
	return callTracers, nil
}

// TraceBlockByNumberFlatCallTracer replays a block by number using the flatCallTracer (go-ethereum only).
func (d *debugNamespace) TraceBlockByNumberFlatCallTracer(
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallTracer, error) {
	return d.TraceBlockByNumberWithContextFlatCallTracer(context.Background(), blockNumber, timeout, reexec, cfg)
}

// TraceBlockByNumberWithContextFlatCallTracer is the context-aware variant of [debugNamespace.TraceBlockByNumberFlatCallTracer].
// Only supported by go-ethereum clients.
func (d *debugNamespace) TraceBlockByNumberWithContextFlatCallTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallTracer, error) {
	return d.traceBlockByNumberFlatCallTracer(ctx, blockNumber, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByNumberFlatCallTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallTracer, error) {
	flatCallTracers := []*evmctypes.FlatCallTracer{}
	traceCfg := newTracerConfig(FlatCallTracer, timeout, reexec, cfg)
	if err := d.traceBlockByNumber(ctx, blockNumber, traceCfg, &flatCallTracers); err != nil {
		return nil, err
	}
	assignIndexFlatCalls(flatCallTracers)
	return flatCallTracers, nil
}

// TraceBlockByNumberPrestateTracer replays a block by number using the prestateTracer.
func (d *debugNamespace) TraceBlockByNumberPrestateTracer(
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) ([]*evmctypes.PrestateTracer, error) {
	return d.TraceBlockByNumberWithContextPrestateTracer(context.Background(), blockNumber, timeout, reexec, cfg)
}

// TraceBlockByNumberWithContextPrestateTracer is the context-aware variant of [debugNamespace.TraceBlockByNumberPrestateTracer].
func (d *debugNamespace) TraceBlockByNumberWithContextPrestateTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) ([]*evmctypes.PrestateTracer, error) {
	return d.traceBlockByNumberPrestateTracer(ctx, blockNumber, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByNumberPrestateTracer(
	ctx context.Context,
	blockNumber uint64,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) ([]*evmctypes.PrestateTracer, error) {
	prestateTracers := []*evmctypes.PrestateTracer{}
	traceCfg := newTracerConfig(PrestateTracer, timeout, reexec, cfg)
	if err := d.traceBlockByNumber(ctx, blockNumber, traceCfg, &prestateTracers); err != nil {
		return nil, err
	}
	return prestateTracers, nil
}

// TraceBlockByNumberCustomTracer replays a block by number using a custom JavaScript tracer.
func (d *debugNamespace) TraceBlockByNumberCustomTracer(
	blockNumber uint64,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) ([]*evmctypes.CustomTraceResult, error) {
	return d.TraceBlockByNumberWithContextCustomTracer(context.Background(), blockNumber, jsTracer, timeout, reexec)
}

// TraceBlockByNumberWithContextCustomTracer is the context-aware variant of [debugNamespace.TraceBlockByNumberCustomTracer].
func (d *debugNamespace) TraceBlockByNumberWithContextCustomTracer(
	ctx context.Context,
	blockNumber uint64,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) ([]*evmctypes.CustomTraceResult, error) {
	result := []*evmctypes.CustomTraceResult{}
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
	result any,
) error {
	params := []any{hexutil.EncodeUint64(blockNumber)}
	if traceCfg != nil {
		params = append(params, *traceCfg)
	}
	return d.c.call(ctx, result, DebugTraceBlockByNumber, params...)
}

// ─── TraceBlockByHash ────────────────────────────────────────────────────────

// TraceBlockByHash replays a block identified by hash and returns raw trace
// results using cfg (nil uses the default struct-logger tracer).
func (d *debugNamespace) TraceBlockByHash(hash string, cfg *TraceConfig) ([]*evmctypes.TraceResult, error) {
	return d.TraceBlockByHashWithContext(context.Background(), hash, cfg)
}

// TraceBlockByHashWithContext is the context-aware variant of [debugNamespace.TraceBlockByHash].
func (d *debugNamespace) TraceBlockByHashWithContext(
	ctx context.Context,
	hash string,
	cfg *TraceConfig,
) ([]*evmctypes.TraceResult, error) {
	result := []*evmctypes.TraceResult{}
	if err := d.traceBlockByHash(ctx, hash, cfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TraceBlockByHashCallTracer replays a block by hash using the callTracer.
func (d *debugNamespace) TraceBlockByHashCallTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) ([]*evmctypes.CallTracer, error) {
	return d.TraceBlockByHashWithContextCallTracer(context.Background(), hash, timeout, reexec, cfg)
}

// TraceBlockByHashWithContextCallTracer is the context-aware variant of [debugNamespace.TraceBlockByHashCallTracer].
func (d *debugNamespace) TraceBlockByHashWithContextCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) ([]*evmctypes.CallTracer, error) {
	return d.traceBlockByHashCallTracer(ctx, hash, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByHashCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) ([]*evmctypes.CallTracer, error) {
	callTracers := []*evmctypes.CallTracer{}
	traceCfg := newTracerConfig(CallTracer, timeout, reexec, cfg)
	if err := d.traceBlockByHash(ctx, hash, traceCfg, &callTracers); err != nil {
		return nil, err
	}
	for _, callTracer := range callTracers {
		assignIndexCalls(callTracer.Result)
	}
	return callTracers, nil
}

// TraceBlockByHashFlatCallTracer replays a block by hash using the flatCallTracer (go-ethereum only).
func (d *debugNamespace) TraceBlockByHashFlatCallTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallTracer, error) {
	return d.TraceBlockByHashWithContextFlatCallTracer(context.Background(), hash, timeout, reexec, cfg)
}

// TraceBlockByHashWithContextFlatCallTracer is the context-aware variant of [debugNamespace.TraceBlockByHashFlatCallTracer].
// Only supported by go-ethereum clients.
func (d *debugNamespace) TraceBlockByHashWithContextFlatCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallTracer, error) {
	return d.traceBlockByHashFlatCallTracer(ctx, hash, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByHashFlatCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallTracer, error) {
	flatCallTracers := []*evmctypes.FlatCallTracer{}
	traceCfg := newTracerConfig(FlatCallTracer, timeout, reexec, cfg)
	if err := d.traceBlockByHash(ctx, hash, traceCfg, &flatCallTracers); err != nil {
		return nil, err
	}
	assignIndexFlatCalls(flatCallTracers)
	return flatCallTracers, nil
}

// TraceBlockByHashPrestateTracer replays a block by hash using the prestateTracer.
func (d *debugNamespace) TraceBlockByHashPrestateTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) ([]*evmctypes.PrestateTracer, error) {
	return d.TraceBlockByHashWithContextPrestateTracer(context.Background(), hash, timeout, reexec, cfg)
}

// TraceBlockByHashWithContextPrestateTracer is the context-aware variant of [debugNamespace.TraceBlockByHashPrestateTracer].
func (d *debugNamespace) TraceBlockByHashWithContextPrestateTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) ([]*evmctypes.PrestateTracer, error) {
	return d.traceBlockByHashPrestateTracer(ctx, hash, timeout, reexec, cfg)
}

func (d *debugNamespace) traceBlockByHashPrestateTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) ([]*evmctypes.PrestateTracer, error) {
	prestateTracers := []*evmctypes.PrestateTracer{}
	traceCfg := newTracerConfig(PrestateTracer, timeout, reexec, cfg)
	if err := d.traceBlockByHash(ctx, hash, traceCfg, &prestateTracers); err != nil {
		return nil, err
	}
	return prestateTracers, nil
}

// TraceBlockByHashCustomTracer replays a block by hash using a custom JavaScript tracer.
func (d *debugNamespace) TraceBlockByHashCustomTracer(
	hash string,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) ([]*evmctypes.CustomTraceResult, error) {
	return d.TraceBlockByHashWithContextCustomTracer(context.Background(), hash, jsTracer, timeout, reexec)
}

// TraceBlockByHashWithContextCustomTracer is the context-aware variant of [debugNamespace.TraceBlockByHashCustomTracer].
func (d *debugNamespace) TraceBlockByHashWithContextCustomTracer(
	ctx context.Context,
	hash string,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) ([]*evmctypes.CustomTraceResult, error) {
	result := []*evmctypes.CustomTraceResult{}
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
	result any,
) error {
	params := []any{hash}
	if traceCfg != nil {
		params = append(params, *traceCfg)
	}
	return d.c.call(ctx, result, DebugTraceBlockByHash, params...)
}

// ─── TraceTransaction ────────────────────────────────────────────────────────

// TraceTransaction replays a transaction and returns the raw trace result
// using cfg (nil uses the default struct-logger tracer).
func (d *debugNamespace) TraceTransaction(hash string, cfg *TraceConfig) (any, error) {
	return d.TraceTransactionWithContext(context.Background(), hash, cfg)
}

// TraceTransactionWithContext is the context-aware variant of [debugNamespace.TraceTransaction].
func (d *debugNamespace) TraceTransactionWithContext(
	ctx context.Context,
	hash string,
	cfg *TraceConfig,
) (any, error) {
	var result any
	if err := d.traceTransaction(ctx, hash, cfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TraceTransactionCallTracer replays a transaction using the callTracer.
func (d *debugNamespace) TraceTransactionCallTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (*evmctypes.CallFrame, error) {
	return d.TraceTransactionWithContextCallTracer(context.Background(), hash, timeout, reexec, cfg)
}

// TraceTransactionWithContextCallTracer is the context-aware variant of [debugNamespace.TraceTransactionCallTracer].
func (d *debugNamespace) TraceTransactionWithContextCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (*evmctypes.CallFrame, error) {
	callFrame := &evmctypes.CallFrame{}
	traceCfg := newTracerConfig(CallTracer, timeout, reexec, cfg)
	if err := d.traceTransaction(ctx, hash, traceCfg, callFrame); err != nil {
		return nil, err
	}
	assignIndexCalls(callFrame)
	return callFrame, nil
}

// TraceTransactionFlatCallTracer replays a transaction using the flatCallTracer.
func (d *debugNamespace) TraceTransactionFlatCallTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallFrame, error) {
	return d.TraceTransactionWithContextFlatCallTracer(context.Background(), hash, timeout, reexec, cfg)
}

// TraceTransactionWithContextFlatCallTracer is the context-aware variant of [debugNamespace.TraceTransactionFlatCallTracer].
func (d *debugNamespace) TraceTransactionWithContextFlatCallTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallFrame, error) {
	flatCallFrames := []*evmctypes.FlatCallFrame{}
	traceCfg := newTracerConfig(FlatCallTracer, timeout, reexec, cfg)
	if err := d.traceTransaction(ctx, hash, traceCfg, &flatCallFrames); err != nil {
		return nil, err
	}
	return flatCallFrames, nil
}

// TraceTransactionPrestateTracer replays a transaction using the prestateTracer.
func (d *debugNamespace) TraceTransactionPrestateTracer(
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (*evmctypes.PrestateResult, error) {
	return d.TraceTransactionWithContextPrestateTracer(context.Background(), hash, timeout, reexec, cfg)
}

// TraceTransactionWithContextPrestateTracer is the context-aware variant of [debugNamespace.TraceTransactionPrestateTracer].
func (d *debugNamespace) TraceTransactionWithContextPrestateTracer(
	ctx context.Context,
	hash string,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (*evmctypes.PrestateResult, error) {
	prestateResult := &evmctypes.PrestateResult{}
	traceCfg := newTracerConfig(PrestateTracer, timeout, reexec, cfg)
	if err := d.traceTransaction(ctx, hash, traceCfg, prestateResult); err != nil {
		return nil, err
	}
	return prestateResult, nil
}

// TraceTransactionCustomTracer replays a transaction using a custom JavaScript tracer.
func (d *debugNamespace) TraceTransactionCustomTracer(
	hash string,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (json.RawMessage, error) {
	return d.TraceTransactionWithContextCustomTracer(context.Background(), hash, jsTracer, timeout, reexec)
}

// TraceTransactionWithContextCustomTracer is the context-aware variant of [debugNamespace.TraceTransactionCustomTracer].
func (d *debugNamespace) TraceTransactionWithContextCustomTracer(
	ctx context.Context,
	hash string,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (json.RawMessage, error) {
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
	result any,
) error {
	params := []any{hash}
	if traceCfg != nil {
		params = append(params, *traceCfg)
	}
	return d.c.call(ctx, result, DebugTraceTransaction, params...)
}

// ─── TraceCall ───────────────────────────────────────────────────────────────

// TraceCall simulates tx against blockAndTag and returns the raw trace result
// using cfg (nil uses the default struct-logger tracer).
func (d *debugNamespace) TraceCall(tx *Tx, blockAndTag evmctypes.BlockAndTag, cfg *TraceConfig) (any, error) {
	return d.TraceCallWithContext(context.Background(), tx, blockAndTag, cfg)
}

// TraceCallWithContext is the context-aware variant of [debugNamespace.TraceCall].
func (d *debugNamespace) TraceCallWithContext(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	cfg *TraceConfig,
) (any, error) {
	var result any
	if err := d.traceCall(ctx, tx, blockAndTag, cfg, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// TraceCallCallTracer simulates tx using the callTracer.
func (d *debugNamespace) TraceCallCallTracer(
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (*evmctypes.CallFrame, error) {
	return d.TraceCallWithContextCallTracer(context.Background(), tx, blockAndTag, timeout, reexec, cfg)
}

// TraceCallWithContextCallTracer is the context-aware variant of [debugNamespace.TraceCallCallTracer].
func (d *debugNamespace) TraceCallWithContextCallTracer(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *CallTracerConfig,
) (*evmctypes.CallFrame, error) {
	callFrame := &evmctypes.CallFrame{}
	traceCfg := newTracerConfig(CallTracer, timeout, reexec, cfg)
	if err := d.traceCall(ctx, tx, blockAndTag, traceCfg, callFrame); err != nil {
		return nil, err
	}
	assignIndexCalls(callFrame)
	return callFrame, nil
}

// TraceCallFlatCallTracer simulates tx using the flatCallTracer.
func (d *debugNamespace) TraceCallFlatCallTracer(
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallFrame, error) {
	return d.TraceCallWithContextFlatCallTracer(context.Background(), tx, blockAndTag, timeout, reexec, cfg)
}

// TraceCallWithContextFlatCallTracer is the context-aware variant of [debugNamespace.TraceCallFlatCallTracer].
func (d *debugNamespace) TraceCallWithContextFlatCallTracer(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *FlatCallTracerConfig,
) ([]*evmctypes.FlatCallFrame, error) {
	flatCallFrames := []*evmctypes.FlatCallFrame{}
	traceCfg := newTracerConfig(FlatCallTracer, timeout, reexec, cfg)
	if err := d.traceCall(ctx, tx, blockAndTag, traceCfg, &flatCallFrames); err != nil {
		return nil, err
	}
	return flatCallFrames, nil
}

// TraceCallPrestateTracer simulates tx using the prestateTracer.
func (d *debugNamespace) TraceCallPrestateTracer(
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (*evmctypes.PrestateResult, error) {
	return d.TraceCallWithContextPrestateTracer(context.Background(), tx, blockAndTag, timeout, reexec, cfg)
}

// TraceCallWithContextPrestateTracer is the context-aware variant of [debugNamespace.TraceCallPrestateTracer].
func (d *debugNamespace) TraceCallWithContextPrestateTracer(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	timeout time.Duration,
	reexec *uint64,
	cfg *PrestateTracerConfig,
) (*evmctypes.PrestateResult, error) {
	prestateResult := &evmctypes.PrestateResult{}
	traceCfg := newTracerConfig(PrestateTracer, timeout, reexec, cfg)
	if err := d.traceCall(ctx, tx, blockAndTag, traceCfg, prestateResult); err != nil {
		return nil, err
	}
	return prestateResult, nil
}

// TraceCallCustomTracer simulates tx using a custom JavaScript tracer.
func (d *debugNamespace) TraceCallCustomTracer(
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (json.RawMessage, error) {
	return d.TraceCallWithContextCustomTracer(context.Background(), tx, blockAndTag, jsTracer, timeout, reexec)
}

// TraceCallWithContextCustomTracer is the context-aware variant of [debugNamespace.TraceCallCustomTracer].
func (d *debugNamespace) TraceCallWithContextCustomTracer(
	ctx context.Context,
	tx *Tx,
	blockAndTag evmctypes.BlockAndTag,
	jsTracer string,
	timeout time.Duration,
	reexec *uint64,
) (json.RawMessage, error) {
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
	result any,
) error {
	msg, err := tx.parseCallMsg()
	if err != nil {
		return err
	}
	params := []any{msg, blockAndTag.String()}
	if traceCfg != nil {
		params = append(params, *traceCfg)
	}
	return d.c.call(ctx, result, DebugTraceCall, params...)
}

// ─── Raw / Getter Methods ───────────────────────────────────────────────────

// GetRawHeader returns the RLP-encoded header for blockAndTag as a hex string.
func (d *debugNamespace) GetRawHeader(blockAndTag evmctypes.BlockAndTag) (string, error) {
	return d.GetRawHeaderWithContext(context.Background(), blockAndTag)
}

// GetRawHeaderWithContext is the context-aware variant of [debugNamespace.GetRawHeader].
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

// GetRawBlock returns the RLP-encoded block for blockAndTag as a hex string.
func (d *debugNamespace) GetRawBlock(blockAndTag evmctypes.BlockAndTag) (string, error) {
	return d.GetRawBlockWithContext(context.Background(), blockAndTag)
}

// GetRawBlockWithContext is the context-aware variant of [debugNamespace.GetRawBlock].
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

// GetRawTransaction returns the RLP-encoded transaction identified by hash as a hex string.
func (d *debugNamespace) GetRawTransaction(hash string) (string, error) {
	return d.GetRawTransactionWithContext(context.Background(), hash)
}

// GetRawTransactionWithContext is the context-aware variant of [debugNamespace.GetRawTransaction].
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

// GetRawReceipts returns RLP-encoded receipts for all transactions in blockAndTag.
func (d *debugNamespace) GetRawReceipts(blockAndTag evmctypes.BlockAndTag) ([]string, error) {
	return d.GetRawReceiptsWithContext(context.Background(), blockAndTag)
}

// GetRawReceiptsWithContext is the context-aware variant of [debugNamespace.GetRawReceipts].
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

// GetBadBlocks returns the list of bad blocks known to the node.
func (d *debugNamespace) GetBadBlocks() ([]*evmctypes.BadBlock, error) {
	return d.GetBadBlocksWithContext(context.Background())
}

// GetBadBlocksWithContext is the context-aware variant of [debugNamespace.GetBadBlocks].
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
