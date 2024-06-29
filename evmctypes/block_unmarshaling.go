package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (b *block) UnmarshalJSON(input []byte) error {
	type block struct {
		Number                *string       `json:"number" validate:"-"`
		Hash                  *string       `json:"hash" validate:"required"`
		ParentHash            *string       `json:"parentHash" validate:"required"`
		Nonce                 *string       `json:"nonce" validate:"required"`
		MixHash               *string       `json:"mixHash" validate:"required"`
		Sha3Uncles            *string       `json:"sha3Uncles" validate:"required"`
		LogsBloom             *string       `json:"logsBloom" validate:"required"`
		StateRoot             *string       `json:"stateRoot" validate:"required"`
		Miner                 *string       `json:"miner" validate:"required"`
		Difficulty            *string       `json:"difficulty" validate:"required"`
		ExtraData             *string       `json:"extraData" validate:"required"`
		GasLimit              *string       `json:"gasLimit" validate:"required"`
		GasUsed               *string       `json:"gasUsed" validate:"required"`
		Timestamp             *string       `json:"timestamp" validate:"required"`
		TransactionsRoot      *string       `json:"transactionsRoot" validate:"required"`
		ReceiptsRoot          *string       `json:"receiptsRoot" validate:"required"`
		TotalDifficulty       *string       `json:"totalDifficulty" validate:"required"`
		Size                  *string       `json:"size"`
		Uncles                []string      `json:"uncles"`
		BaseFeePerGas         *string       `json:"baseFeePerGas,omitempty"`
		WithdrawalsRoot       *string       `json:"withdrawalsRoot,omitempty"`
		Withdrawals           []*Withdrawal `json:"withdrawals,omitempty"`
		BlobGasUsed           *string       `json:"blobGasUsed,omitempty"`
		ExcessBlobGas         *string       `json:"excessBlobGas,omitempty"`
		ParentBeaconBlockRoot *string       `json:"parentBeaconBlockRoot,omitempty"`
		L1BlockNumber         *string       `json:"l1BlockNumber,omitempty"`
		SendCount             *string       `json:"sendCount,omitempty"`
		SendRoot              *string       `json:"sendRoot,omitempty"`
	}
	var dec block
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Number != nil {
		number, err := hexutil.DecodeUint64(*dec.Number)
		if err != nil {
			return err
		}
		b.Number = number
	}
	if dec.Hash != nil {
		b.Hash = *dec.Hash
	}
	if dec.ParentHash != nil {
		b.ParentHash = *dec.ParentHash
	}
	if dec.Nonce != nil {
		b.Nonce = *dec.Nonce
	}
	if dec.MixHash != nil {
		b.MixHash = *dec.MixHash
	}
	if dec.Sha3Uncles != nil {
		b.Sha3Uncles = *dec.Sha3Uncles
	}
	if dec.LogsBloom != nil {
		b.LogsBloom = *dec.LogsBloom
	}
	if dec.StateRoot != nil {
		b.StateRoot = *dec.StateRoot
	}
	if dec.Miner != nil {
		b.Miner = *dec.Miner
	}
	if dec.Difficulty != nil {
		b.Difficulty = *dec.Difficulty
	}
	if dec.ExtraData != nil {
		b.ExtraData = *dec.ExtraData
	}
	if dec.GasLimit != nil {
		b.GasLimit = *dec.GasLimit
	}
	if dec.GasUsed != nil {
		b.GasUsed = *dec.GasUsed
	}
	if dec.Timestamp != nil {
		timestamp, err := hexutil.DecodeUint64(*dec.Timestamp)
		if err != nil {
			return err
		}
		b.Timestamp = timestamp
	}
	if dec.TransactionsRoot != nil {
		b.TransactionsRoot = *dec.TransactionsRoot
	}
	if dec.ReceiptsRoot != nil {
		b.ReceiptsRoot = *dec.ReceiptsRoot
	}
	if dec.TotalDifficulty != nil {
		b.TotalDifficulty = *dec.TotalDifficulty
	}
	if dec.Size != nil {
		b.Size = *dec.Size
	}
	if dec.Uncles != nil {
		b.Uncles = dec.Uncles
	}
	if dec.BaseFeePerGas != nil {
		b.BaseFeePerGas = dec.BaseFeePerGas
	}
	if dec.WithdrawalsRoot != nil {
		b.WithdrawalsRoot = dec.WithdrawalsRoot
	}
	if dec.Withdrawals != nil {
		b.Withdrawals = dec.Withdrawals
	}
	if dec.BlobGasUsed != nil {
		b.BlobGasUsed = dec.BlobGasUsed
	}
	if dec.ExcessBlobGas != nil {
		b.ExcessBlobGas = dec.ExcessBlobGas
	}
	if dec.ParentBeaconBlockRoot != nil {
		b.ParentBeaconBlockRoot = dec.ParentBeaconBlockRoot
	}
	if dec.L1BlockNumber != nil {
		l1BlockNumber, err := hexutil.DecodeUint64(*dec.L1BlockNumber)
		if err != nil {
			return err
		}
		b.L1BlockNumber = &l1BlockNumber
	}
	if dec.SendCount != nil {
		b.SendCount = dec.SendCount
	}
	if dec.SendRoot != nil {
		b.SendRoot = dec.SendRoot
	}
	return nil
}
