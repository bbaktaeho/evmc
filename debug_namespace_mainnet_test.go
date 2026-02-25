package evmc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCustomJSTracer는 gasUsed를 반환하는 간단한 JavaScript 트레이서다.
const testCustomJSTracer = `{result: function(ctx) { return {gasUsed: ctx.gasUsed}; }, fault: function() {}}`

// 블록 20,000,000 은 트랜잭션이 134개로 trace가 무거우므로
// 더 가벼운 블록을 사용한다.
const (
	// 블록 20,000,010: 트랜잭션이 적은 블록
	debugRefBlockNumber = uint64(20_000_010)
	debugRefBlockHash   = "0xfdf90530404e89e281d0b64b941d2601de1be4f7ec9dccfe935e19fc4f5ab3bf"

	// 간단한 ETH 전송 트랜잭션 (block 20,000,000 내 첫 번째 tx)
	debugRefTxHash = mainnetRefEip1559TxHash
)

// ─── TraceBlockByNumber ──────────────────────────────────────────────────────

func Test_Mainnet_Debug_TraceBlockByNumber_callTracer(t *testing.T) {
	c := mainnetClient(t)
	results, err := c.Debug().TraceBlockByNumber_callTracer(debugRefBlockNumber, time.Minute, nil, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.NotEmpty(t, r.TxHash)
		require.NotNil(t, r.Result)
		assert.NotEmpty(t, r.Result.From)
		assert.NotEmpty(t, r.Result.Type)
	}
	// Index 검증: 루트 CallFrame의 Index는 0
	assert.Equal(t, uint64(0), results[0].Result.Index)
}

func Test_Mainnet_Debug_TraceBlockByNumber_callTracer_withLog(t *testing.T) {
	c := mainnetClient(t)
	results, err := c.Debug().TraceBlockByNumber_callTracer(debugRefBlockNumber, time.Minute, nil, &CallTracerConfig{WithLog: true})
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	// 하나 이상의 트랜잭션에 로그가 있어야 함
	var hasLogs bool
	for _, r := range results {
		if r.Result != nil && len(r.Result.Logs) > 0 {
			hasLogs = true
			break
		}
	}
	assert.True(t, hasLogs, "at least one transaction should have logs with withLog=true")
}

func Test_Mainnet_Debug_TraceBlockByNumber_prestateTracer(t *testing.T) {
	c := mainnetClient(t)
	results, err := c.Debug().TraceBlockByNumber_prestateTracer(debugRefBlockNumber, time.Minute, nil, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.NotEmpty(t, r.TxHash)
		assert.NotEmpty(t, r.Result, "prestate result should not be empty")
		frame, err := r.ParseFrames()
		require.NoError(t, err)
		assert.NotEmpty(t, frame, "parsed prestate frame should have accounts")
	}
}

func Test_Mainnet_Debug_TraceBlockByNumber_prestateTracer_diffMode(t *testing.T) {
	c := mainnetClient(t)
	results, err := c.Debug().TraceBlockByNumber_prestateTracer(debugRefBlockNumber, time.Minute, nil, &PrestateTracerConfig{DiffMode: true})
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.NotEmpty(t, r.TxHash)
		diffFrame, err := r.ParseDiffFrames()
		require.NoError(t, err)
		require.NotNil(t, diffFrame, "diff frame should not be nil")
		assert.NotEmpty(t, diffFrame.Pre, "pre state should have accounts")
	}
}

// ─── TraceBlockByHash ────────────────────────────────────────────────────────

func Test_Mainnet_Debug_TraceBlockByHash_callTracer(t *testing.T) {
	c := mainnetClient(t)
	results, err := c.Debug().TraceBlockByHash_callTracer(debugRefBlockHash, time.Minute, nil, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.NotEmpty(t, r.TxHash)
		require.NotNil(t, r.Result)
		assert.NotEmpty(t, r.Result.From)
		assert.NotEmpty(t, r.Result.Type)
	}
}

func Test_Mainnet_Debug_TraceBlockByHash_prestateTracer(t *testing.T) {
	c := mainnetClient(t)
	results, err := c.Debug().TraceBlockByHash_prestateTracer(debugRefBlockHash, time.Minute, nil, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.NotEmpty(t, r.TxHash)
		assert.NotEmpty(t, r.Result)
	}
}

// ─── TraceTransaction ────────────────────────────────────────────────────────

func Test_Mainnet_Debug_TraceTransaction_callTracer(t *testing.T) {
	c := mainnetClient(t)
	callFrame, err := c.Debug().TraceTransaction_callTracer(debugRefTxHash, time.Minute, nil, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, callFrame.From)
	assert.NotEmpty(t, callFrame.Type)
	assert.NotEmpty(t, callFrame.Gas)
	assert.NotEmpty(t, callFrame.GasUsed)
	assert.Equal(t, uint64(0), callFrame.Index)
}

func Test_Mainnet_Debug_TraceTransaction_callTracer_withLog(t *testing.T) {
	c := mainnetClient(t)
	callFrame, err := c.Debug().TraceTransaction_callTracer(debugRefTxHash, time.Minute, nil, &CallTracerConfig{WithLog: true})
	require.NoError(t, err)

	assert.NotEmpty(t, callFrame.From)
	assert.NotEmpty(t, callFrame.Type)
}

func Test_Mainnet_Debug_TraceTransaction_prestateTracer(t *testing.T) {
	c := mainnetClient(t)
	result, err := c.Debug().TraceTransaction_prestateTracer(debugRefTxHash, time.Minute, nil, nil)
	require.NoError(t, err)

	require.NotNil(t, result)
	frame, err := result.ParseFrame()
	require.NoError(t, err)
	assert.NotEmpty(t, frame, "prestate frame should have accounts")
}

func Test_Mainnet_Debug_TraceTransaction_prestateTracer_diffMode(t *testing.T) {
	c := mainnetClient(t)
	result, err := c.Debug().TraceTransaction_prestateTracer(debugRefTxHash, time.Minute, nil, &PrestateTracerConfig{DiffMode: true})
	require.NoError(t, err)

	require.NotNil(t, result)
	diffFrame, err := result.ParseDiffFrame()
	require.NoError(t, err)
	require.NotNil(t, diffFrame)
	assert.NotEmpty(t, diffFrame.Pre, "pre state should have accounts")
}

func Test_Mainnet_Debug_TraceTransaction_flatCallTracer(t *testing.T) {
	c := mainnetClient(t)
	flatCalls, err := c.Debug().TraceTransaction_flatCallTracer(debugRefTxHash, time.Minute, nil, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, flatCalls)
	assert.NotEmpty(t, flatCalls[0].Type)
}

func Test_Mainnet_Debug_TraceTransaction_default(t *testing.T) {
	c := mainnetClient(t)
	result, err := c.Debug().TraceTransaction(debugRefTxHash, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// ─── TraceCall ───────────────────────────────────────────────────────────────

func Test_Mainnet_Debug_TraceCall_callTracer(t *testing.T) {
	c := mainnetClient(t)
	// USDT totalSupply() 호출 trace
	tx := &Tx{
		To:   mainnetRefUSDT,
		Data: "0x18160ddd",
	}
	callFrame, err := c.Debug().TraceCall_callTracer(tx, evmctypes.FormatNumber(mainnetRefBlockNumber), 0, nil, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, callFrame.Type)
	assert.NotEmpty(t, callFrame.From)
	require.NotNil(t, callFrame.To)
	assert.Equal(t, "0xdac17f958d2ee523a2206206994597c13d831ec7", *callFrame.To)
	assert.NotEmpty(t, callFrame.Gas)
	assert.NotEmpty(t, callFrame.GasUsed)
	require.NotNil(t, callFrame.Output)
	assert.NotEmpty(t, *callFrame.Output)
	assert.Equal(t, uint64(0), callFrame.Index)
}

func Test_Mainnet_Debug_TraceCall_prestateTracer(t *testing.T) {
	c := mainnetClient(t)
	tx := &Tx{
		To:   mainnetRefUSDT,
		Data: "0x18160ddd",
	}
	result, err := c.Debug().TraceCall_prestateTracer(tx, evmctypes.Latest, 0, nil, nil)
	require.NoError(t, err)

	require.NotNil(t, result)
	frame, err := result.ParseFrame()
	require.NoError(t, err)
	assert.NotEmpty(t, frame, "prestate frame should have accounts")
}

func Test_Mainnet_Debug_TraceCall_prestateTracer_diffMode(t *testing.T) {
	c := mainnetClient(t)
	tx := &Tx{
		To:   mainnetRefUSDT,
		Data: "0x18160ddd",
	}
	result, err := c.Debug().TraceCall_prestateTracer(tx, evmctypes.Latest, 0, nil, &PrestateTracerConfig{DiffMode: true})
	require.NoError(t, err)

	require.NotNil(t, result)
	diffFrame, err := result.ParseDiffFrame()
	require.NoError(t, err)
	require.NotNil(t, diffFrame)
	assert.NotEmpty(t, diffFrame.Pre, "pre state should have accounts")
}

func Test_Mainnet_Debug_TraceCall_default(t *testing.T) {
	c := mainnetClient(t)
	tx := &Tx{
		To:   mainnetRefUSDT,
		Data: "0x18160ddd",
	}
	result, err := c.Debug().TraceCall(tx, evmctypes.Latest, nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// ─── Custom Tracer ───────────────────────────────────────────────────────────

func Test_Mainnet_Debug_TraceTransaction_customTracer(t *testing.T) {
	c := mainnetClient(t)
	result, err := c.Debug().TraceTransaction_customTracer(debugRefTxHash, testCustomJSTracer, time.Minute, nil)
	require.NoError(t, err)
	require.NotNil(t, result)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &parsed))
	assert.Contains(t, parsed, "gasUsed")
}

func Test_Mainnet_Debug_TraceBlockByNumber_customTracer(t *testing.T) {
	c := mainnetClient(t)
	results, err := c.Debug().TraceBlockByNumber_customTracer(debugRefBlockNumber, testCustomJSTracer, time.Minute, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.NotEmpty(t, r.TxHash)
		assert.NotEmpty(t, r.Result)
	}
}

func Test_Mainnet_Debug_TraceBlockByHash_customTracer(t *testing.T) {
	c := mainnetClient(t)
	results, err := c.Debug().TraceBlockByHash_customTracer(debugRefBlockHash, testCustomJSTracer, time.Minute, nil)
	require.NoError(t, err)

	assert.NotEmpty(t, results)
	for _, r := range results {
		assert.NotEmpty(t, r.TxHash)
		assert.NotEmpty(t, r.Result)
	}
}

func Test_Mainnet_Debug_TraceCall_customTracer(t *testing.T) {
	c := mainnetClient(t)
	tx := &Tx{
		To:   mainnetRefUSDT,
		Data: "0x18160ddd",
	}
	result, err := c.Debug().TraceCall_customTracer(tx, evmctypes.FormatNumber(mainnetRefBlockNumber), testCustomJSTracer, time.Minute, nil)
	require.NoError(t, err)
	require.NotNil(t, result)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal(result, &parsed))
	assert.Contains(t, parsed, "gasUsed")
}

// ─── Raw / Getter Methods ───────────────────────────────────────────────────

func Test_Mainnet_Debug_GetRawHeader(t *testing.T) {
	c := mainnetClient(t)
	raw, err := c.Debug().GetRawHeader(evmctypes.FormatNumber(mainnetRefBlockNumber))
	require.NoError(t, err)

	assert.NotEmpty(t, raw)
	assert.Equal(t, "0x", raw[:2])
}

func Test_Mainnet_Debug_GetRawBlock(t *testing.T) {
	c := mainnetClient(t)
	raw, err := c.Debug().GetRawBlock(evmctypes.FormatNumber(mainnetRefBlockNumber))
	require.NoError(t, err)

	assert.NotEmpty(t, raw)
	assert.Equal(t, "0x", raw[:2])
}

func Test_Mainnet_Debug_GetRawTransaction(t *testing.T) {
	c := mainnetClient(t)
	raw, err := c.Debug().GetRawTransaction(debugRefTxHash)
	require.NoError(t, err)

	assert.NotEmpty(t, raw)
	assert.Equal(t, "0x", raw[:2])
}

func Test_Mainnet_Debug_GetRawReceipts(t *testing.T) {
	c := mainnetClient(t)
	receipts, err := c.Debug().GetRawReceipts(evmctypes.FormatNumber(mainnetRefBlockNumber))
	require.NoError(t, err)

	// 블록 20M에는 134개 tx가 있으므로 134개 영수증
	assert.Len(t, receipts, 134)
	for _, r := range receipts {
		assert.Equal(t, "0x", r[:2])
	}
}

func Test_Mainnet_Debug_GetBadBlocks(t *testing.T) {
	c := mainnetClient(t)
	badBlocks, err := c.Debug().GetBadBlocks()
	require.NoError(t, err)

	// 정상 노드에서는 빈 배열이 일반적
	assert.NotNil(t, badBlocks)
}
