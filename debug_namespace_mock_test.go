package evmc

import (
	"encoding/json"
	"testing"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// callTracerResultJSON은 callTracer 결과 mock 데이터.
// txHash는 결과 객체를 구분하는 트랜잭션 해시다.
func callTracerResultJSON(txHash string) map[string]interface{} {
	return map[string]interface{}{
		"txHash": txHash,
		"result": map[string]interface{}{
			"from": "0xfrom",
			// 0x5208 = 21000 (standard ETH transfer gas)
			"gas":     "0x5208",
			"gasUsed": "0x5208",
			"to":      "0xto",
			"input":   "0x",
			"output":  "0xresult",
			// 0xde0b6b3a7640000 = 1000000000000000000 (1 ETH in wei)
			"value": "0xde0b6b3a7640000",
			"type":  "CALL",
		},
	}
}

// traceResultJSON은 structLogger(기본 tracer) 결과 mock 데이터.
// txHash는 결과 객체를 구분하는 트랜잭션 해시다.
func traceResultJSON(txHash string) map[string]interface{} {
	return map[string]interface{}{
		"txHash": txHash,
		"result": map[string]interface{}{
			// 21000: standard ETH transfer gas cost
			"gas":         21000,
			"failed":      false,
			"returnValue": "",
			"structLogs":  []interface{}{},
		},
	}
}

// ─── TraceBlockByNumber ──────────────────────────────────────────────────────

func Test_debugNamespace_mock_TraceBlockByNumber(t *testing.T) {
	client := testWithMock(t, "debug_traceBlockByNumber", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			traceResultJSON("0xtx1"),
			traceResultJSON("0xtx2"),
		}
	})
	results, err := client.Debug().TraceBlockByNumber(100, nil)
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, "0xtx1", results[0].TxHash)
	assert.Equal(t, "0xtx2", results[1].TxHash)
}

func Test_debugNamespace_mock_TraceBlockByNumber_callTracer(t *testing.T) {
	client := testWithMock(t, "debug_traceBlockByNumber", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			callTracerResultJSON("0xtx1"),
			callTracerResultJSON("0xtx2"),
		}
	})
	results, err := client.Debug().TraceBlockByNumber_callTracer(100, 0, nil, nil)
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, "0xtx1", results[0].TxHash)
	require.NotNil(t, results[0].Result)
	assert.Equal(t, "CALL", results[0].Result.Type)
	assert.Equal(t, "0xfrom", results[0].Result.From)
}

// ─── TraceBlockByHash ────────────────────────────────────────────────────────

func Test_debugNamespace_mock_TraceBlockByHash(t *testing.T) {
	client := testWithMock(t, "debug_traceBlockByHash", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			traceResultJSON("0xtx1"),
			traceResultJSON("0xtx2"),
		}
	})
	results, err := client.Debug().TraceBlockByHash("0xblockhash", nil)
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, "0xtx1", results[0].TxHash)
	assert.Equal(t, "0xtx2", results[1].TxHash)
}

func Test_debugNamespace_mock_TraceBlockByHash_callTracer(t *testing.T) {
	client := testWithMock(t, "debug_traceBlockByHash", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			callTracerResultJSON("0xtx1"),
		}
	})
	results, err := client.Debug().TraceBlockByHash_callTracer("0xblockhash", 0, nil, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "0xtx1", results[0].TxHash)
}

func Test_debugNamespace_mock_TraceBlockByHash_prestateTracer(t *testing.T) {
	client := testWithMock(t, "debug_traceBlockByHash", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			{
				"txHash": "0xtx1",
				"result": map[string]interface{}{
					"0xaddr": map[string]interface{}{
						"balance": "0xde0b6b3a7640000",
						"nonce":   "0x5",
					},
				},
			},
		}
	})
	results, err := client.Debug().TraceBlockByHash_prestateTracer("0xblockhash", 0, nil, nil)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "0xtx1", results[0].TxHash)
}

// ─── TraceTransaction ────────────────────────────────────────────────────────

func Test_debugNamespace_mock_TraceTransaction(t *testing.T) {
	client := testWithMock(t, "debug_traceTransaction", func(params json.RawMessage) interface{} {
		return map[string]interface{}{
			// 21000: standard ETH transfer gas cost
			"gas":         21000,
			"failed":      false,
			"returnValue": "",
			"structLogs":  []interface{}{},
		}
	})
	result, err := client.Debug().TraceTransaction("0xtxhash", nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func Test_debugNamespace_mock_TraceTransaction_callTracer(t *testing.T) {
	client := testWithMock(t, "debug_traceTransaction", func(params json.RawMessage) interface{} {
		return map[string]interface{}{
			"from": "0xfrom",
			// 0x5208 = 21000 (standard ETH transfer gas)
			"gas":     "0x5208",
			"gasUsed": "0x5208",
			"to":      "0xto",
			"input":   "0x",
			"output":  "0xresult",
			// 0xde0b6b3a7640000 = 1000000000000000000 (1 ETH in wei)
			"value": "0xde0b6b3a7640000",
			"type":  "CALL",
		}
	})
	callFrame, err := client.Debug().TraceTransaction_callTracer("0xtxhash", 0, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, callFrame)
	assert.Equal(t, "0xfrom", callFrame.From)
	assert.Equal(t, "CALL", callFrame.Type)
}

func Test_debugNamespace_mock_TraceTransaction_prestateTracer(t *testing.T) {
	client := testWithMock(t, "debug_traceTransaction", func(params json.RawMessage) interface{} {
		return map[string]interface{}{
			"0xfrom": map[string]interface{}{
				// 0xde0b6b3a7640000 = 1000000000000000000 (1 ETH in wei)
				"balance": "0xde0b6b3a7640000",
				// 0x5 = 5
				"nonce": "0x5",
			},
			"0xto": map[string]interface{}{
				"balance": "0x0",
				"nonce":   "0x0",
				"code":    "0x60806040",
			},
		}
	})
	result, err := client.Debug().TraceTransaction_prestateTracer("0xtxhash", 0, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

// ─── TraceCall ───────────────────────────────────────────────────────────────

func Test_debugNamespace_mock_TraceCall(t *testing.T) {
	client := testWithMock(t, "debug_traceCall", func(params json.RawMessage) interface{} {
		return map[string]interface{}{
			"gas":         21000,
			"failed":      false,
			"returnValue": "0x",
			"structLogs":  []interface{}{},
		}
	})
	tx := &Tx{To: "0x000000000000000000000000000000000000dead", Data: "0x"}
	result, err := client.Debug().TraceCall(tx, evmctypes.Latest, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

func Test_debugNamespace_mock_TraceCall_callTracer(t *testing.T) {
	client := testWithMock(t, "debug_traceCall", func(params json.RawMessage) interface{} {
		return map[string]interface{}{
			"from":    "0x0000000000000000000000000000000000000000",
			"gas":     "0x5208",
			"gasUsed": "0x5208",
			"to":      "0x000000000000000000000000000000000000dead",
			"input":   "0x18160ddd",
			"output":  "0x",
			"value":   "0x0",
			"type":    "CALL",
		}
	})
	tx := &Tx{To: "0x000000000000000000000000000000000000dead", Data: "0x18160ddd"}
	callFrame, err := client.Debug().TraceCall_callTracer(tx, evmctypes.Latest, 0, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, callFrame)
	assert.Equal(t, "CALL", callFrame.Type)
	assert.Equal(t, "0x000000000000000000000000000000000000dead", *callFrame.To)
}

func Test_debugNamespace_mock_TraceCall_prestateTracer(t *testing.T) {
	client := testWithMock(t, "debug_traceCall", func(params json.RawMessage) interface{} {
		return map[string]interface{}{
			"0xfrom": map[string]interface{}{
				"balance": "0xde0b6b3a7640000",
				"nonce":   "0x1",
			},
		}
	})
	tx := &Tx{To: "0x000000000000000000000000000000000000dead", Data: "0x"}
	result, err := client.Debug().TraceCall_prestateTracer(tx, evmctypes.Latest, 0, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
}

// ─── Raw / Getter Methods ───────────────────────────────────────────────────

func Test_debugNamespace_mock_GetRawHeader(t *testing.T) {
	client := testWithMock(t, "debug_getRawHeader", func(params json.RawMessage) interface{} {
		return "0xf9020aa0"
	})
	raw, err := client.Debug().GetRawHeader(evmctypes.Latest)
	require.NoError(t, err)
	assert.Equal(t, "0xf9020aa0", raw)
}

func Test_debugNamespace_mock_GetRawBlock(t *testing.T) {
	client := testWithMock(t, "debug_getRawBlock", func(params json.RawMessage) interface{} {
		return "0xf9020af9"
	})
	raw, err := client.Debug().GetRawBlock(evmctypes.Latest)
	require.NoError(t, err)
	assert.Equal(t, "0xf9020af9", raw)
}

func Test_debugNamespace_mock_GetRawTransaction(t *testing.T) {
	client := testWithMock(t, "debug_getRawTransaction", func(params json.RawMessage) interface{} {
		return "0x02f8748201"
	})
	raw, err := client.Debug().GetRawTransaction("0xtxhash")
	require.NoError(t, err)
	assert.Equal(t, "0x02f8748201", raw)
}

func Test_debugNamespace_mock_GetRawReceipts(t *testing.T) {
	client := testWithMock(t, "debug_getRawReceipts", func(params json.RawMessage) interface{} {
		return []string{"0xf9010a", "0xf9020b"}
	})
	receipts, err := client.Debug().GetRawReceipts(evmctypes.Latest)
	require.NoError(t, err)
	require.Len(t, receipts, 2)
	assert.Equal(t, "0xf9010a", receipts[0])
	assert.Equal(t, "0xf9020b", receipts[1])
}

func Test_debugNamespace_mock_GetBadBlocks(t *testing.T) {
	client := testWithMock(t, "debug_getBadBlocks", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			{
				"hash": "0xbadblockhash",
				"rlp":  "0xf9020a",
				"block": map[string]interface{}{
					"number":           "0x1",
					"hash":             "0xbadblockhash",
					"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
					"nonce":            "0x0000000000000000",
					"mixHash":          "0x0000000000000000000000000000000000000000000000000000000000000000",
					"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
					"logsBloom":        "0x00000000",
					"stateRoot":        "0xabcd",
					"miner":            "0x1234567890123456789012345678901234567890",
					"difficulty":       "0x0",
					"extraData":        "0x",
					"gasLimit":         "0x1c9c380",
					"gasUsed":          "0x0",
					"timestamp":        "0x64",
					"transactionsRoot": "0xabc",
					"receiptsRoot":     "0xdef",
					"totalDifficulty":  "0x0",
					"size":             "0x200",
					"uncles":           []string{},
					"transactions":     []string{},
				},
			},
		}
	})
	badBlocks, err := client.Debug().GetBadBlocks()
	require.NoError(t, err)
	require.Len(t, badBlocks, 1)
	assert.Equal(t, "0xbadblockhash", badBlocks[0].Hash)
	assert.Equal(t, "0xf9020a", badBlocks[0].RLP)
	require.NotNil(t, badBlocks[0].Block)
	assert.Equal(t, uint64(1), badBlocks[0].Block.Number)
}
