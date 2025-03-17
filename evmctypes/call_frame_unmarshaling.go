package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

func (c *CallFrame) UnmarshalJSON(input []byte) error {
	type CallFrame0 struct {
		From         *string      `json:"from" validate:"required"`
		Gas          *string      `json:"gas" validate:"required"`
		GasUsed      *string      `json:"gasUsed" validate:"required"`
		To           *string      `json:"to,omitempty"`
		Input        *string      `json:"input"`
		Output       *string      `json:"output,omitempty"`
		Error        *string      `json:"error,omitempty"`
		RevertReason *string      `json:"revertReason,omitempty"`
		Calls        []*CallFrame `json:"calls,omitempty"`
		Logs         []*CallLog   `json:"logs,omitempty"`
		Value        *string      `json:"value,omitempty"`
		Type         *string      `json:"type"`

		// Arbitrum
		AfterEVMTransfers  []*EVMTransfer `json:"afterEVMTransfers,omitempty"`
		BeforeEVMTransfers []*EVMTransfer `json:"beforeEVMTransfers,omitempty"`
	}
	var dec CallFrame0
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.From != nil {
		c.From = *dec.From
	}
	if dec.Gas != nil {
		c.Gas = *dec.Gas
	}
	if dec.GasUsed != nil {
		c.GasUsed = *dec.GasUsed
	}
	if dec.To != nil {
		c.To = dec.To
	}
	if dec.Input != nil {
		c.Input = dec.Input
	}
	if dec.Output != nil {
		c.Output = dec.Output
	}
	if dec.Error != nil {
		c.Error = dec.Error
	}
	if dec.RevertReason != nil {
		c.RevertReason = dec.RevertReason
	}
	if dec.Calls != nil {
		c.Calls = dec.Calls
	}
	if dec.Logs != nil {
		c.Logs = dec.Logs
	}
	if dec.Value != nil {
		valueBig, err := hexutil.DecodeBig(*dec.Value)
		if err != nil {
			return err
		}
		value := decimal.NewFromBigInt(valueBig, 0)
		c.Value = &value
	}
	if dec.Type != nil {
		c.Type = *dec.Type
	}
	if dec.AfterEVMTransfers != nil {
		c.AfterEVMTransfers = dec.AfterEVMTransfers
	}
	if dec.BeforeEVMTransfers != nil {
		c.BeforeEVMTransfers = dec.BeforeEVMTransfers
	}
	return nil
}
