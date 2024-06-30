package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

type block struct {
	Number           uint64 `json:"number" validate:"-"`
	Hash             string `json:"hash" validate:"required"`
	ParentHash       string `json:"parentHash" validate:"required"`
	Nonce            string `json:"nonce" validate:"required"`
	MixHash          string `json:"mixHash" validate:"required"`
	Sha3Uncles       string `json:"sha3Uncles" validate:"required"`
	LogsBloom        string `json:"logsBloom" validate:"required"`
	StateRoot        string `json:"stateRoot" validate:"required"`
	Miner            string `json:"miner" validate:"required"`
	Difficulty       string `json:"difficulty" validate:"required"`
	ExtraData        string `json:"extraData" validate:"required"`
	GasLimit         string `json:"gasLimit" validate:"required"`
	GasUsed          string `json:"gasUsed" validate:"required"`
	Timestamp        uint64 `json:"timestamp" validate:"required"`
	TransactionsRoot string `json:"transactionsRoot" validate:"required"`
	ReceiptsRoot     string `json:"receiptsRoot" validate:"required"`

	TotalDifficulty string   `json:"totalDifficulty,omitempty"` // Unused fields in newHeads
	Size            string   `json:"size,omitempty"`            // Unused fields in newHeads
	Uncles          []string `json:"uncles,omitempty"`          // Unused fields in newHeads

	BaseFeePerGas         *string `json:"baseFeePerGas,omitempty"`         // EIP-1559
	WithdrawalsRoot       *string `json:"withdrawalsRoot,omitempty"`       // EIP-4895
	BlobGasUsed           *string `json:"blobGasUsed,omitempty"`           // EIP-4844
	ExcessBlobGas         *string `json:"excessBlobGas,omitempty"`         // EIP-4844
	ParentBeaconBlockRoot *string `json:"parentBeaconBlockRoot,omitempty"` // EIP-4788

	L1BlockNumber *uint64 `json:"l1BlockNumber,omitempty"` // Arbitrum
	SendCount     *string `json:"sendCount,omitempty"`     // Arbitrum
	SendRoot      *string `json:"sendRoot,omitempty"`      // Arbitrum
}

func (b *block) NextBaseFee() decimal.Decimal {
	if b.BaseFeePerGas == nil {
		return decimal.Zero
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*b.BaseFeePerGas), 0).Mul(decimal.NewFromInt(2))
}

// _block is raw data from the blockchain RPC calls.
type _block struct {
	Number           *string  `json:"number"`
	Hash             *string  `json:"hash"`
	ParentHash       *string  `json:"parentHash"`
	Nonce            *string  `json:"nonce"`
	MixHash          *string  `json:"mixHash"`
	Sha3Uncles       *string  `json:"sha3Uncles"`
	LogsBloom        *string  `json:"logsBloom"`
	StateRoot        *string  `json:"stateRoot"`
	Miner            *string  `json:"miner"`
	Difficulty       *string  `json:"difficulty"`
	ExtraData        *string  `json:"extraData"`
	GasLimit         *string  `json:"gasLimit"`
	GasUsed          *string  `json:"gasUsed"`
	Timestamp        *string  `json:"timestamp"`
	TransactionsRoot *string  `json:"transactionsRoot"`
	ReceiptsRoot     *string  `json:"receiptsRoot"`
	TotalDifficulty  *string  `json:"totalDifficulty"`
	Size             *string  `json:"size"`
	Uncles           []string `json:"uncles"`

	BaseFeePerGas         *string `json:"baseFeePerGas,omitempty"`
	WithdrawalsRoot       *string `json:"withdrawalsRoot,omitempty"`
	BlobGasUsed           *string `json:"blobGasUsed,omitempty"`
	ExcessBlobGas         *string `json:"excessBlobGas,omitempty"`
	ParentBeaconBlockRoot *string `json:"parentBeaconBlockRoot,omitempty"`

	L1BlockNumber *string `json:"l1BlockNumber,omitempty"`
	SendCount     *string `json:"sendCount,omitempty"`
	SendRoot      *string `json:"sendRoot,omitempty"`
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
	if _b.Nonce != nil {
		b.Nonce = *_b.Nonce
	}
	if _b.MixHash != nil {
		b.MixHash = *_b.MixHash
	}
	if _b.Sha3Uncles != nil {
		b.Sha3Uncles = *_b.Sha3Uncles
	}
	if _b.LogsBloom != nil {
		b.LogsBloom = *_b.LogsBloom
	}
	if _b.StateRoot != nil {
		b.StateRoot = *_b.StateRoot
	}
	if _b.Miner != nil {
		b.Miner = *_b.Miner
	}
	if _b.Difficulty != nil {
		b.Difficulty = *_b.Difficulty
	}
	if _b.ExtraData != nil {
		b.ExtraData = *_b.ExtraData
	}
	if _b.GasLimit != nil {
		b.GasLimit = *_b.GasLimit
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
	if _b.TransactionsRoot != nil {
		b.TransactionsRoot = *_b.TransactionsRoot
	}
	if _b.ReceiptsRoot != nil {
		b.ReceiptsRoot = *_b.ReceiptsRoot
	}
	if _b.TotalDifficulty != nil {
		b.TotalDifficulty = *_b.TotalDifficulty
	}
	if _b.Size != nil {
		b.Size = *_b.Size
	}
	if _b.Uncles != nil {
		b.Uncles = _b.Uncles
	}
	if _b.BaseFeePerGas != nil {
		b.BaseFeePerGas = _b.BaseFeePerGas
	}
	if _b.WithdrawalsRoot != nil {
		b.WithdrawalsRoot = _b.WithdrawalsRoot
	}
	if _b.BlobGasUsed != nil {
		b.BlobGasUsed = _b.BlobGasUsed
	}
	if _b.ExcessBlobGas != nil {
		b.ExcessBlobGas = _b.ExcessBlobGas
	}
	if _b.ParentBeaconBlockRoot != nil {
		b.ParentBeaconBlockRoot = _b.ParentBeaconBlockRoot
	}
	if _b.L1BlockNumber != nil {
		l1BlockNumber, err := hexutil.DecodeUint64(*_b.L1BlockNumber)
		if err != nil {
			return err
		}
		b.L1BlockNumber = &l1BlockNumber
	}
	if _b.SendCount != nil {
		b.SendCount = _b.SendCount
	}
	if _b.SendRoot != nil {
		b.SendRoot = _b.SendRoot
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
		Withdrawals  []*Withdrawal `json:"withdrawals,omitempty"`
		Transactions []string      `json:"transactions"`
		UncleBlocks  []*Block      `json:"uncleBlocks"`
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
	if dec.UncleBlocks != nil {
		b.UncleBlocks = dec.UncleBlocks
	}
	if dec.Withdrawals != nil {
		b.Withdrawals = dec.Withdrawals
	}
	return nil
}

func (b *BlockIncTx) UnmarshalJSON(input []byte) error {
	type blockIncTx struct {
		_block
		Transactions []*Transaction `json:"transactions"`
		Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"`
		UncleBlocks  []*BlockIncTx  `json:"uncleBlocks,omitempty"`
	}
	var dec blockIncTx
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if err := dec.unmarshal(&b.block); err != nil {
		return err
	}
	if dec.Transactions != nil {
		b.Transactions = dec.Transactions
	}
	if dec.UncleBlocks != nil {
		b.UncleBlocks = dec.UncleBlocks
	}
	if dec.Withdrawals != nil {
		b.Withdrawals = dec.Withdrawals
	}
	return nil
}

func (w *Withdrawal) UnmarshalJSON(input []byte) error {
	type withdrawal struct {
		Index          *string `json:"index"`
		ValidatorIndex *string `json:"validatorIndex"`
		Address        *string `json:"address"`
		Amount         *string `json:"amount"`
	}
	var dec withdrawal
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Index != nil {
		index, err := hexutil.DecodeUint64(*dec.Index)
		if err != nil {
			return err
		}
		w.Index = index
	}
	if dec.ValidatorIndex != nil {
		validatorIndex, err := hexutil.DecodeUint64(*dec.ValidatorIndex)
		if err != nil {
			return err
		}
		w.ValidatorIndex = validatorIndex
	}
	if dec.Address != nil {
		w.Address = *dec.Address
	}
	if dec.Amount != nil {
		amount, err := hexutil.DecodeUint64(*dec.Amount)
		if err != nil {
			return err
		}
		w.Amount = amount
	}
	return nil
}
