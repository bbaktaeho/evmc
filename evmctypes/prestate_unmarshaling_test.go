package evmctypes

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrestateAccount_UnmarshalJSON(t *testing.T) {
	raw := map[string]interface{}{
		"balance":  "0xde0b6b3a7640000", // 1 ETH in wei
		"code":     "0x60806040",
		"codeHash": "0x1234567890abcdef",
		"nonce":    uint64(10), // geth returns nonce as integer
		"storage": map[string]string{
			"0x0": "0x1",
			"0x1": "0x2",
		},
	}

	account := new(PrestateAccount)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, account)
	require.NoError(t, err)

	// Check balance (1 ETH = 1000000000000000000 wei)
	expectedBalance := decimal.RequireFromString("1000000000000000000")
	assert.NotNil(t, account.Balance)
	assert.True(t, account.Balance.Equal(expectedBalance), "Balance mismatch: expected %s, got %s", expectedBalance, account.Balance)

	// Check code
	assert.Equal(t, "0x60806040", account.Code)

	// Check codeHash
	assert.Equal(t, "0x1234567890abcdef", account.CodeHash)

	// Check nonce
	assert.Equal(t, uint64(10), account.Nonce)

	// Check storage
	assert.Equal(t, 2, len(account.Storage))
	assert.Equal(t, "0x1", account.Storage["0x0"])
	assert.Equal(t, "0x2", account.Storage["0x1"])
}

func TestPrestateAccount_UnmarshalJSON_MinimalFields(t *testing.T) {
	// Test with only required/minimal fields
	raw := map[string]interface{}{
		"nonce": uint64(0),
	}

	account := new(PrestateAccount)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, account)
	require.NoError(t, err)

	assert.Nil(t, account.Balance)
	assert.Equal(t, "", account.Code)
	assert.Equal(t, "", account.CodeHash)
	assert.Equal(t, uint64(0), account.Nonce)
	assert.Nil(t, account.Storage)
}

func TestPrestateTracer_UnmarshalJSON_StandardMode(t *testing.T) {
	// Standard prestate mode (not diff mode)
	raw := map[string]interface{}{
		"txHash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"result": map[string]interface{}{
			"0x0000000000000000000000000000000000000001": map[string]interface{}{
				"balance": "0xde0b6b3a7640000",
				"nonce":   uint64(1),
			},
			"0x0000000000000000000000000000000000000002": map[string]interface{}{
				"balance": "0x1bc16d674ec80000",
				"nonce":   uint64(5),
				"code":    "0x60806040",
				"storage": map[string]string{
					"0x0": "0x1",
				},
			},
		},
	}

	tracer := new(PrestateTracer)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, tracer)
	require.NoError(t, err)

	// Check txHash
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", tracer.TxHash)

	// Check error is nil
	assert.Nil(t, tracer.Error)

	// Parse as standard prestate frame
	frame, err := tracer.ParseFrames()
	require.NoError(t, err)
	require.NotNil(t, frame)

	// Check accounts
	assert.Equal(t, 2, len(frame))

	// Check first account
	account1 := (frame)["0x0000000000000000000000000000000000000001"]
	require.NotNil(t, account1)
	assert.Equal(t, uint64(1), account1.Nonce)
	expectedBalance1 := decimal.RequireFromString("1000000000000000000")
	assert.True(t, account1.Balance.Equal(expectedBalance1))

	// Check second account
	account2 := (frame)["0x0000000000000000000000000000000000000002"]
	require.NotNil(t, account2)
	assert.Equal(t, uint64(5), account2.Nonce)
	expectedBalance2 := decimal.RequireFromString("2000000000000000000")
	assert.True(t, account2.Balance.Equal(expectedBalance2))
	assert.Equal(t, "0x60806040", account2.Code)
	assert.Equal(t, 1, len(account2.Storage))
	assert.Equal(t, "0x1", account2.Storage["0x0"])
}

func TestPrestateTracer_UnmarshalJSON_DiffMode(t *testing.T) {
	// Diff mode prestate
	raw := map[string]interface{}{
		"txHash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"result": map[string]interface{}{
			"pre": map[string]interface{}{
				"0x0000000000000000000000000000000000000001": map[string]interface{}{
					"balance": "0xde0b6b3a7640000",
					"nonce":   uint64(1),
				},
			},
			"post": map[string]interface{}{
				"0x0000000000000000000000000000000000000001": map[string]interface{}{
					"balance": "0xd2f13f7789f0000", // Changed balance
					"nonce":   uint64(2),           // Incremented nonce
				},
				"0x0000000000000000000000000000000000000002": map[string]interface{}{
					"balance": "0x16345785d8a0000", // New account
					"nonce":   uint64(1),
				},
			},
		},
	}

	tracer := new(PrestateTracer)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, tracer)
	require.NoError(t, err)

	// Check txHash
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", tracer.TxHash)

	// Parse as diff frame
	diffFrame, err := tracer.ParseDiffFrames()
	require.NoError(t, err)
	require.NotNil(t, diffFrame)

	// Check pre-state
	require.NotNil(t, diffFrame.Pre)
	assert.Equal(t, 1, len(diffFrame.Pre))

	preAccount := diffFrame.Pre["0x0000000000000000000000000000000000000001"]
	require.NotNil(t, preAccount)
	assert.Equal(t, uint64(1), preAccount.Nonce)
	expectedPreBalance := decimal.RequireFromString("1000000000000000000")
	assert.True(t, preAccount.Balance.Equal(expectedPreBalance))

	// Check post-state
	require.NotNil(t, diffFrame.Post)
	assert.Equal(t, 2, len(diffFrame.Post))

	postAccount1 := diffFrame.Post["0x0000000000000000000000000000000000000001"]
	require.NotNil(t, postAccount1)
	assert.Equal(t, uint64(2), postAccount1.Nonce)
	expectedPostBalance1 := decimal.RequireFromString("950000000000000000")
	assert.True(t, postAccount1.Balance.Equal(expectedPostBalance1))

	postAccount2 := diffFrame.Post["0x0000000000000000000000000000000000000002"]
	require.NotNil(t, postAccount2)
	assert.Equal(t, uint64(1), postAccount2.Nonce)
	expectedPostBalance2 := decimal.RequireFromString("100000000000000000")
	assert.True(t, postAccount2.Balance.Equal(expectedPostBalance2))
}

func TestPrestateTracer_UnmarshalJSON_WithError(t *testing.T) {
	t.Run("with null result", func(t *testing.T) {
		raw := map[string]interface{}{
			"txHash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			"error":  "execution reverted",
			"result": nil,
		}

		tracer := new(PrestateTracer)
		rawBytes, err := json.Marshal(raw)
		require.NoError(t, err)

		err = json.Unmarshal(rawBytes, tracer)
		require.NoError(t, err)

		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", tracer.TxHash)
		assert.Equal(t, "execution reverted", tracer.Error)
		// JSON null becomes []byte("null"), not nil
		assert.Equal(t, json.RawMessage("null"), tracer.Result)
	})

	t.Run("without result field", func(t *testing.T) {
		jsonStr := `{
			"txHash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			"error": "execution reverted"
		}`

		tracer := new(PrestateTracer)
		err := json.Unmarshal([]byte(jsonStr), tracer)
		require.NoError(t, err)

		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", tracer.TxHash)
		assert.Equal(t, "execution reverted", tracer.Error)
		// When field is omitted, json.RawMessage is nil
		assert.Nil(t, tracer.Result)
	})

	t.Run("error as map", func(t *testing.T) {
		// Test with error as map (erigon format)
		jsonStr := `{
			"txHash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			"error": {
				"message": "execution reverted",
				"code": -32000
			}
		}`

		tracer := new(PrestateTracer)
		err := json.Unmarshal([]byte(jsonStr), tracer)
		require.NoError(t, err)

		assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", tracer.TxHash)
		assert.NotNil(t, tracer.Error)
		// Error is interface{}, so it can be a map
		errorMap, ok := tracer.Error.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "execution reverted", errorMap["message"])
	})
}

func TestPrestateTracer_ParseResult(t *testing.T) {
	t.Run("empty result", func(t *testing.T) {
		tracer := &PrestateTracer{}
		frame, err := tracer.ParseFrames()
		require.NoError(t, err)
		assert.Nil(t, frame)

		diffFrame, err := tracer.ParseDiffFrames()
		require.NoError(t, err)
		assert.Nil(t, diffFrame)
	})

	t.Run("null result", func(t *testing.T) {
		tracer := &PrestateTracer{
			Result: json.RawMessage("null"),
		}

		frame, err := tracer.ParseFrames()
		require.NoError(t, err)
		assert.Nil(t, frame)

		diffFrame, err := tracer.ParseDiffFrames()
		require.NoError(t, err)
		assert.Nil(t, diffFrame)
	})

	t.Run("invalid json", func(t *testing.T) {
		tracer := &PrestateTracer{
			Result: json.RawMessage(`{invalid json}`),
		}

		_, err := tracer.ParseFrames()
		assert.Error(t, err)

		_, err = tracer.ParseDiffFrames()
		assert.Error(t, err)
	})
}

func TestPrestateAccount_EmptyStorage(t *testing.T) {
	raw := map[string]interface{}{
		"balance": "0x0",
		"nonce":   uint64(0),
		"storage": map[string]string{},
	}

	account := new(PrestateAccount)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, account)
	require.NoError(t, err)

	assert.NotNil(t, account.Storage)
	assert.Equal(t, 0, len(account.Storage))
}

func TestPrestateAccount_LargeBalance(t *testing.T) {
	raw := map[string]interface{}{
		"balance": "0x152d02c7e14af6800000", // 100000 ETH
		"nonce":   uint64(0),
	}

	account := new(PrestateAccount)
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)

	err = json.Unmarshal(rawBytes, account)
	require.NoError(t, err)

	expectedBalance := decimal.RequireFromString("100000000000000000000000")
	assert.True(t, account.Balance.Equal(expectedBalance))
}

func TestPrestateTracer_UnmarshalJSON_WithResult(t *testing.T) {
	t.Run("result as prestate map", func(t *testing.T) {
		// PrestateTracer with result field containing account map
		raw := map[string]interface{}{
			"txHash": "0xabc",
			"result": map[string]interface{}{
				"0x0000000000000000000000000000000000000001": map[string]interface{}{
					"balance": "0xde0b6b3a7640000",
					"nonce":   uint64(1),
				},
				"0x0000000000000000000000000000000000000002": map[string]interface{}{
					"balance": "0x1bc16d674ec80000",
					"nonce":   uint64(5),
					"code":    "0x60806040",
				},
			},
		}

		tracer := new(PrestateTracer)
		rawBytes, err := json.Marshal(raw)
		require.NoError(t, err)
		err = json.Unmarshal(rawBytes, tracer)
		require.NoError(t, err)

		assert.Equal(t, "0xabc", tracer.TxHash)
		assert.NotNil(t, tracer.Result)

		frame, err := tracer.ParseFrames()
		require.NoError(t, err)
		require.NotNil(t, frame)

		assert.Equal(t, 2, len(frame))
		account1 := frame["0x0000000000000000000000000000000000000001"]
		require.NotNil(t, account1)
		assert.Equal(t, uint64(1), account1.Nonce)
	})

	t.Run("result as diff mode map", func(t *testing.T) {
		raw := map[string]interface{}{
			"txHash": "0xdef",
			"result": map[string]interface{}{
				"pre": map[string]interface{}{
					"0x0000000000000000000000000000000000000001": map[string]interface{}{
						"balance": "0xde0b6b3a7640000",
						"nonce":   uint64(1),
					},
				},
				"post": map[string]interface{}{
					"0x0000000000000000000000000000000000000001": map[string]interface{}{
						"balance": "0xd2f13f7789f0000",
						"nonce":   uint64(2),
					},
				},
			},
		}

		tracer := new(PrestateTracer)
		rawBytes, err := json.Marshal(raw)
		require.NoError(t, err)
		err = json.Unmarshal(rawBytes, tracer)
		require.NoError(t, err)

		assert.Equal(t, "0xdef", tracer.TxHash)
		assert.NotNil(t, tracer.Result)

		diffFrame, err := tracer.ParseDiffFrames()
		require.NoError(t, err)
		require.NotNil(t, diffFrame)

		assert.NotNil(t, diffFrame.Pre)
		assert.NotNil(t, diffFrame.Post)
		assert.Equal(t, 1, len(diffFrame.Pre))
		assert.Equal(t, 1, len(diffFrame.Post))
	})
}
