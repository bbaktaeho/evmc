package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (r *Receipt) UnmarshalJSON(input []byte) error {
	type Receipt struct {
		BlockHash             *string `json:"blockHash" validate:"required"`
		BlockNumber           *string `json:"blockNumber" validate:"-"`
		TransactionHash       *string `json:"transactionHash" validate:"required"`
		TransactionIndex      *string `json:"transactionIndex" validate:"-"`
		From                  *string `json:"from" validate:"required"`
		To                    *string `json:"to" validate:"required"`
		GasUsed               *string `json:"gasUsed" validate:"required"`
		CumulativeGasUsed     *string `json:"cumulativeGasUsed" validate:"required"`
		ContractAddress       *string `json:"contractAddress,omitempty"`
		Logs                  []*Log  `json:"logs" validate:"required"`
		Type                  *string `json:"type" validate:"required"`
		EffectiveGasPrice     *string `json:"effectiveGasPrice" validate:"required"`
		Root                  *string `json:"root,omitempty"`
		Status                *string `json:"status,omitempty"`
		LogsBloom             *string `json:"logsBloom"`
		BlobGasPrice          *string `json:"blobGasPrice,omitempty"`
		BlobGasUsed           *string `json:"blobGasUsed,omitempty"`
		GasUsedForL1          *string `json:"gasUsedForL1,omitempty"`
		L1BlockNumber         *string `json:"l1BlockNumber,omitempty"`
		L1GasPrice            *string `json:"l1GasPrice,omitempty"`
		L1GasUsed             *string `json:"l1GasUsed,omitempty"`
		L1FeeScalar           *string `json:"l1FeeScalar,omitempty"`
		L1Fee                 *string `json:"l1Fee,omitempty"`
		DepositNonce          *string `json:"depositNonce,omitempty"`
		DepositReceiptVersion *string `json:"depositReceiptVersion,omitempty"`
		L1BlobBaseFee         *string `json:"l1BlobBaseFee,omitempty"`
		L1BaseFeeScalar       *string `json:"l1BaseFeeScalar,omitempty"`
		L1BlobBaseFeeScalar   *string `json:"l1BlobBaseFeeScalar,omitempty"`
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
	if dec.From != nil {
		r.From = *dec.From
	}
	if dec.To != nil {
		r.To = *dec.To
	}
	if dec.GasUsed != nil {
		r.GasUsed = *dec.GasUsed
	}
	if dec.CumulativeGasUsed != nil {
		r.CumulativeGasUsed = *dec.CumulativeGasUsed
	}
	if dec.ContractAddress != nil {
		r.ContractAddress = dec.ContractAddress
	}
	if dec.Logs != nil {
		r.Logs = dec.Logs
	}
	if dec.Type != nil {
		r.Type = *dec.Type
	}
	if dec.EffectiveGasPrice != nil {
		r.EffectiveGasPrice = *dec.EffectiveGasPrice
	}
	if dec.Root != nil {
		r.Root = dec.Root
	}
	if dec.Status != nil {
		r.Status = dec.Status
	}
	if dec.LogsBloom != nil {
		r.LogsBloom = *dec.LogsBloom
	}
	if dec.BlobGasPrice != nil {
		r.BlobGasPrice = dec.BlobGasPrice
	}
	if dec.BlobGasUsed != nil {
		r.BlobGasUsed = dec.BlobGasUsed
	}
	if dec.GasUsedForL1 != nil {
		r.GasUsedForL1 = dec.GasUsedForL1
	}
	if dec.L1BlockNumber != nil {
		l1BlockNumber, err := hexutil.DecodeUint64(*dec.L1BlockNumber)
		if err != nil {
			return err
		}
		r.L1BlockNumber = &l1BlockNumber
	}
	if dec.L1GasPrice != nil {
		r.L1GasPrice = dec.L1GasPrice
	}
	if dec.L1GasUsed != nil {
		r.L1GasUsed = dec.L1GasUsed
	}
	if dec.L1FeeScalar != nil {
		r.L1FeeScalar = dec.L1FeeScalar
	}
	if dec.L1Fee != nil {
		r.L1Fee = dec.L1Fee
	}
	if dec.DepositNonce != nil {
		r.DepositNonce = dec.DepositNonce
	}
	if dec.DepositReceiptVersion != nil {
		r.DepositReceiptVersion = dec.DepositReceiptVersion
	}
	if dec.L1BlobBaseFee != nil {
		r.L1BlobBaseFee = dec.L1BlobBaseFee
	}
	if dec.L1BaseFeeScalar != nil {
		r.L1BaseFeeScalar = dec.L1BaseFeeScalar
	}
	if dec.L1BlobBaseFeeScalar != nil {
		r.L1BlobBaseFeeScalar = dec.L1BlobBaseFeeScalar
	}
	return nil
}
