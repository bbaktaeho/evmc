package evmctypes

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, uint64(1), block.Withdrawals[0].Index)
	assert.Equal(t, uint64(2), block.Withdrawals[0].ValidatorIndex)
	assert.Equal(t, "0x3", block.Withdrawals[0].Address)
	assert.Equal(t, uint64(4), block.Withdrawals[0].Amount)
	assert.Nil(t, block.BaseFeePerGas)
	assert.Nil(t, block.L1BlockNumber)
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
