package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

func (f *FlatCallFrame) UnmarshalJSON(input []byte) error {
	type FlatCallFrame struct {
		Action *struct {
			Author         *string `json:"author,omitempty"`
			RewardType     *string `json:"rewardType,omitempty"`
			Address        *string `json:"address,omitempty"`
			Balance        *string `json:"balance,omitempty"`
			CreationMethod *string `json:"creationMethod,omitempty"`
			RefundAddress  *string `json:"refundAddress,omitempty"`
			CallType       *string `json:"callType,omitempty"`
			From           *string `json:"from,omitempty"`
			Gas            *string `json:"gas,omitempty"`
			Input          *string `json:"input,omitempty"`
			To             *string `json:"to,omitempty"`
			Init           *string `json:"init,omitempty"`
			Value          *string `json:"value,omitempty"`
		} `json:"action"`
		BlockHash   *string `json:"blockHash" validate:"required"`
		BlockNumber *uint64 `json:"blockNumber" validate:"-"`
		Error       *string `json:"error,omitempty"`
		Result      *struct {
			Address *string `json:"address,omitempty"`
			Code    *string `json:"code,omitempty"`
			GasUsed *string `json:"gasUsed,omitempty"`
			Output  *string `json:"output,omitempty"`
		} `json:"result,omitempty"`
		Subtraces           *uint64  `json:"subtraces" validate:"required"`
		TraceAddress        []uint64 `json:"traceAddress" validate:"required"`
		TransactionHash     *string  `json:"transactionHash" validate:"required"`
		TransactionPosition *uint64  `json:"transactionPosition" validate:"-"`
		Type                *string  `json:"type" validate:"required"`

		// Arbitrum
		AfterEVMTransfers  []*EVMTransfer `json:"afterEVMTransfers,omitempty"`
		BeforeEVMTransfers []*EVMTransfer `json:"beforeEVMTransfers,omitempty"`
	}
	var dec FlatCallFrame
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Action != nil {
		f.Action.Author = dec.Action.Author
		f.Action.RewardType = dec.Action.RewardType
		f.Action.Address = dec.Action.Address
		f.Action.Balance = dec.Action.Balance
		f.Action.CreationMethod = dec.Action.CreationMethod
		f.Action.RefundAddress = dec.Action.RefundAddress
		f.Action.CallType = dec.Action.CallType
		f.Action.From = dec.Action.From
		f.Action.Gas = dec.Action.Gas
		f.Action.Input = dec.Action.Input
		f.Action.To = dec.Action.To
		f.Action.Init = dec.Action.Init
		if *dec.Action.Value == "" {
			zero := decimal.Zero
			f.Action.Value = &zero
		} else {
			valueBig, err := hexutil.DecodeBig(*dec.Action.Value)
			if err != nil {
				return err
			}
			value := decimal.NewFromBigInt(valueBig, 0)
			f.Action.Value = &value
		}
	}
	if dec.BlockHash != nil {
		f.BlockHash = *dec.BlockHash
	}
	if dec.BlockNumber != nil {
		f.BlockNumber = *dec.BlockNumber
	}
	if dec.Error != nil {
		f.Error = dec.Error
	}
	if dec.Result != nil {
		f.Result = dec.Result
	}
	if dec.Subtraces != nil {
		f.Subtraces = *dec.Subtraces
	}
	if dec.TraceAddress != nil {
		f.TraceAddress = dec.TraceAddress
	}
	if dec.TransactionHash != nil {
		f.TransactionHash = *dec.TransactionHash
	}
	if dec.TransactionPosition != nil {
		f.TransactionPosition = *dec.TransactionPosition
	}
	if dec.Type != nil {
		f.Type = *dec.Type
	}
	if dec.AfterEVMTransfers != nil {
		f.AfterEVMTransfers = dec.AfterEVMTransfers
	}
	if dec.BeforeEVMTransfers != nil {
		f.BeforeEVMTransfers = dec.BeforeEVMTransfers
	}
	return nil
}
