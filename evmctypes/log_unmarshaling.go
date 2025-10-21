package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// UnmarshalJSON unmarshals from JSON.
func (l *Log) UnmarshalJSON(input []byte) error {
	type Log struct {
		Address          *string  `json:"address" validate:"required"`
		Topics           []string `json:"topics" validate:"required"`
		Data             *string  `json:"data" validate:"required"`
		BlockNumber      *string  `json:"blockNumber" validate:"-"`
		TransactionHash  *string  `json:"transactionHash" validate:"required"`
		TransactionIndex *string  `json:"transactionIndex" validate:"-"`
		BlockHash        *string  `json:"blockHash" validate:"required"`
		LogIndex         *string  `json:"logIndex" validate:"-"`
		Removed          *bool    `json:"removed" validate:"-"`
		BlockTimestamp   *string  `json:"blockTimestamp" validate:"-"`
	}
	var dec Log
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Address != nil {
		l.Address = *dec.Address
	}
	if dec.Topics != nil {
		l.Topics = dec.Topics
	}
	if dec.Data != nil {
		l.Data = *dec.Data
	}
	if dec.BlockNumber != nil {
		blockNumber, err := hexutil.DecodeUint64(*dec.BlockNumber)
		if err != nil {
			return err
		}
		l.BlockNumber = blockNumber
	}
	if dec.TransactionHash != nil {
		l.TransactionHash = *dec.TransactionHash
	}
	if dec.TransactionIndex != nil {
		transactionIndex, err := hexutil.DecodeUint64(*dec.TransactionIndex)
		if err != nil {
			return err
		}
		l.TransactionIndex = transactionIndex
	}
	if dec.BlockHash != nil {
		l.BlockHash = *dec.BlockHash
	}
	if dec.LogIndex != nil {
		logIndex, err := hexutil.DecodeUint64(*dec.LogIndex)
		if err != nil {
			return err
		}
		l.LogIndex = logIndex
	}
	if dec.Removed != nil {
		l.Removed = *dec.Removed
	}
	if dec.BlockTimestamp != nil {
		blockTimestamp, err := hexutil.DecodeUint64(*dec.BlockTimestamp)
		if err != nil {
			return err
		}
		l.BlockTimestamp = blockTimestamp
	}
	return nil
}
