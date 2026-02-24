package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

func (t *Transaction) UnmarshalJSON(input []byte) error {
	type txWire struct {
		// hex→uint64 후처리 대상
		BlockNumber      *string `json:"blockNumber"`
		TransactionIndex *string `json:"transactionIndex"`
		L1BlockNumber    *string `json:"l1BlockNumber,omitempty"`
		Index            *string `json:"index,omitempty"`
		QueueIndex       *string `json:"queueIndex,omitempty"`
		// hex→decimal 후처리 대상
		Value *string `json:"value"`
		// 나머지: 단순 string 직접 할당
		BlockHash            string           `json:"blockHash"`
		From                 string           `json:"from"`
		Gas                  string           `json:"gas"`
		GasPrice             string           `json:"gasPrice"`
		Hash                 string           `json:"hash"`
		Input                string           `json:"input"`
		Nonce                string           `json:"nonce"`
		To                   string           `json:"to"`
		Type                 string           `json:"type"`
		V                    string           `json:"v"`
		R                    string           `json:"r"`
		S                    string           `json:"s"`
		YParity              *string          `json:"yParity,omitempty"`
		ChainID              *string          `json:"chainId,omitempty"`
		AccessList           []*Access        `json:"accessList,omitempty"`
		MaxFeePerGas         *string          `json:"maxFeePerGas,omitempty"`
		MaxPriorityFeePerGas *string          `json:"maxPriorityFeePerGas,omitempty"`
		MaxFeePerBlobGas     *string          `json:"maxFeePerBlobGas,omitempty"`
		BlobVersionedHashes  []string         `json:"blobVersionedHashes,omitempty"`
		AuthorizationList    []*Authorization `json:"authorizationList,omitempty"`
		RequestID            *string          `json:"requestId,omitempty"`
		TicketID             *string          `json:"ticketId,omitempty"`
		MaxRefund            *string          `json:"maxRefund,omitempty"`
		SubmissionFeeRefund  *string          `json:"submissionFeeRefund,omitempty"`
		RefundTo             *string          `json:"refundTo,omitempty"`
		L1BaseFee            *string          `json:"l1BaseFee,omitempty"`
		DepositValue         *string          `json:"depositValue,omitempty"`
		RetryTo              *string          `json:"retryTo,omitempty"`
		RetryValue           *string          `json:"retryValue,omitempty"`
		RetryData            *string          `json:"retryData,omitempty"`
		Beneficiary          *string          `json:"beneficiary,omitempty"`
		MaxSubmissionFee     *string          `json:"maxSubmissionFee,omitempty"`
		EffectiveGasPrice    *string          `json:"effectiveGasPrice,omitempty"`
		QueueOrigin          *string          `json:"queueOrigin,omitempty"`
		L1TxOrigin           *string          `json:"l1TxOrigin,omitempty"`
		L1BlockTimestamp     *string          `json:"l1Timestamp,omitempty"`
		SourceHash           *string          `json:"sourceHash,omitempty"`
		Mint                 *string          `json:"mint,omitempty"`
		IsSystemTx           *bool            `json:"isSystemTx,omitempty"`
		DepositReceiptVersion *string         `json:"depositReceiptVersion,omitempty"`
	}
	var dec txWire
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// 단순 string 필드 직접 할당
	t.BlockHash = dec.BlockHash
	t.From = dec.From
	t.Gas = dec.Gas
	t.GasPrice = dec.GasPrice
	t.Hash = dec.Hash
	t.Input = dec.Input
	t.Nonce = dec.Nonce
	t.To = dec.To
	t.Type = dec.Type
	t.V = dec.V
	t.R = dec.R
	t.S = dec.S
	t.YParity = dec.YParity
	t.ChainID = dec.ChainID
	t.AccessList = dec.AccessList
	t.MaxFeePerGas = dec.MaxFeePerGas
	t.MaxPriorityFeePerGas = dec.MaxPriorityFeePerGas
	t.MaxFeePerBlobGas = dec.MaxFeePerBlobGas
	t.BlobVersionedHashes = dec.BlobVersionedHashes
	t.AuthorizationList = dec.AuthorizationList
	t.RequestID = dec.RequestID
	t.TicketID = dec.TicketID
	t.MaxRefund = dec.MaxRefund
	t.SubmissionFeeRefund = dec.SubmissionFeeRefund
	t.RefundTo = dec.RefundTo
	t.L1BaseFee = dec.L1BaseFee
	t.DepositValue = dec.DepositValue
	t.RetryTo = dec.RetryTo
	t.RetryValue = dec.RetryValue
	t.RetryData = dec.RetryData
	t.Beneficiary = dec.Beneficiary
	t.MaxSubmissionFee = dec.MaxSubmissionFee
	t.EffectiveGasPrice = dec.EffectiveGasPrice
	t.QueueOrigin = dec.QueueOrigin
	t.L1TxOrigin = dec.L1TxOrigin
	t.L1BlockTimestamp = dec.L1BlockTimestamp
	t.SourceHash = dec.SourceHash
	t.Mint = dec.Mint
	t.IsSystemTx = dec.IsSystemTx
	t.DepositReceiptVersion = dec.DepositReceiptVersion

	// hex→uint64 특수 처리
	if dec.BlockNumber != nil {
		blockNumber, err := hexutil.DecodeUint64(*dec.BlockNumber)
		if err != nil {
			return err
		}
		t.BlockNumber = blockNumber
	}
	if dec.TransactionIndex != nil {
		transactionIndex, err := hexutil.DecodeUint64(*dec.TransactionIndex)
		if err != nil {
			return err
		}
		t.TransactionIndex = transactionIndex
	}
	if dec.L1BlockNumber != nil {
		l1BlockNumber, err := hexutil.DecodeUint64(*dec.L1BlockNumber)
		if err != nil {
			return err
		}
		t.L1BlockNumber = &l1BlockNumber
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

	// hex→decimal 특수 처리
	if dec.Value != nil {
		valueBig, err := hexutil.DecodeBig(*dec.Value)
		if err != nil {
			return err
		}
		t.Value = decimal.NewFromBigInt(valueBig, 0)
	}

	return nil
}
