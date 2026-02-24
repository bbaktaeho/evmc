package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (l *Log) UnmarshalJSON(input []byte) error {
	type logWire struct {
		// hex→uint64 후처리 대상
		BlockNumber      *string `json:"blockNumber"`
		TransactionIndex *string `json:"transactionIndex"`
		LogIndex         *string `json:"logIndex"`
		BlockTimestamp   *string `json:"blockTimestamp"`
		// 나머지: 단순 string 직접 할당
		Address         string   `json:"address"`
		Topics          []string `json:"topics"`
		Data            string   `json:"data"`
		TransactionHash string   `json:"transactionHash"`
		BlockHash       string   `json:"blockHash"`
		Removed         *bool    `json:"removed"`
	}
	var dec logWire
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// 단순 string 필드 직접 할당
	l.Address = dec.Address
	l.Topics = dec.Topics
	l.Data = dec.Data
	l.TransactionHash = dec.TransactionHash
	l.BlockHash = dec.BlockHash
	if dec.Removed != nil {
		l.Removed = *dec.Removed
	}

	// hex→uint64 특수 처리
	if dec.BlockNumber != nil {
		blockNumber, err := hexutil.DecodeUint64(*dec.BlockNumber)
		if err != nil {
			return err
		}
		l.BlockNumber = blockNumber
	}
	if dec.TransactionIndex != nil {
		transactionIndex, err := hexutil.DecodeUint64(*dec.TransactionIndex)
		if err != nil {
			return err
		}
		l.TransactionIndex = transactionIndex
	}
	if dec.LogIndex != nil {
		logIndex, err := hexutil.DecodeUint64(*dec.LogIndex)
		if err != nil {
			return err
		}
		l.LogIndex = logIndex
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
