package evmc

import (
	"encoding/json"
	"testing"

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

func Test_debugNamespace_mock_TraceBlockByHash(t *testing.T) {
	client := testWithMock(t, "debug_traceBlockByHash", func(params json.RawMessage) interface{} {
		// TraceBlockByHash returns a single TraceResult (not an array)
		return traceResultJSON("0xtx1")
	})
	result, err := client.Debug().TraceBlockByHash("0xblockhash", nil)
	require.NoError(t, err)
	require.NotNil(t, result)
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
