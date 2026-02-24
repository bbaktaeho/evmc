package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (s *SimulateBlockResult) UnmarshalJSON(input []byte) error {
	type wire struct {
		_block
		Number        *hexutil.Big         `json:"number"`
		Hash          *common.Hash         `json:"hash"`
		Timestamp     *hexutil.Uint64      `json:"timestamp"`
		GasLimit      *hexutil.Uint64      `json:"gasLimit"`
		GasUsed       *hexutil.Uint64      `json:"gasUsed"`
		FeeRecipient  *common.Address      `json:"miner"`
		BaseFeePerGas *hexutil.Big         `json:"baseFeePerGas,omitempty"`
		Calls         []*SimulateCallResult `json:"calls"`
		Transactions  interface{}          `json:"transactions"`
		Withdrawals   []*Withdrawal        `json:"withdrawals,omitempty"`
	}
	var dec wire
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Number != nil {
		s.Number = dec.Number.ToInt().Uint64()
	}
	if dec.Hash != nil {
		s.Hash = dec.Hash.Hex()
	}
	if dec.Timestamp != nil {
		s.Timestamp = uint64(*dec.Timestamp)
	}
	if dec.GasLimit != nil {
		s.GasLimit = hexutil.EncodeUint64(uint64(*dec.GasLimit))
	}
	if dec.GasUsed != nil {
		s.GasUsed = hexutil.EncodeUint64(uint64(*dec.GasUsed))
	}
	if dec.FeeRecipient != nil {
		s.FeeRecipient = dec.FeeRecipient.Hex()
	}
	if dec.BaseFeePerGas != nil {
		s.BaseFeePerGas = hexutil.EncodeBig(dec.BaseFeePerGas.ToInt())
	}
	s.Calls = dec.Calls
	s.Transactions = dec.Transactions
	s.Withdrawals = dec.Withdrawals
	return nil
}

func (s *SimulateCallResult) UnmarshalJSON(input []byte) error {
	type wire struct {
		ReturnValue *hexutil.Bytes  `json:"returnData"`
		Logs        []*Log          `json:"logs"`
		GasUsed     *hexutil.Uint64 `json:"gasUsed"`
		Status      *hexutil.Uint64 `json:"status"`
		Error       *CallError      `json:"error,omitempty"`
	}
	var dec wire
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ReturnValue != nil {
		s.ReturnValue = dec.ReturnValue.String()
	}
	if dec.Logs != nil {
		s.Logs = dec.Logs
	}
	if dec.GasUsed != nil {
		s.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.Status != nil {
		s.Status = uint64(*dec.Status)
	}
	if dec.Error != nil {
		s.Error = dec.Error
	}
	return nil
}

func (c *CallError) UnmarshalJSON(input []byte) error {
	type wire struct {
		Code    *int    `json:"code"`
		Message *string `json:"message"`
		Data    *string `json:"data,omitempty"`
	}
	var dec wire
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Code != nil {
		c.Code = *dec.Code
	}
	if dec.Message != nil {
		c.Message = *dec.Message
	}
	if dec.Data != nil {
		c.Data = *dec.Data
	}
	return nil
}
