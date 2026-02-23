package evmc

import (
	"encoding/json"
	"testing"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// blockJSON은 테스트용 블록 JSON (헥스 인코딩 필드 포함).
func blockJSON(number string, hash string, withTxHashes bool) map[string]interface{} {
	b := map[string]interface{}{
		"number":           number,
		"hash":             hash,
		"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce":            "0x0000000000000000",
		"mixHash":          "0x0000000000000000000000000000000000000000000000000000000000000000",
		"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"logsBloom":        "0x00000000000000000000000000000000",
		"stateRoot":        "0xabcd",
		"miner":            "0x1234567890123456789012345678901234567890",
		"difficulty":       "0x0",
		"extraData":        "0x",
		"gasLimit":         "0x1c9c380",
		"gasUsed":          "0x5208",
		"timestamp":        "0x64",
		"transactionsRoot": "0xabc",
		"receiptsRoot":     "0xdef",
		"totalDifficulty":  "0x0",
		"size":             "0x200",
		"uncles":           []string{},
		"baseFeePerGas":    "0x3b9aca00",
	}
	if withTxHashes {
		b["transactions"] = []string{"0xtx1", "0xtx2"}
	} else {
		b["transactions"] = []map[string]interface{}{
			{
				"blockHash":        hash,
				"blockNumber":      number,
				"from":             "0xfrom",
				"gas":              "0x5208",
				"gasPrice":         "0x3b9aca00",
				"hash":             "0xtx1",
				"input":            "0x",
				"nonce":            "0x0",
				"to":               "0xto",
				"transactionIndex": "0x0",
				"value":            "0xde0b6b3a7640000",
				"type":             "0x2",
				"v":                "0x1",
				"r":                "0x1",
				"s":                "0x1",
			},
		}
	}
	return b
}

func receiptJSON(txHash string, blockNumber string, blockHash string) map[string]interface{} {
	return map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       blockNumber,
		"transactionHash":   txHash,
		"transactionIndex":  "0x0",
		"from":              "0xfrom",
		"to":                "0xto",
		"gasUsed":           "0x5208",
		"cumulativeGasUsed": "0x5208",
		"logs":              []interface{}{},
		"type":              "0x2",
		"effectiveGasPrice": "0x3b9aca00",
		"status":            "0x1",
		"logsBloom":         "0x00000000",
	}
}

// testWithMock은 단일 메서드 핸들러를 등록하고 Evmc 클라이언트를 반환하는 헬퍼.
// 하나의 RPC 메서드만 모킹하는 간단한 테스트 케이스에 사용한다.
func testWithMock(t *testing.T, method string, handler func(params json.RawMessage) interface{}) *Evmc {
	t.Helper()
	mock := newMockRPCServer(t)
	mock.on(method, handler)
	return testEvmc(mock.url())
}

func Test_ethNamespace_mock_BlockNumber(t *testing.T) {
	client := testWithMock(t, "eth_blockNumber", func(params json.RawMessage) interface{} {
		return "0x64"
	})
	n, err := client.Eth().BlockNumber()
	require.NoError(t, err)
	// 0x64 = 100
	assert.Equal(t, uint64(100), n)
}

func Test_ethNamespace_mock_ChainID(t *testing.T) {
	client := testWithMock(t, "eth_chainId", func(params json.RawMessage) interface{} {
		return "0x1"
	})
	chainID, err := client.Eth().ChainID()
	require.NoError(t, err)
	assert.Equal(t, uint64(1), chainID)
}

func Test_ethNamespace_mock_GetBlockByNumber(t *testing.T) {
	client := testWithMock(t, "eth_getBlockByNumber", func(params json.RawMessage) interface{} {
		return blockJSON("0x64", "0xblockhash", true)
	})
	block, err := client.Eth().GetBlockByNumber(100)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), block.Number)
	assert.Equal(t, "0xblockhash", block.Hash)
	assert.Len(t, block.Transactions, 2)
	assert.Equal(t, "0x3b9aca00", *block.BaseFeePerGas)
}

func Test_ethNamespace_mock_GetBlockByNumber_tag(t *testing.T) {
	mock := newMockRPCServer(t)
	mock.on("eth_getBlockByNumber", func(params json.RawMessage) interface{} {
		var args []interface{}
		json.Unmarshal(params, &args)
		tag := args[0].(string)
		if tag == "latest" {
			return blockJSON("0x100", "0xlatesthash", true)
		}
		return nil
	})

	client := testEvmc(mock.url())
	block, err := client.Eth().GetBlockByTag(evmctypes.Latest)
	require.NoError(t, err)
	// 0x100 = 256
	assert.Equal(t, uint64(256), block.Number)
	assert.Equal(t, "0xlatesthash", block.Hash)
}

func Test_ethNamespace_mock_GetBlockIncTxByNumber(t *testing.T) {
	client := testWithMock(t, "eth_getBlockByNumber", func(params json.RawMessage) interface{} {
		return blockJSON("0x1", "0xhash1", false)
	})
	blockIncTx, err := client.Eth().GetBlockIncTxByNumber(1)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), blockIncTx.Number)
	require.Len(t, blockIncTx.Transactions, 1)
	assert.Equal(t, "0xtx1", blockIncTx.Transactions[0].Hash)
	// 0xde0b6b3a7640000 = 1000000000000000000 (1 ETH in wei)
	assert.True(t, decimal.RequireFromString("1000000000000000000").Equal(blockIncTx.Transactions[0].Value))
}

func Test_ethNamespace_mock_GetBlockByHash(t *testing.T) {
	client := testWithMock(t, "eth_getBlockByHash", func(params json.RawMessage) interface{} {
		return blockJSON("0xa", "0xtesthash", true)
	})
	block, err := client.Eth().GetBlockByHash("0xtesthash")
	require.NoError(t, err)
	// 0xa = 10
	assert.Equal(t, uint64(10), block.Number)
	assert.Equal(t, "0xtesthash", block.Hash)
}

func Test_ethNamespace_mock_GetTransactionByHash(t *testing.T) {
	client := testWithMock(t, "eth_getTransactionByHash", func(params json.RawMessage) interface{} {
		return map[string]interface{}{
			"blockHash":        "0xblockhash",
			"blockNumber":      "0x5",
			"from":             "0xfrom",
			"gas":              "0x5208",
			"gasPrice":         "0x3b9aca00",
			"hash":             "0xtxhash",
			"input":            "0x",
			"nonce":            "0x1",
			"to":               "0xto",
			"transactionIndex": "0x0",
			"value":            "0xde0b6b3a7640000",
			"type":             "0x0",
			"v":                "0x25",
			"r":                "0xabc",
			"s":                "0xdef",
		}
	})
	tx, err := client.Eth().GetTransactionByHash("0xtxhash")
	require.NoError(t, err)
	assert.Equal(t, "0xtxhash", tx.Hash)
	assert.Equal(t, "0xfrom", tx.From)
	assert.Equal(t, "0xto", tx.To)
	assert.Equal(t, uint64(5), tx.BlockNumber)
	// 0xde0b6b3a7640000 = 1000000000000000000 (1 ETH in wei)
	assert.True(t, decimal.RequireFromString("1000000000000000000").Equal(tx.Value))
}

func Test_ethNamespace_mock_GetTransactionReceipt(t *testing.T) {
	client := testWithMock(t, "eth_getTransactionReceipt", func(params json.RawMessage) interface{} {
		return receiptJSON("0xtxhash", "0x1", "0xblockhash")
	})
	receipt, err := client.Eth().GetTransactionReceipt("0xtxhash")
	require.NoError(t, err)
	assert.Equal(t, "0xtxhash", receipt.TransactionHash)
	assert.Equal(t, uint64(1), receipt.BlockNumber)
	assert.Equal(t, "0x1", *receipt.Status)
}

func Test_ethNamespace_mock_GetBalance(t *testing.T) {
	client := testWithMock(t, "eth_getBalance", func(params json.RawMessage) interface{} {
		// 0xde0b6b3a7640000 = 1000000000000000000 (1 ETH in wei)
		return "0xde0b6b3a7640000"
	})
	balance, err := client.Eth().GetBalance("0xaddr", evmctypes.Latest)
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("1000000000000000000"), balance)
}

func Test_ethNamespace_mock_GetBalance_Zero(t *testing.T) {
	client := testWithMock(t, "eth_getBalance", func(params json.RawMessage) interface{} {
		return "0x0"
	})
	balance, err := client.Eth().GetBalance("0xaddr", evmctypes.Latest)
	require.NoError(t, err)
	assert.True(t, balance.IsZero())
}

func Test_ethNamespace_mock_GasPrice(t *testing.T) {
	client := testWithMock(t, "eth_gasPrice", func(params json.RawMessage) interface{} {
		// 0x3b9aca00 = 1000000000 (1 Gwei)
		return "0x3b9aca00"
	})
	gasPrice, err := client.Eth().GasPrice()
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("1000000000"), gasPrice)
}

func Test_ethNamespace_mock_MaxPriorityFeePerGas(t *testing.T) {
	client := testWithMock(t, "eth_maxPriorityFeePerGas", func(params json.RawMessage) interface{} {
		// 0x77359400 = 2000000000 (2 Gwei)
		return "0x77359400"
	})
	fee, err := client.Eth().MaxPriorityFeePerGas()
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("2000000000"), fee)
}

func Test_ethNamespace_mock_GetTransactionCount(t *testing.T) {
	client := testWithMock(t, "eth_getTransactionCount", func(params json.RawMessage) interface{} {
		// 0x2a = 42
		return "0x2a"
	})
	nonce, err := client.Eth().GetTransactionCount("0xaddr", evmctypes.Latest)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), nonce)
}

func Test_ethNamespace_mock_GetCode(t *testing.T) {
	client := testWithMock(t, "eth_getCode", func(params json.RawMessage) interface{} {
		return "0x60806040"
	})
	code, err := client.Eth().GetCode("0xcontract", evmctypes.Latest)
	require.NoError(t, err)
	assert.Equal(t, "0x60806040", code)
}

func Test_ethNamespace_mock_GetStorageAt(t *testing.T) {
	client := testWithMock(t, "eth_getStorageAt", func(params json.RawMessage) interface{} {
		return "0x0000000000000000000000000000000000000000000000000000000000000001"
	})
	storage, err := client.Eth().GetStorageAt("0xaddr", "0x0", evmctypes.Latest)
	require.NoError(t, err)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", storage)
}

func Test_ethNamespace_mock_GetLogs(t *testing.T) {
	client := testWithMock(t, "eth_getLogs", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			{
				"address":          "0xcontract",
				"topics":           []string{"0xtopic1"},
				"data":             "0xdata",
				"blockNumber":      "0x1",
				"transactionHash":  "0xtxhash",
				"transactionIndex": "0x0",
				"blockHash":        "0xblockhash",
				"logIndex":         "0x0",
				"removed":          false,
				// 0x64 = 100 (unix timestamp)
				"blockTimestamp": "0x64",
			},
		}
	})
	fromBlock := uint64(1)
	toBlock := uint64(10)
	addr := "0xcontract"
	logs, err := client.Eth().GetLogs(&evmctypes.LogFilter{
		FromBlock: &fromBlock,
		ToBlock:   &toBlock,
		Address:   &addr,
	})
	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Equal(t, "0xcontract", logs[0].Address)
	assert.Equal(t, []string{"0xtopic1"}, logs[0].Topics)
	assert.Equal(t, uint64(1), logs[0].BlockNumber)
	assert.Equal(t, uint64(100), logs[0].BlockTimestamp)
}

func Test_ethNamespace_mock_GetBlockReceipts(t *testing.T) {
	client := testWithMock(t, "eth_getBlockReceipts", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			receiptJSON("0xtx1", "0x1", "0xblockhash"),
			receiptJSON("0xtx2", "0x1", "0xblockhash"),
		}
	})
	receipts, err := client.Eth().GetBlockReceipts(1)
	require.NoError(t, err)
	require.Len(t, receipts, 2)
	assert.Equal(t, "0xtx1", receipts[0].TransactionHash)
	assert.Equal(t, "0xtx2", receipts[1].TransactionHash)
}

func Test_ethNamespace_mock_BlockTransactionCountByNumber(t *testing.T) {
	client := testWithMock(t, "eth_getBlockTransactionCountByNumber", func(params json.RawMessage) interface{} {
		// 0xa = 10
		return "0xa"
	})
	count, err := client.Eth().BlockTransactionCountByNumber(1)
	require.NoError(t, err)
	assert.Equal(t, uint64(10), count)
}

func Test_ethNamespace_mock_BlockTransactionCountByHash(t *testing.T) {
	client := testWithMock(t, "eth_getBlockTransactionCountByHash", func(params json.RawMessage) interface{} {
		return "0x5"
	})
	count, err := client.Eth().BlockTransactionCountByHash("0xhash")
	require.NoError(t, err)
	assert.Equal(t, uint64(5), count)
}

func Test_ethNamespace_mock_FeeHistory(t *testing.T) {
	client := testWithMock(t, "eth_feeHistory", func(params json.RawMessage) interface{} {
		return map[string]interface{}{
			// 0x64 = 100
			"oldestBlock":   "0x64",
			"baseFeePerGas": []string{"0x3b9aca00", "0x77359400"},
			"gasUsedRatio":  []float64{0.5, 0.8},
		}
	})
	feeHistory, err := client.Eth().FeeHistory(2, evmctypes.Latest, nil)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), feeHistory.OldestBlock)
	assert.Len(t, feeHistory.BaseFeePerGas, 2)
	assert.Len(t, feeHistory.GasUsedRatio, 2)
}

func Test_ethNamespace_mock_EstimateGas(t *testing.T) {
	client := testWithMock(t, "eth_estimateGas", func(params json.RawMessage) interface{} {
		// 0x5208 = 21000 (standard ETH transfer gas)
		return "0x5208"
	})
	gas, err := client.Eth().EstimateGas(&Tx{
		From:  ZeroAddress,
		To:    ZeroAddress,
		Data:  "0x",
		Value: decimal.Zero,
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(21000), gas)
}

func Test_ethNamespace_mock_Call(t *testing.T) {
	client := testWithMock(t, "eth_call", func(params json.RawMessage) interface{} {
		return "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000"
	})
	result, err := client.Eth().Call(&Tx{
		From: ZeroAddress,
		To:   "0xcontract",
		Data: "0x18160ddd",
	}, evmctypes.Latest)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func Test_ethNamespace_mock_GetBlockRange(t *testing.T) {
	mock := newMockRPCServer(t)
	callCount := 0
	mock.on("eth_getBlockByNumber", func(params json.RawMessage) interface{} {
		callCount++
		var args []interface{}
		json.Unmarshal(params, &args)
		hexNum := args[0].(string)
		return blockJSON(hexNum, "0xhash"+hexNum, true)
	})

	client := testEvmc(mock.url())
	blocks, err := client.Eth().GetBlockRange(1, 3)
	require.NoError(t, err)
	assert.Len(t, blocks, 3)
	assert.Equal(t, uint64(1), blocks[0].Number)
	assert.Equal(t, uint64(2), blocks[1].Number)
	assert.Equal(t, uint64(3), blocks[2].Number)
}

func Test_ethNamespace_mock_GetBlockRange_InvalidRange(t *testing.T) {
	mock := newMockRPCServer(t)
	client := testEvmc(mock.url())
	_, err := client.Eth().GetBlockRange(10, 5)
	require.ErrorIs(t, err, ErrInvalidRange)
}

func Test_ethNamespace_mock_Syncing_false(t *testing.T) {
	client := testWithMock(t, "eth_syncing", func(params json.RawMessage) interface{} {
		return false
	})
	syncing, _, err := client.Eth().Syncing()
	require.NoError(t, err)
	assert.False(t, syncing)
}

func Test_ethNamespace_mock_GetLogsByBlockNumber(t *testing.T) {
	client := testWithMock(t, "eth_getLogs", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			{
				"address": "0xcontract",
				"topics":  []string{},
				"data":    "0x",
				// 0xa = 10
				"blockNumber":      "0xa",
				"transactionHash":  "0xtxhash",
				"transactionIndex": "0x0",
				"blockHash":        "0xblockhash",
				"logIndex":         "0x0",
				"removed":          false,
			},
		}
	})
	logs, err := client.Eth().GetLogsByBlockNumber(10)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Equal(t, uint64(10), logs[0].BlockNumber)
}

func Test_ethNamespace_mock_GetLogsByBlockHash(t *testing.T) {
	client := testWithMock(t, "eth_getLogs", func(params json.RawMessage) interface{} {
		return []map[string]interface{}{
			{
				"address": "0xcontract",
				"topics":  []string{},
				"data":    "0x",
				// 0xa = 10
				"blockNumber":      "0xa",
				"transactionHash":  "0xtxhash",
				"transactionIndex": "0x0",
				"blockHash":        "0xtesthash",
				"logIndex":         "0x0",
				"removed":          false,
			},
		}
	})
	logs, err := client.Eth().GetLogsByBlockHash("0xtesthash")
	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Equal(t, "0xtesthash", logs[0].BlockHash)
}

func Test_ethNamespace_mock_BlobBaseFee(t *testing.T) {
	client := testWithMock(t, "eth_blobBaseFee", func(params json.RawMessage) interface{} {
		// 0x3b9aca00 = 1000000000 (1 Gwei)
		return "0x3b9aca00"
	})
	fee, err := client.Eth().BlobBaseFee()
	require.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("1000000000"), fee)
}

func Test_ethNamespace_mock_SendRawTransaction(t *testing.T) {
	client := testWithMock(t, "eth_sendRawTransaction", func(params json.RawMessage) interface{} {
		return "0xsentTxHash"
	})
	txHash, err := client.Eth().SendRawTransaction("0xrawTx")
	require.NoError(t, err)
	assert.Equal(t, "0xsentTxHash", txHash)
}
