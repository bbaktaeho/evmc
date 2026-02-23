package evmctypes

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Block(t *testing.T) {
	raw := make(map[string]interface{})
	raw["number"] = "0x1"
	raw["hash"] = "0x2"
	raw["parentHash"] = "0x3"
	raw["nonce"] = "0x4"
	raw["mixHash"] = "0x5"
	raw["sha3Uncles"] = "0x6"
	raw["logsBloom"] = "0x7"
	raw["stateRoot"] = "0x8"
	raw["miner"] = "0x9"
	raw["difficulty"] = "0xa"
	raw["extraData"] = "0xb"
	raw["gasLimit"] = "0xc"
	raw["gasUsed"] = "0xd"
	raw["timestamp"] = "0xe"
	raw["transactionsRoot"] = "0xf"
	raw["receiptsRoot"] = "0x10"
	raw["totalDifficulty"] = "0x11"
	raw["size"] = "0x12"
	raw["uncles"] = []string{"0x13"}
	raw["baseFeePerGas"] = "0x14"
	raw["withdrawalsRoot"] = "0x15"
	raw["blobGasUsed"] = "0x16"
	raw["excessBlobGas"] = "0x17"
	raw["parentBeaconBlockRoot"] = "0x18"
	raw["l1BlockNumber"] = "0x19"
	raw["sendCount"] = "0x1a"
	raw["sendRoot"] = "0x1b"
	raw["transactions"] = []string{"0x10", "0x11", "0x12"}
	raw["withdrawals"] = []map[string]interface{}{
		{
			"index":          "0x1",
			"validatorIndex": "0x2",
			"address":        "0x3",
			"amount":         "0x4",
		},
	}

	block := new(Block)
	rawBytes, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(rawBytes, block); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, uint64(1), block.Number)
	assert.Equal(t, "0x2", block.Hash)
	assert.Equal(t, "0x3", block.ParentHash)
	assert.Equal(t, "0x4", block.Nonce)
	assert.Equal(t, "0x5", block.MixHash)
	assert.Equal(t, "0x6", block.Sha3Uncles)
	assert.Equal(t, "0x7", block.LogsBloom)
	assert.Equal(t, "0x8", block.StateRoot)
	assert.Equal(t, "0x9", block.Miner)
	assert.Equal(t, "0xa", block.Difficulty)
	assert.Equal(t, "0xb", block.ExtraData)
	assert.Equal(t, "0xc", block.GasLimit)
	assert.Equal(t, "0xd", block.GasUsed)
	assert.Equal(t, uint64(14), block.Timestamp)
	assert.Equal(t, "0xf", block.TransactionsRoot)
	assert.Equal(t, "0x10", block.ReceiptsRoot)
	assert.Equal(t, "0x11", block.TotalDifficulty)
	assert.Equal(t, "0x12", block.Size)
	assert.Equal(t, []string{"0x13"}, block.Uncles)
	assert.Equal(t, "0x14", *block.BaseFeePerGas)
	assert.Equal(t, "0x15", *block.WithdrawalsRoot)
	assert.Equal(t, "0x16", *block.BlobGasUsed)
	assert.Equal(t, "0x17", *block.ExcessBlobGas)
	assert.Equal(t, "0x18", *block.ParentBeaconBlockRoot)
	assert.Equal(t, uint64(25), *block.L1BlockNumber)
	assert.Equal(t, "0x1a", *block.SendCount)
	assert.Equal(t, "0x1b", *block.SendRoot)
	assert.Equal(t, []string{"0x10", "0x11", "0x12"}, block.Transactions)
	assert.Equal(t, uint64(1), block.Withdrawals[0].Index)
	assert.Equal(t, uint64(2), block.Withdrawals[0].ValidatorIndex)
	assert.Equal(t, "0x3", block.Withdrawals[0].Address)
	assert.Equal(t, uint64(4), block.Withdrawals[0].Amount)
}

func Test_block_NextBaseFee(t *testing.T) {
	raw := make(map[string]interface{})
	raw["baseFeePerGas"] = "0x1"
	block := new(Block)
	rawBytes, _ := json.Marshal(raw)
	json.Unmarshal(rawBytes, block)

	assert.Equal(t, decimal.NewFromInt(2), block.NextBaseFee())
}

func Test_Transaction(t *testing.T) {
	raw := make(map[string]interface{})
	raw["blockHash"] = "0x1"
	raw["blockNumber"] = "0x2"
	raw["from"] = "0x3"
	raw["gas"] = "0x4"
	raw["gasPrice"] = "0x5"
	raw["hash"] = "0x6"
	raw["input"] = "0x7"
	raw["nonce"] = "0x8"
	raw["to"] = "0x9"
	raw["transactionIndex"] = "0xa"
	raw["value"] = "0xb"
	raw["type"] = "0xc"
	raw["v"] = "0xd"
	raw["r"] = "0xe"

	tx := new(Transaction)
	rawBytes, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(rawBytes, tx); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "0x1", tx.BlockHash)
	assert.Equal(t, uint64(2), tx.BlockNumber)
	assert.Equal(t, "0x3", tx.From)
	assert.Equal(t, "0x4", tx.Gas)
	assert.Equal(t, "0x5", tx.GasPrice)
	assert.Equal(t, "0x6", tx.Hash)
	assert.Equal(t, "0x7", tx.Input)
	assert.Equal(t, "0x8", tx.Nonce)
	assert.Equal(t, "0x9", tx.To)
	assert.Equal(t, uint64(10), tx.TransactionIndex)
	assert.Equal(t, decimal.RequireFromString("11"), tx.Value)
	assert.Equal(t, "0xc", tx.Type)
	assert.Equal(t, "0xd", tx.V)
	assert.Equal(t, "0xe", tx.R)
}

func Test_Receipt(t *testing.T) {
	raw := make(map[string]interface{})
	raw["blockHash"] = "0x1"
	raw["blockNumber"] = "0x2"
	raw["transactionHash"] = "0x3"
	raw["transactionIndex"] = "0x4"
	raw["from"] = "0x5"
	raw["to"] = "0x6"
	raw["gasUsed"] = "0x7"
	raw["cumulativeGasUsed"] = "0x8"
	raw["contractAddress"] = "0x9"
	raw["logs"] = []map[string]interface{}{
		{
			"address":          "0xa",
			"topics":           []string{"0xb"},
			"data":             "0xc",
			"blockNumber":      "0xd",
			"transactionHash":  "0xe",
			"transactionIndex": "0xf",
			"blockHash":        "0x10",
			"logIndex":         "0x11",
			"removed":          false,
		},
	}
	raw["type"] = "0xd"
	raw["effectiveGasPrice"] = "0xe"
	raw["root"] = "0xf"
	raw["status"] = "0x10"
	raw["logsBloom"] = "0x11"

	receipt := new(Receipt)
	rawBytes, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(rawBytes, receipt); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "0x1", receipt.BlockHash)
	assert.Equal(t, uint64(2), receipt.BlockNumber)
	assert.Equal(t, "0x3", receipt.TransactionHash)
	assert.Equal(t, uint64(4), receipt.TransactionIndex)
	assert.Equal(t, "0x5", receipt.From)
	assert.Equal(t, "0x6", receipt.To)
	assert.Equal(t, "0x7", receipt.GasUsed)
	assert.Equal(t, "0x8", receipt.CumulativeGasUsed)
	assert.Equal(t, "0x9", *receipt.ContractAddress)
	assert.Equal(t, "0xa", receipt.Logs[0].Address)
	assert.Equal(t, []string{"0xb"}, receipt.Logs[0].Topics)
	assert.Equal(t, "0xc", receipt.Logs[0].Data)
	assert.Equal(t, uint64(13), receipt.Logs[0].BlockNumber)
	assert.Equal(t, "0xe", receipt.Logs[0].TransactionHash)
	assert.Equal(t, uint64(15), receipt.Logs[0].TransactionIndex)
	assert.Equal(t, "0x10", receipt.Logs[0].BlockHash)
	assert.Equal(t, uint64(17), receipt.Logs[0].LogIndex)
	assert.False(t, receipt.Logs[0].Removed)
	assert.Equal(t, "0xd", receipt.Type)
	assert.Equal(t, "0xe", receipt.EffectiveGasPrice)
	assert.Equal(t, "0xf", *receipt.Root)
	assert.Equal(t, "0x10", *receipt.Status)
	assert.Equal(t, "0x11", receipt.LogsBloom)
}

// mustMarshalUnmarshal은 raw 맵을 JSON으로 직렬화한 후 target에 역직렬화한다.
// 직렬화나 역직렬화 중 에러가 발생하면 즉시 테스트를 실패시킨다.
func mustMarshalUnmarshal[T any](t *testing.T, raw map[string]interface{}, target *T) {
	t.Helper()
	rawBytes, err := json.Marshal(raw)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rawBytes, target))
}

func Test_Log(t *testing.T) {
	raw := map[string]interface{}{
		"address":          "0xabc",
		"topics":           []string{"0x1", "0x2"},
		"data":             "0xdeadbeef",
		"blockNumber":      "0xa",
		"transactionHash":  "0xhash",
		"transactionIndex": "0x3",
		"blockHash":        "0xblockhash",
		"logIndex":         "0x5",
		"removed":          false,
		"blockTimestamp":   "0x64",
	}

	log := new(Log)
	mustMarshalUnmarshal(t, raw, log)

	assert.Equal(t, "0xabc", log.Address)
	assert.Equal(t, []string{"0x1", "0x2"}, log.Topics)
	assert.Equal(t, "0xdeadbeef", log.Data)
	assert.Equal(t, uint64(10), log.BlockNumber)
	assert.Equal(t, "0xhash", log.TransactionHash)
	assert.Equal(t, uint64(3), log.TransactionIndex)
	assert.Equal(t, "0xblockhash", log.BlockHash)
	assert.Equal(t, uint64(5), log.LogIndex)
	assert.False(t, log.Removed)
	assert.Equal(t, uint64(100), log.BlockTimestamp)
}

func Test_Log_Removed(t *testing.T) {
	raw := map[string]interface{}{
		"address":         "0xabc",
		"topics":          []string{},
		"data":            "0x",
		"blockNumber":     "0x1",
		"transactionHash": "0xhash",
		"blockHash":       "0xblockhash",
		"removed":         true,
	}

	log := new(Log)
	mustMarshalUnmarshal(t, raw, log)

	assert.True(t, log.Removed)
}

func Test_FeeHistory(t *testing.T) {
	raw := map[string]interface{}{
		"oldestBlock":   "0x64",
		"baseFeePerGas": []string{"0x1", "0x2", "0x3"},
		"gasUsedRatio":  []float64{0.1, 0.5, 0.9},
		"reward":        [][]string{{"0xa"}, {"0xb"}},
	}

	feeHistory := new(FeeHistory)
	mustMarshalUnmarshal(t, raw, feeHistory)

	assert.Equal(t, uint64(100), feeHistory.OldestBlock)
	assert.Equal(t, []string{"0x1", "0x2", "0x3"}, feeHistory.BaseFeePerGas)
	assert.Equal(t, []float64{0.1, 0.5, 0.9}, feeHistory.GasUsedRatio)
	assert.Equal(t, [][]string{{"0xa"}, {"0xb"}}, feeHistory.Reward)
}

func Test_FeeHistory_WithBlobGas(t *testing.T) {
	raw := map[string]interface{}{
		"oldestBlock":       "0x1",
		"baseFeePerBlobGas": []string{"0x100", "0x200"},
		"blobGasUsedRatio":  []float64{0.3, 0.7},
	}

	feeHistory := new(FeeHistory)
	mustMarshalUnmarshal(t, raw, feeHistory)

	assert.Equal(t, uint64(1), feeHistory.OldestBlock)
	assert.Equal(t, []string{"0x100", "0x200"}, feeHistory.BaseFeePerBlobGas)
	assert.Equal(t, []float64{0.3, 0.7}, feeHistory.BlobGasUsedRatio)
}

func Test_Syncing(t *testing.T) {
	raw := map[string]interface{}{
		"startingBlock":          "0x1",
		"currentBlock":           "0xa",
		"highestBlock":           "0x64",
		"syncedAccounts":         "0x10",
		"syncedAccountBytes":     "0x20",
		"syncedBytecodes":        "0x30",
		"syncedBytecodeBytes":    "0x40",
		"syncedStorage":          "0x50",
		"syncedStorageBytes":     "0x60",
		"healedTrienodes":        "0x70",
		"healedTrienodeBytes":    "0x80",
		"healedBytecodes":        "0x90",
		"healedBytecodeBytes":    "0xa0",
		"healingTrienodes":       "0xb0",
		"healingBytecode":        "0xc0",
		"txIndexFinishedBlocks":  "0xd0",
		"txIndexRemainingBlocks": "0xe0",
	}

	syncing := new(Syncing)
	mustMarshalUnmarshal(t, raw, syncing)

	assert.Equal(t, uint64(1), syncing.StartingBlock)
	assert.Equal(t, uint64(10), syncing.CurrentBlock)
	assert.Equal(t, uint64(100), syncing.HighestBlock)
	assert.Equal(t, uint64(16), syncing.SyncedAccounts)
	assert.Equal(t, uint64(32), syncing.SyncedAccountBytes)
	assert.Equal(t, uint64(48), syncing.SyncedBytecodes)
	assert.Equal(t, uint64(64), syncing.SyncedBytecodeBytes)
	assert.Equal(t, uint64(80), syncing.SyncedStorage)
	assert.Equal(t, uint64(96), syncing.SyncedStorageBytes)
	assert.Equal(t, uint64(112), syncing.HealedTrienodes)
	assert.Equal(t, uint64(128), syncing.HealedTrienodeBytes)
	assert.Equal(t, uint64(144), syncing.HealedBytecodes)
	assert.Equal(t, uint64(160), syncing.HealedBytecodeBytes)
	assert.Equal(t, uint64(176), syncing.HealingTrienodes)
	assert.Equal(t, uint64(192), syncing.HealingBytecode)
	assert.Equal(t, uint64(208), syncing.TxIndexFinishedBlocks)
	assert.Equal(t, uint64(224), syncing.TxIndexRemainingBlocks)
}

func Test_BlockIncTx(t *testing.T) {
	raw := map[string]interface{}{
		"number":           "0x1",
		"hash":             "0xhash",
		"parentHash":       "0xparent",
		"nonce":            "0x0",
		"mixHash":          "0x0",
		"sha3Uncles":       "0x0",
		"logsBloom":        "0x0",
		"stateRoot":        "0x0",
		"miner":            "0xminer",
		"difficulty":       "0x0",
		"extraData":        "0x",
		"gasLimit":         "0x1000",
		"gasUsed":          "0x500",
		"timestamp":        "0x1",
		"transactionsRoot": "0x0",
		"receiptsRoot":     "0x0",
		"transactions": []map[string]interface{}{
			{
				"blockHash":        "0xhash",
				"blockNumber":      "0x1",
				"from":             "0xfrom",
				"gas":              "0x5208",
				"gasPrice":         "0x1",
				"hash":             "0xtxhash",
				"input":            "0x",
				"nonce":            "0x0",
				"to":               "0xto",
				"transactionIndex": "0x0",
				"value":            "0x0",
				"type":             "0x0",
				"v":                "0x1",
				"r":                "0x1",
				"s":                "0x1",
			},
		},
	}

	blockIncTx := new(BlockIncTx)
	mustMarshalUnmarshal(t, raw, blockIncTx)

	assert.Equal(t, uint64(1), blockIncTx.Number)
	assert.Equal(t, "0xhash", blockIncTx.Hash)
	require.Len(t, blockIncTx.Transactions, 1)
	assert.Equal(t, "0xtxhash", blockIncTx.Transactions[0].Hash)
	assert.Equal(t, uint64(0), blockIncTx.Transactions[0].TransactionIndex)
}

func Test_Header(t *testing.T) {
	raw := map[string]interface{}{
		"number":           "0xa",
		"hash":             "0xheaderhash",
		"parentHash":       "0xparent",
		"nonce":            "0x0",
		"mixHash":          "0x0",
		"sha3Uncles":       "0x0",
		"logsBloom":        "0x0",
		"stateRoot":        "0x0",
		"miner":            "0xminer",
		"difficulty":       "0x0",
		"extraData":        "0x",
		"gasLimit":         "0x1000",
		"gasUsed":          "0x500",
		"timestamp":        "0x5f5e100",
		"transactionsRoot": "0x0",
		"receiptsRoot":     "0x0",
	}

	header := new(Header)
	mustMarshalUnmarshal(t, raw, header)

	assert.Equal(t, uint64(10), header.Number)
	assert.Equal(t, "0xheaderhash", header.Hash)
	assert.Equal(t, uint64(100000000), header.Timestamp)
}

func Test_CallFrame(t *testing.T) {
	raw := map[string]interface{}{
		"from":    "0xfrom",
		"gas":     "0x5208",
		"gasUsed": "0x5208",
		"to":      "0xto",
		"input":   "0x",
		"output":  "0xresult",
		"value":   "0xde0b6b3a7640000",
		"type":    "CALL",
		"calls": []map[string]interface{}{
			{
				"from":    "0xfrom2",
				"gas":     "0x1000",
				"gasUsed": "0x500",
				"input":   "0x",
				"type":    "STATICCALL",
			},
		},
	}

	callFrame := new(CallFrame)
	mustMarshalUnmarshal(t, raw, callFrame)

	assert.Equal(t, "0xfrom", callFrame.From)
	assert.Equal(t, "0x5208", callFrame.Gas)
	assert.Equal(t, "0x5208", callFrame.GasUsed)
	assert.Equal(t, "CALL", callFrame.Type)
	require.NotNil(t, callFrame.To)
	assert.Equal(t, "0xto", *callFrame.To)
	require.NotNil(t, callFrame.Output)
	assert.Equal(t, "0xresult", *callFrame.Output)
	require.NotNil(t, callFrame.Value)
	assert.Equal(t, decimal.RequireFromString("1000000000000000000"), *callFrame.Value)
	require.Len(t, callFrame.Calls, 1)
	assert.Equal(t, "STATICCALL", callFrame.Calls[0].Type)
}

func Test_CallFrame_WithError(t *testing.T) {
	raw := map[string]interface{}{
		"from":         "0xfrom",
		"gas":          "0x5208",
		"gasUsed":      "0x5208",
		"input":        "0x",
		"type":         "CALL",
		"error":        "execution reverted",
		"revertReason": "Ownable: caller is not the owner",
	}

	callFrame := new(CallFrame)
	mustMarshalUnmarshal(t, raw, callFrame)

	require.NotNil(t, callFrame.Error)
	assert.Equal(t, "execution reverted", *callFrame.Error)
	require.NotNil(t, callFrame.RevertReason)
	assert.Equal(t, "Ownable: caller is not the owner", *callFrame.RevertReason)
}

func Test_FlatCallFrame(t *testing.T) {
	callType := "call"
	from := "0xfrom"
	to := "0xto"
	gas := "0x5208"
	input := "0x"
	value := "0xde0b6b3a7640000"

	raw := map[string]interface{}{
		"action": map[string]interface{}{
			"callType": callType,
			"from":     from,
			"to":       to,
			"gas":      gas,
			"input":    input,
			"value":    value,
		},
		"blockHash":           "0xblockhash",
		"blockNumber":         uint64(100),
		"subtraces":           uint64(0),
		"traceAddress":        []uint64{},
		"transactionHash":     "0xtxhash",
		"transactionPosition": uint64(0),
		"type":                "call",
	}

	flatFrame := new(FlatCallFrame)
	mustMarshalUnmarshal(t, raw, flatFrame)

	assert.Equal(t, "0xblockhash", flatFrame.BlockHash)
	assert.Equal(t, uint64(100), flatFrame.BlockNumber)
	assert.Equal(t, "0xtxhash", flatFrame.TransactionHash)
	assert.Equal(t, "call", flatFrame.Type)
	require.NotNil(t, flatFrame.Action.CallType)
	assert.Equal(t, "call", *flatFrame.Action.CallType)
	require.NotNil(t, flatFrame.Action.Value)
	assert.Equal(t, decimal.RequireFromString("1000000000000000000"), *flatFrame.Action.Value)
}

func Test_BlockAndTag(t *testing.T) {
	t.Run("constants", func(t *testing.T) {
		assert.Equal(t, "pending", Pending.String())
		assert.Equal(t, "earliest", Earliest.String())
		assert.Equal(t, "latest", Latest.String())
		assert.Equal(t, "safe", Safe.String())
		assert.Equal(t, "finalized", Finalized.String())
	})

	t.Run("FormatNumber", func(t *testing.T) {
		tag := FormatNumber(100)
		assert.Equal(t, BlockAndTag("0x64"), tag)

		tag0 := FormatNumber(0)
		assert.Equal(t, BlockAndTag("0x0"), tag0)
	})

	t.Run("Uint64", func(t *testing.T) {
		tag := BlockAndTag("0x64")
		n, err := tag.Uint64()
		require.NoError(t, err)
		assert.Equal(t, uint64(100), n)
	})

	t.Run("ParseBlockAndTag with BlockAndTag", func(t *testing.T) {
		result := ParseBlockAndTag(Latest)
		assert.Equal(t, "latest", result)
	})

	t.Run("ParseBlockAndTag with string", func(t *testing.T) {
		result := ParseBlockAndTag("0x64")
		assert.Equal(t, "0x64", result)
	})

	t.Run("ParseBlockAndTag with non-string falls back to latest", func(t *testing.T) {
		result := ParseBlockAndTag(42)
		assert.Equal(t, "latest", result)
	})
}

func Test_Transaction_EIP1559(t *testing.T) {
	maxFeePerGas := "0x77359400"
	maxPriorityFeePerGas := "0x3b9aca00"
	chainID := "0x1"
	yParity := "0x1"

	raw := map[string]interface{}{
		"blockHash":            "0x1",
		"blockNumber":          "0x1",
		"from":                 "0xfrom",
		"gas":                  "0x5208",
		"gasPrice":             "0x0",
		"hash":                 "0xhash",
		"input":                "0x",
		"nonce":                "0x0",
		"to":                   "0xto",
		"transactionIndex":     "0x0",
		"value":                "0x0",
		"type":                 "0x2",
		"v":                    "0x0",
		"r":                    "0x1",
		"s":                    "0x1",
		"maxFeePerGas":         maxFeePerGas,
		"maxPriorityFeePerGas": maxPriorityFeePerGas,
		"chainId":              chainID,
		"yParity":              yParity,
		"accessList":           []interface{}{},
	}

	tx := new(Transaction)
	mustMarshalUnmarshal(t, raw, tx)

	assert.Equal(t, "0x2", tx.Type)
	require.NotNil(t, tx.MaxFeePerGas)
	assert.Equal(t, maxFeePerGas, *tx.MaxFeePerGas)
	require.NotNil(t, tx.MaxPriorityFeePerGas)
	assert.Equal(t, maxPriorityFeePerGas, *tx.MaxPriorityFeePerGas)
	require.NotNil(t, tx.ChainID)
	assert.Equal(t, chainID, *tx.ChainID)
	require.NotNil(t, tx.YParity)
	assert.Equal(t, yParity, *tx.YParity)
	assert.NotNil(t, tx.AccessList)
}

func Test_Receipt_L1Fields(t *testing.T) {
	raw := map[string]interface{}{
		"blockHash":         "0x1",
		"blockNumber":       "0x1",
		"transactionHash":   "0x2",
		"transactionIndex":  "0x0",
		"from":              "0xfrom",
		"to":                "0xto",
		"gasUsed":           "0x5208",
		"cumulativeGasUsed": "0x5208",
		"logs":              []interface{}{},
		"type":              "0x7e",
		"effectiveGasPrice": "0x0",
		"logsBloom":         "0x0",
		"l1BlockNumber":     "0x64",
		"l1GasPrice":        "0xabcd",
		"l1GasUsed":         "0x1234",
		"l1Fee":             "0x5678",
	}

	receipt := new(Receipt)
	mustMarshalUnmarshal(t, raw, receipt)

	require.NotNil(t, receipt.L1BlockNumber)
	assert.Equal(t, uint64(100), *receipt.L1BlockNumber)
	require.NotNil(t, receipt.L1GasPrice)
	assert.Equal(t, "0xabcd", *receipt.L1GasPrice)
	require.NotNil(t, receipt.L1GasUsed)
	assert.Equal(t, "0x1234", *receipt.L1GasUsed)
	require.NotNil(t, receipt.L1Fee)
	assert.Equal(t, "0x5678", *receipt.L1Fee)
}

func Test_Block_MilliTimestamp(t *testing.T) {
	raw := map[string]interface{}{
		"number":           "0x1",
		"hash":             "0xhash",
		"parentHash":       "0x0",
		"nonce":            "0x0",
		"mixHash":          "0x0",
		"sha3Uncles":       "0x0",
		"logsBloom":        "0x0",
		"stateRoot":        "0x0",
		"miner":            "0xminer",
		"difficulty":       "0x0",
		"extraData":        "0x",
		"gasLimit":         "0x1000",
		"gasUsed":          "0x0",
		"timestamp":        "0x64",
		"transactionsRoot": "0x0",
		"receiptsRoot":     "0x0",
		"transactions":     []string{},
		"milliTimestamp":   "0x17a4f8dadb8",
	}

	block := new(Block)
	mustMarshalUnmarshal(t, raw, block)

	require.NotNil(t, block.MilliTimestamp)
	assert.Equal(t, uint64(1624832323000), *block.MilliTimestamp)
}

func Test_Block_WithdrawalsAndUncles(t *testing.T) {
	raw := map[string]interface{}{
		"number":           "0x1",
		"hash":             "0xhash",
		"parentHash":       "0x0",
		"nonce":            "0x0",
		"mixHash":          "0x0",
		"sha3Uncles":       "0x0",
		"logsBloom":        "0x0",
		"stateRoot":        "0x0",
		"miner":            "0xminer",
		"difficulty":       "0x0",
		"extraData":        "0x",
		"gasLimit":         "0x1000",
		"gasUsed":          "0x0",
		"timestamp":        "0x1",
		"transactionsRoot": "0x0",
		"receiptsRoot":     "0x0",
		"transactions":     []string{"0xtx1"},
		"uncles":           []string{"0xuncle1"},
		"withdrawals": []map[string]interface{}{
			{
				"index":          "0x0",
				"validatorIndex": "0x10",
				"address":        "0xvalidator",
				"amount":         "0x3b9aca00",
			},
		},
	}

	block := new(Block)
	mustMarshalUnmarshal(t, raw, block)

	assert.Equal(t, []string{"0xuncle1"}, block.Uncles)
	require.Len(t, block.Withdrawals, 1)
	assert.Equal(t, uint64(0), block.Withdrawals[0].Index)
	assert.Equal(t, uint64(16), block.Withdrawals[0].ValidatorIndex)
	assert.Equal(t, "0xvalidator", block.Withdrawals[0].Address)
	assert.Equal(t, uint64(1000000000), block.Withdrawals[0].Amount)
}
