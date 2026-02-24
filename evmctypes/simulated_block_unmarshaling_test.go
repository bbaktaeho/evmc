package evmctypes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimulateBlockResult_UnmarshalJSON(t *testing.T) {
	raw := map[string]interface{}{
		"number":        "0x16e3600", // 24000000
		"hash":          "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"timestamp":     "0x65a8c3a0", // 1705558944
		"gasLimit":      "0x1c9c380",  // 30000000
		"gasUsed":       "0xf4240",    // 1000000
		"miner":         "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5",
		"baseFeePerGas": "0x3b9aca00", // 1000000000
		"calls": []interface{}{
			map[string]interface{}{
				"returnData": "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000",
				"logs":       []interface{}{},
				"gasUsed":    "0x5208", // 21000
				"status":     "0x1",    // success
			},
		},
	}

	block := new(SimulateBlockResult)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, block)
	require.NoError(t, err)

	// Check number
	assert.Equal(t, uint64(24000000), block.Number)

	// Check hash
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", block.Hash)

	// Check timestamp
	assert.Equal(t, uint64(1705558944), block.Timestamp)

	// Check gasLimit (hex string)
	assert.Equal(t, "0x1c9c380", block.GasLimit)

	// Check gasUsed (hex string)
	assert.Equal(t, "0xf4240", block.GasUsed)

	// Check feeRecipient (checksum address)
	assert.Equal(t, "0x95222290DD7278Aa3Ddd389Cc1E1d165CC4BAfe5", block.FeeRecipient)

	// Check baseFeePerGas (hex string)
	assert.Equal(t, "0x3b9aca00", block.BaseFeePerGas)

	// Check calls
	require.Len(t, block.Calls, 1)
	assert.Equal(t, uint64(21000), block.Calls[0].GasUsed)
	assert.Equal(t, uint64(1), block.Calls[0].Status)
}

func TestSimulateCallResult_UnmarshalJSON(t *testing.T) {
	raw := map[string]interface{}{
		"returnData": "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000",
		"logs": []interface{}{
			map[string]interface{}{
				"address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				"topics":  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"},
				"data":    "0x0000000000000000000000000000000000000000000000000000000000000064",
			},
		},
		"gasUsed": "0x5208", // 21000
		"status":  "0x1",    // success
	}

	call := new(SimulateCallResult)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, call)
	require.NoError(t, err)

	// Check returnValue
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000", call.ReturnValue)

	// Check logs
	require.Len(t, call.Logs, 1)
	assert.Equal(t, "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", call.Logs[0].Address)

	// Check gasUsed
	assert.Equal(t, uint64(21000), call.GasUsed)

	// Check status
	assert.Equal(t, uint64(1), call.Status)

	// Check error (should be nil for successful call)
	assert.Nil(t, call.Error)
}

func TestSimulateCallResult_UnmarshalJSON_WithError(t *testing.T) {
	raw := map[string]interface{}{
		"returnData": "0x",
		"logs":       []interface{}{},
		"gasUsed":    "0x5208", // 21000
		"status":     "0x0",    // failure
		"error": map[string]interface{}{
			"code":    3, // execution reverted
			"message": "execution reverted",
			"data":    "0x08c379a00000000000000000000000000000000000000000000000000000000000000020",
		},
	}

	call := new(SimulateCallResult)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, call)
	require.NoError(t, err)

	// Check returnValue
	assert.Equal(t, "0x", call.ReturnValue)

	// Check logs (empty)
	assert.Empty(t, call.Logs)

	// Check gasUsed
	assert.Equal(t, uint64(21000), call.GasUsed)

	// Check status (failure)
	assert.Equal(t, uint64(0), call.Status)

	// Check error
	require.NotNil(t, call.Error)
	assert.Equal(t, 3, call.Error.Code)
	assert.Equal(t, "execution reverted", call.Error.Message)
	assert.Equal(t, "0x08c379a00000000000000000000000000000000000000000000000000000000000000020", call.Error.Data)
}

func TestCallError_UnmarshalJSON(t *testing.T) {
	raw := map[string]interface{}{
		"code":    3, // 3
		"message": "execution reverted: insufficient balance",
		"data":    "0x08c379a00000000000000000000000000000000000000000000000000000000000000020",
	}

	errorResult := new(CallError)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, errorResult)
	require.NoError(t, err)

	// Check code
	assert.Equal(t, 3, errorResult.Code)

	// Check message
	assert.Equal(t, "execution reverted: insufficient balance", errorResult.Message)

	// Check data
	assert.Equal(t, "0x08c379a00000000000000000000000000000000000000000000000000000000000000020", errorResult.Data)
}

func TestCallError_UnmarshalJSON_MinimalFields(t *testing.T) {
	// Test with only required fields (no data)
	raw := map[string]interface{}{
		"code":    1,
		"message": "out of gas",
	}

	errorResult := new(CallError)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, errorResult)
	require.NoError(t, err)

	// Check code
	assert.Equal(t, 1, errorResult.Code)

	// Check message
	assert.Equal(t, "out of gas", errorResult.Message)

	// Check data (should be empty)
	assert.Empty(t, errorResult.Data)
}

func TestSimulateBlockResult_UnmarshalJSON_MinimalFields(t *testing.T) {
	// Test with minimal required fields
	raw := map[string]interface{}{
		"number":    "0x0",
		"hash":      "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"timestamp": "0x0",
		"gasLimit":  "0x0",
		"gasUsed":   "0x0",
		"miner":     "0x0000000000000000000000000000000000000000",
		"calls":     []interface{}{},
	}

	block := new(SimulateBlockResult)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, block)
	require.NoError(t, err)

	// Check values
	assert.Equal(t, uint64(0), block.Number)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", block.Hash)
	assert.Equal(t, uint64(0), block.Timestamp)
	assert.Equal(t, "0x0", block.GasLimit)
	assert.Equal(t, "0x0", block.GasUsed)
	assert.Equal(t, "0x0000000000000000000000000000000000000000", block.FeeRecipient)
	assert.Empty(t, block.BaseFeePerGas)
	assert.Empty(t, block.Calls)
}
