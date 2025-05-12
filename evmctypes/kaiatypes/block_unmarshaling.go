package kaiatypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type block struct {
	Number           uint64 `json:"number" validate:"-"`
	Hash             string `json:"hash" validate:"required"`
	ParentHash       string `json:"parentHash" validate:"required"`
	MixHash          string `json:"mixHash,omitempty"`
	RandomReveal     string `json:"randomReveal,omitempty"`
	LogsBloom        string `json:"logsBloom,omitempty"`
	TransactionsRoot string `json:"transactionsRoot" validate:"required"`
	StateRoot        string `json:"stateRoot" validate:"required"`
	ReceiptsRoot     string `json:"receiptsRoot" validate:"required"`
	Reward           string `json:"reward" validate:"required"`
	BlockScore       string `json:"blockScore" validate:"required"`
	TotalBlockScore  string `json:"totalBlockScore" validate:"required"`
	ExtraData        string `json:"extraData,omitempty"`
	Size             string `json:"size" validate:"required"`
	GasUsed          string `json:"gasUsed" validate:"required"`
	Timestamp        uint64 `json:"timestamp" validate:"required"`
	TimestampFoS     string `json:"timestampFoS" validate:"required"`
	GovernanceData   string `json:"governanceData" validate:"required"`
	VoteData         string `json:"voteData" validate:"required"`
	BaseFeePerGas    string `json:"baseFeePerGas" validate:"required"` // EIP-1559
}

// _block is raw data from the blockchain RPC calls.
type _block struct {
	Number           *string `json:"number"`
	Hash             *string `json:"hash"`
	ParentHash       *string `json:"parentHash"`
	MixHash          *string `json:"mixHash"`
	RandomReveal     *string `json:"randomReveal"`
	LogsBloom        *string `json:"logsBloom"`
	TransactionsRoot *string `json:"transactionsRoot"`
	StateRoot        *string `json:"stateRoot"`
	ReceiptsRoot     *string `json:"receiptsRoot"`
	Reward           *string `json:"reward"`
	BlockScore       *string `json:"blockScore"`
	TotalBlockScore  *string `json:"totalBlockScore"`
	ExtraData        *string `json:"extraData"`
	Size             *string `json:"size"`
	GasUsed          *string `json:"gasUsed"`
	Timestamp        *string `json:"timestamp"`
	TimestampFoS     *string `json:"timestampFoS"`
	GovernanceData   *string `json:"governanceData"`
	VoteData         *string `json:"voteData"`
	BaseFeePerGas    *string `json:"baseFeePerGas"`
}

func (_b *_block) unmarshal(b *block) error {
	if _b.Number != nil {
		number, err := hexutil.DecodeUint64(*_b.Number)
		if err != nil {
			return err
		}
		b.Number = number
	}
	if _b.Hash != nil {
		b.Hash = *_b.Hash
	}
	if _b.ParentHash != nil {
		b.ParentHash = *_b.ParentHash
	}
	if _b.MixHash != nil {
		b.MixHash = *_b.MixHash
	}
	if _b.RandomReveal != nil {
		b.RandomReveal = *_b.RandomReveal
	}
	if _b.LogsBloom != nil {
		b.LogsBloom = *_b.LogsBloom
	}
	if _b.TransactionsRoot != nil {
		b.TransactionsRoot = *_b.TransactionsRoot
	}
	if _b.StateRoot != nil {
		b.StateRoot = *_b.StateRoot
	}
	if _b.ReceiptsRoot != nil {
		b.ReceiptsRoot = *_b.ReceiptsRoot
	}
	if _b.Reward != nil {
		b.Reward = *_b.Reward
	}
	if _b.BlockScore != nil {
		b.BlockScore = *_b.BlockScore
	}
	if _b.TotalBlockScore != nil {
		b.TotalBlockScore = *_b.TotalBlockScore
	}
	if _b.ExtraData != nil {
		b.ExtraData = *_b.ExtraData
	}
	if _b.Size != nil {
		b.Size = *_b.Size
	}
	if _b.GasUsed != nil {
		b.GasUsed = *_b.GasUsed
	}
	if _b.Timestamp != nil {
		timestamp, err := hexutil.DecodeUint64(*_b.Timestamp)
		if err != nil {
			return err
		}
		b.Timestamp = timestamp
	}
	if _b.TimestampFoS != nil {
		b.TimestampFoS = *_b.TimestampFoS
	}
	if _b.GovernanceData != nil {
		b.GovernanceData = *_b.GovernanceData
	}
	if _b.VoteData != nil {
		b.VoteData = *_b.VoteData
	}
	if _b.BaseFeePerGas != nil {
		b.BaseFeePerGas = *_b.BaseFeePerGas
	}
	return nil
}

func (h *Header) UnmarshalJSON(input []byte) error {
	type header struct {
		_block
	}
	var dec header
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	return dec.unmarshal(&h.block)
}

func (b *Block) UnmarshalJSON(input []byte) error {
	type block struct {
		_block
		Transactions []string `json:"transactions"`
	}
	var dec block
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if err := dec.unmarshal(&b.block); err != nil {
		return err
	}
	if dec.Transactions != nil {
		b.Transactions = dec.Transactions
	}
	return nil
}

func (b *BlockIncTxs) UnmarshalJSON(input []byte) error {
	type block struct {
		_block
		Transactions []*Transaction `json:"transactions"`
	}
	var dec block
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if err := dec.unmarshal(&b.block); err != nil {
		return err
	}
	if dec.Transactions != nil {
		b.Transactions = dec.Transactions
	}
	return nil
}
