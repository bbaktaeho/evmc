package kaiatypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

func (t *Transaction) UnmarshalJSON(input []byte) error {
	type Transaction struct {
		BlockHash            *string      `json:"blockHash"`
		BlockNumber          *string      `json:"blockNumber"`
		CodeFormat           *string      `json:"codeFormat,omitempty"`
		FeePayer             *string      `json:"feePayer,omitempty"`
		FeePayerSignatures   []*Signature `json:"feePayerSignatures,omitempty"`
		FeeRatio             *string      `json:"feeRatio,omitempty"`
		From                 *string      `json:"from"`
		Gas                  *string      `json:"gas"`
		GasPrice             *string      `json:"gasPrice"`
		Hash                 *string      `json:"hash"`
		HumanReadable        *bool        `json:"humanReadable,omitempty"`
		Key                  *string      `json:"key,omitempty"`
		Input                *string      `json:"input"`
		Nonce                *string      `json:"nonce"`
		SenderTxHash         *string      `json:"senderTxHash"`
		Signatures           []*Signature `json:"signatures"`
		To                   *string      `json:"to"`
		TransactionIndex     *string      `json:"transactionIndex"`
		Type                 *string      `json:"type"`
		TypeInt              *uint64      `json:"typeInt"`
		Value                *string      `json:"value"`
		ChainID              *string      `json:"chainId,omitempty"`              // EIP-155
		AccessList           []*Access    `json:"accessList,omitempty"`           // EIP-2930
		MaxFeePerGas         *string      `json:"maxFeePerGas,omitempty"`         // EIP-1559
		MaxPriorityFeePerGas *string      `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559
	}
	var dec Transaction
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.BlockHash != nil {
		t.BlockHash = *dec.BlockHash
	}
	if dec.BlockNumber != nil {
		blockNumber, err := hexutil.DecodeUint64(*dec.BlockNumber)
		if err != nil {
			return err
		}
		t.BlockNumber = blockNumber
	}
	if dec.CodeFormat != nil {
		t.CodeFormat = *dec.CodeFormat
	}
	if dec.FeePayer != nil {
		t.FeePayer = *dec.FeePayer
	}
	if dec.FeePayerSignatures != nil {
		t.FeePayerSignatures = dec.FeePayerSignatures
	}
	if dec.FeeRatio != nil {
		t.FeeRatio = *dec.FeeRatio
	}
	if dec.From != nil {
		t.From = *dec.From
	}
	if dec.Gas != nil {
		t.Gas = *dec.Gas
	}
	if dec.GasPrice != nil {
		t.GasPrice = *dec.GasPrice
	}
	if dec.Hash != nil {
		t.Hash = *dec.Hash
	}
	if dec.HumanReadable != nil {
		t.HumanReadable = *dec.HumanReadable
	}
	if dec.Key != nil {
		t.Key = *dec.Key
	}
	if dec.Input != nil {
		t.Input = *dec.Input
	}
	if dec.Nonce != nil {
		t.Nonce = *dec.Nonce
	}
	if dec.SenderTxHash != nil {
		t.SenderTxHash = *dec.SenderTxHash
	}
	if dec.Signatures != nil {
		t.Signatures = dec.Signatures
	}
	if dec.To != nil {
		t.To = *dec.To
	}
	if dec.TransactionIndex != nil {
		transactionIndex, err := hexutil.DecodeUint64(*dec.TransactionIndex)
		if err != nil {
			return err
		}
		t.TransactionIndex = transactionIndex
	}
	if dec.Type != nil {
		t.Type = *dec.Type
	}
	if dec.TypeInt != nil {
		t.TypeInt = *dec.TypeInt
	}
	if dec.Value != nil {
		valueBig, err := hexutil.DecodeBig(*dec.Value)
		if err != nil {
			return err
		}
		t.Value = decimal.NewFromBigInt(valueBig, 0)
	}
	if dec.ChainID != nil {
		t.ChainID = dec.ChainID
	}
	if dec.AccessList != nil {
		t.AccessList = dec.AccessList
	}
	if dec.MaxFeePerGas != nil {
		t.MaxFeePerGas = dec.MaxFeePerGas
	}
	if dec.MaxPriorityFeePerGas != nil {
		t.MaxPriorityFeePerGas = dec.MaxPriorityFeePerGas
	}
	return nil
}
