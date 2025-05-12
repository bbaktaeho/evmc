package kaiatypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (r *Receipt) UnmarshalJSON(input []byte) error {
	type Receipt struct {
		BlockHash          *string      `json:"blockHash" validate:"required"`
		BlockNumber        *string      `json:"blockNumber" validate:"-"`
		CodeFormat         *string      `json:"codeFormat,omitempty"`
		ContractAddress    *string      `json:"contractAddress,omitempty"`
		FeePayer           *string      `json:"feePayer,omitempty"`
		FeePayerSignatures []*Signature `json:"feePayerSignatures,omitempty"`
		FeeRatio           *string      `json:"feeRatio,omitempty"`
		From               *string      `json:"from" validate:"required"`
		Gas                *string      `json:"gas" validate:"required"`
		EffectiveGasPrice  *string      `json:"effectiveGasPrice" validate:"required"`
		GasPrice           *string      `json:"gasPrice" validate:"required"`
		GasUsed            *string      `json:"gasUsed" validate:"required"`
		HumanReadable      *bool        `json:"humanReadable,omitempty"`
		Key                *string      `json:"key,omitempty"`
		Logs               []*Log       `json:"logs" validate:"required"`
		LogsBloom          *string      `json:"logsBloom"`
		Nonce              *string      `json:"nonce" validate:"required"`
		SenderTxHash       *string      `json:"senderTxHash" validate:"required"`
		Signature          []*Signature `json:"signature" validate:"required"`
		Status             *string      `json:"status" validate:"required"`
		TxError            *string      `json:"txError,omitempty"`
		To                 *string      `json:"to" validate:"required"`
		TransactionHash    *string      `json:"transactionHash" validate:"required"`
		TransactionIndex   *string      `json:"transactionIndex" validate:"-"`
		Type               *string      `json:"type" validate:"required"`
		TypeInt            *uint64      `json:"typeInt" validate:"required"`
		Value              *string      `json:"value" validate:"required"`

		Input                string    `json:"input"`
		ChainID              *string   `json:"chainId,omitempty"`              // EIP-155
		AccessList           []*Access `json:"accessList,omitempty"`           // EIP-2930
		MaxFeePerGas         *string   `json:"maxFeePerGas,omitempty"`         // EIP-1559
		MaxPriorityFeePerGas *string   `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559
	}
	var dec Receipt
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.BlockHash != nil {
		r.BlockHash = *dec.BlockHash
	}
	if dec.BlockNumber != nil {
		blockNumber, err := hexutil.DecodeUint64(*dec.BlockNumber)
		if err != nil {
			return err
		}
		r.BlockNumber = blockNumber
	}
	if dec.CodeFormat != nil {
		r.CodeFormat = *dec.CodeFormat
	}
	if dec.ContractAddress != nil {
		r.ContractAddress = *dec.ContractAddress
	}
	if dec.FeePayer != nil {
		r.FeePayer = *dec.FeePayer
	}
	if dec.FeePayerSignatures != nil {
		r.FeePayerSignatures = dec.FeePayerSignatures
	}
	if dec.FeeRatio != nil {
		r.FeeRatio = *dec.FeeRatio
	}
	if dec.From != nil {
		r.From = *dec.From
	}
	if dec.Gas != nil {
		r.Gas = *dec.Gas
	}
	if dec.EffectiveGasPrice != nil {
		r.EffectiveGasPrice = *dec.EffectiveGasPrice
	}
	if dec.GasPrice != nil {
		r.GasPrice = *dec.GasPrice
	}
	if dec.GasUsed != nil {
		r.GasUsed = *dec.GasUsed
	}
	if dec.HumanReadable != nil {
		r.HumanReadable = *dec.HumanReadable
	}
	if dec.Key != nil {
		r.Key = *dec.Key
	}
	if dec.Logs != nil {
		r.Logs = dec.Logs
	}
	if dec.LogsBloom != nil {
		r.LogsBloom = *dec.LogsBloom
	}
	if dec.Nonce != nil {
		r.Nonce = *dec.Nonce
	}
	if dec.SenderTxHash != nil {
		r.SenderTxHash = *dec.SenderTxHash
	}
	if dec.Signature != nil {
		r.Signatures = dec.Signature
	}
	if dec.Status != nil {
		r.Status = *dec.Status
	}
	if dec.TxError != nil {
		r.TxError = *dec.TxError
	}
	if dec.To != nil {
		r.To = *dec.To
	}
	if dec.TransactionHash != nil {
		r.TransactionHash = *dec.TransactionHash
	}
	if dec.TransactionIndex != nil {
		transactionIndex, err := hexutil.DecodeUint64(*dec.TransactionIndex)
		if err != nil {
			return err
		}
		r.TransactionIndex = transactionIndex
	}
	if dec.Type != nil {
		r.Type = *dec.Type
	}
	if dec.TypeInt != nil {
		r.TypeInt = *dec.TypeInt
	}
	if dec.Value != nil {
		r.Value = *dec.Value
	}
	if dec.Input != "" {
		r.Input = dec.Input
	}
	if dec.ChainID != nil {
		r.ChainID = dec.ChainID
	}
	if dec.AccessList != nil {
		r.AccessList = dec.AccessList
	}
	if dec.MaxFeePerGas != nil {
		r.MaxFeePerGas = dec.MaxFeePerGas
	}
	if dec.MaxPriorityFeePerGas != nil {
		r.MaxPriorityFeePerGas = dec.MaxPriorityFeePerGas
	}
	return nil
}
