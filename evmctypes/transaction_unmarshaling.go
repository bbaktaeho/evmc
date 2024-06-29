package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

func (t *Transaction) UnmarshalJSON(input []byte) error {
	type Transaction struct {
		BlockHash             *string   `json:"blockHash" validate:"required"`
		BlockNumber           *string   `json:"blockNumber" validate:"-"`
		From                  *string   `json:"from" validate:"required"`
		Gas                   *string   `json:"gas" validate:"required"`
		GasPrice              *string   `json:"gasPrice" validate:"required"`
		Hash                  *string   `json:"hash" validate:"required"`
		Input                 *string   `json:"input" validate:"required"`
		Nonce                 *string   `json:"nonce" validate:"required"`
		To                    *string   `json:"to" validate:"required"`
		TransactionIndex      *string   `json:"transactionIndex" validate:"-"`
		Value                 *string   `json:"value" validate:"required"`
		Type                  *string   `json:"type" validate:"required"`
		V                     *string   `json:"v" validate:"required"`
		R                     *string   `json:"r" validate:"required"`
		S                     *string   `json:"s" validate:"required"`
		YParity               *string   `json:"yParity,omitempty"`
		ChainID               *string   `json:"chainId,omitempty"`
		AccessList            []*Access `json:"accessList,omitempty"`
		MaxFeePerGas          *string   `json:"maxFeePerGas,omitempty"`
		MaxPriorityFeePerGas  *string   `json:"maxPriorityFeePerGas,omitempty"`
		MaxFeePerBlobGas      *string   `json:"maxFeePerBlobGas,omitempty"`
		BlobVersionedHashes   []string  `json:"blobVersionedHashes,omitempty"`
		L1BlockNumber         *string   `json:"l1BlockNumber,omitempty"`
		RequestID             *string   `json:"requestId,omitempty"`
		TicketID              *string   `json:"ticketId,omitempty"`
		MaxRefund             *string   `json:"maxRefund,omitempty"`
		SubmissionFeeRefund   *string   `json:"submissionFeeRefund,omitempty"`
		RefundTo              *string   `json:"refundTo,omitempty"`
		L1BaseFee             *string   `json:"l1BaseFee,omitempty"`
		DepositValue          *string   `json:"depositValue,omitempty"`
		RetryTo               *string   `json:"retryTo,omitempty"`
		RetryValue            *string   `json:"retryValue,omitempty"`
		RetryData             *string   `json:"retryData,omitempty"`
		Beneficiary           *string   `json:"beneficiary,omitempty"`
		MaxSubmissionFee      *string   `json:"maxSubmissionFee,omitempty"`
		EffectiveGasPrice     *string   `json:"effectiveGasPrice,omitempty"`
		QueueOrigin           *string   `json:"queueOrigin,omitempty"`
		L1TxOrigin            *string   `json:"l1TxOrigin,omitempty"`
		L1BlockTimestamp      *string   `json:"l1Timestamp,omitempty"`
		Index                 *string   `json:"index,omitempty"`
		QueueIndex            *string   `json:"queueIndex,omitempty"`
		SourceHash            *string   `json:"sourceHash,omitempty"`
		Mint                  *string   `json:"mint,omitempty"`
		IsSystemTx            *bool     `json:"isSystemTx,omitempty"`
		DepositReceiptVersion *string   `json:"depositReceiptVersion,omitempty"`
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
	if dec.Input != nil {
		t.Input = *dec.Input
	}
	if dec.Nonce != nil {
		t.Nonce = *dec.Nonce
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
	if dec.Value != nil {
		valueBig, err := hexutil.DecodeBig(*dec.Value)
		if err != nil {
			return err
		}
		t.Value = decimal.NewFromBigInt(valueBig, 0)
	}
	if dec.Type != nil {
		t.Type = *dec.Type
	}
	if dec.V != nil {
		t.V = *dec.V
	}
	if dec.R != nil {
		t.R = *dec.R
	}
	if dec.S != nil {
		t.S = *dec.S
	}
	if dec.YParity != nil {
		t.YParity = dec.YParity
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
	if dec.MaxFeePerBlobGas != nil {
		t.MaxFeePerBlobGas = dec.MaxFeePerBlobGas
	}
	if dec.BlobVersionedHashes != nil {
		t.BlobVersionedHashes = dec.BlobVersionedHashes
	}
	if dec.L1BlockNumber != nil {
		l1BlockNumber, err := hexutil.DecodeUint64(*dec.L1BlockNumber)
		if err != nil {
			return err
		}
		t.L1BlockNumber = &l1BlockNumber
	}
	if dec.RequestID != nil {
		t.RequestID = dec.RequestID
	}
	if dec.TicketID != nil {
		t.TicketID = dec.TicketID
	}
	if dec.MaxRefund != nil {
		t.MaxRefund = dec.MaxRefund
	}
	if dec.SubmissionFeeRefund != nil {
		t.SubmissionFeeRefund = dec.SubmissionFeeRefund
	}
	if dec.RefundTo != nil {
		t.RefundTo = dec.RefundTo
	}
	if dec.L1BaseFee != nil {
		t.L1BaseFee = dec.L1BaseFee
	}
	if dec.DepositValue != nil {
		t.DepositValue = dec.DepositValue
	}
	if dec.RetryTo != nil {
		t.RetryTo = dec.RetryTo
	}
	if dec.RetryValue != nil {
		t.RetryValue = dec.RetryValue
	}
	if dec.RetryData != nil {
		t.RetryData = dec.RetryData
	}
	if dec.Beneficiary != nil {
		t.Beneficiary = dec.Beneficiary
	}
	if dec.MaxSubmissionFee != nil {
		t.MaxSubmissionFee = dec.MaxSubmissionFee
	}
	if dec.EffectiveGasPrice != nil {
		t.EffectiveGasPrice = dec.EffectiveGasPrice
	}
	if dec.QueueOrigin != nil {
		t.QueueOrigin = dec.QueueOrigin
	}
	if dec.L1TxOrigin != nil {
		t.L1TxOrigin = dec.L1TxOrigin
	}
	if dec.L1BlockTimestamp != nil {
		t.L1BlockTimestamp = dec.L1BlockTimestamp
	}
	if dec.Index != nil {
		index, err := hexutil.DecodeUint64(*dec.Index)
		if err != nil {
			return err
		}
		t.Index = &index
	}
	if dec.QueueIndex != nil {
		queueIndex, err := hexutil.DecodeUint64(*dec.QueueIndex)
		if err != nil {
			return err
		}
		t.QueueIndex = &queueIndex
	}
	if dec.SourceHash != nil {
		t.SourceHash = dec.SourceHash
	}
	if dec.Mint != nil {
		t.Mint = dec.Mint
	}
	if dec.IsSystemTx != nil {
		t.IsSystemTx = dec.IsSystemTx
	}
	if dec.DepositReceiptVersion != nil {
		t.DepositReceiptVersion = dec.DepositReceiptVersion
	}
	return nil
}
