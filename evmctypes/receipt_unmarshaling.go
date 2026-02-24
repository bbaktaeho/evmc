package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (r *Receipt) UnmarshalJSON(input []byte) error {
	type receiptWire struct {
		// hex→uint64 후처리 대상
		BlockNumber      *string `json:"blockNumber"`
		TransactionIndex *string `json:"transactionIndex"`
		L1BlockNumber    *string `json:"l1BlockNumber,omitempty"`
		// 나머지: 단순 string 직접 할당
		BlockHash             string  `json:"blockHash"`
		TransactionHash       string  `json:"transactionHash"`
		From                  string  `json:"from"`
		To                    string  `json:"to"`
		GasUsed               string  `json:"gasUsed"`
		CumulativeGasUsed     string  `json:"cumulativeGasUsed"`
		ContractAddress       *string `json:"contractAddress,omitempty"`
		Logs                  []*Log  `json:"logs"`
		Type                  string  `json:"type"`
		EffectiveGasPrice     string  `json:"effectiveGasPrice"`
		Root                  *string `json:"root,omitempty"`
		Status                *string `json:"status,omitempty"`
		LogsBloom             string  `json:"logsBloom"`
		BlobGasPrice          *string `json:"blobGasPrice,omitempty"`
		BlobGasUsed           *string `json:"blobGasUsed,omitempty"`
		GasUsedForL1          *string `json:"gasUsedForL1,omitempty"`
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
	var dec receiptWire
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// 단순 string 필드 직접 할당
	r.BlockHash = dec.BlockHash
	r.TransactionHash = dec.TransactionHash
	r.From = dec.From
	r.To = dec.To
	r.GasUsed = dec.GasUsed
	r.CumulativeGasUsed = dec.CumulativeGasUsed
	r.ContractAddress = dec.ContractAddress
	r.Logs = dec.Logs
	r.Type = dec.Type
	r.EffectiveGasPrice = dec.EffectiveGasPrice
	r.Root = dec.Root
	r.Status = dec.Status
	r.LogsBloom = dec.LogsBloom
	r.BlobGasPrice = dec.BlobGasPrice
	r.BlobGasUsed = dec.BlobGasUsed
	r.GasUsedForL1 = dec.GasUsedForL1
	r.L1GasPrice = dec.L1GasPrice
	r.L1GasUsed = dec.L1GasUsed
	r.L1FeeScalar = dec.L1FeeScalar
	r.L1Fee = dec.L1Fee
	r.DepositNonce = dec.DepositNonce
	r.DepositReceiptVersion = dec.DepositReceiptVersion
	r.L1BlobBaseFee = dec.L1BlobBaseFee
	r.L1BaseFeeScalar = dec.L1BaseFeeScalar
	r.L1BlobBaseFeeScalar = dec.L1BlobBaseFeeScalar

	// hex→uint64 특수 처리
	if dec.BlockNumber != nil {
		blockNumber, err := hexutil.DecodeUint64(*dec.BlockNumber)
		if err != nil {
			return err
		}
		r.BlockNumber = blockNumber
	}
	if dec.TransactionIndex != nil {
		transactionIndex, err := hexutil.DecodeUint64(*dec.TransactionIndex)
		if err != nil {
			return err
		}
		r.TransactionIndex = transactionIndex
	}
	if dec.L1BlockNumber != nil {
		l1BlockNumber, err := hexutil.DecodeUint64(*dec.L1BlockNumber)
		if err != nil {
			return err
		}
		r.L1BlockNumber = &l1BlockNumber
	}

	return nil
}
