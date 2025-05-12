package kaiatypes

import (
	"math/big"

	"github.com/shopspring/decimal"
)

type Header struct {
	block
}

type Block struct {
	block
	Transactions []string `json:"transactions"`
}

type BlockIncTx struct {
	block
	Transactions []*Transaction `json:"transactions"`
}

type Transaction struct {
	BlockHash          string          `json:"blockHash" validate:"required"`
	BlockNumber        uint64          `json:"blockNumber" validate:"-"`
	CodeFormat         string          `json:"codeFormat,omitempty"`
	FeePayer           string          `json:"feePayer,omitempty"`
	FeePayerSignatures []*Signature    `json:"feePayerSignatures,omitempty"`
	FeeRatio           string          `json:"feeRatio,omitempty"`
	From               string          `json:"from" validate:"required"`
	Gas                string          `json:"gas" validate:"required"`
	GasPrice           string          `json:"gasPrice" validate:"required"`
	Hash               string          `json:"hash" validate:"required"`
	HumanReadable      bool            `json:"humanReadable,omitempty"`
	Key                string          `json:"key,omitempty"`
	Input              string          `json:"input"`
	Nonce              string          `json:"nonce" validate:"required"`
	SenderTxHash       string          `json:"senderTxHash" validate:"required"`
	Signatures         []*Signature    `json:"signatures" validate:"required"`
	To                 string          `json:"to" validate:"required"`
	TransactionIndex   uint64          `json:"transactionIndex" validate:"-"`
	Type               string          `json:"type" validate:"required"`
	TypeInt            uint64          `json:"typeInt" validate:"required"`
	Value              decimal.Decimal `json:"value" validate:"required"`

	ChainID              *string   `json:"chainId,omitempty"`              // EIP-155
	AccessList           []*Access `json:"accessList,omitempty"`           // EIP-2930
	MaxFeePerGas         *string   `json:"maxFeePerGas,omitempty"`         // EIP-1559
	MaxPriorityFeePerGas *string   `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559
}

type Signature struct {
	V string `json:"v" validate:"required"`
	R string `json:"r" validate:"required"`
	S string `json:"s" validate:"required"`
}

type Access struct {
	Address     string   `json:"address" validate:"required"`
	StorageKeys []string `json:"storageKeys" validate:"required"`
}

type Receipt struct {
	BlockHash          string       `json:"blockHash" validate:"required"`
	BlockNumber        uint64       `json:"blockNumber" validate:"-"`
	CodeFormat         string       `json:"codeFormat,omitempty"`
	ContractAddress    string       `json:"contractAddress,omitempty"`
	FeePayer           string       `json:"feePayer,omitempty"`
	FeePayerSignatures []*Signature `json:"feePayerSignatures,omitempty"`
	FeeRatio           string       `json:"feeRatio,omitempty"`
	From               string       `json:"from" validate:"required"`
	Gas                string       `json:"gas" validate:"required"`
	EffectiveGasPrice  string       `json:"effectiveGasPrice" validate:"required"`
	GasPrice           string       `json:"gasPrice" validate:"required"`
	GasUsed            string       `json:"gasUsed" validate:"required"`
	HumanReadable      bool         `json:"humanReadable,omitempty"`
	Key                string       `json:"key,omitempty"`
	Logs               []*Log       `json:"logs" validate:"required"`
	LogsBloom          string       `json:"logsBloom"`
	Nonce              string       `json:"nonce" validate:"required"`
	SenderTxHash       string       `json:"senderTxHash" validate:"required"`
	Signatures         []*Signature `json:"signatures" validate:"required"`
	Status             string       `json:"status" validate:"required"`
	TxError            string       `json:"txError,omitempty"`
	To                 string       `json:"to" validate:"required"`
	TransactionHash    string       `json:"transactionHash" validate:"required"`
	TransactionIndex   uint64       `json:"transactionIndex" validate:"-"`
	Type               string       `json:"type" validate:"required"`
	TypeInt            uint64       `json:"typeInt" validate:"required"`
	Value              string       `json:"value" validate:"required"`

	Input                string    `json:"input"`
	ChainID              *string   `json:"chainId,omitempty"`              // EIP-155
	AccessList           []*Access `json:"accessList,omitempty"`           // EIP-2930
	MaxFeePerGas         *string   `json:"maxFeePerGas,omitempty"`         // EIP-1559
	MaxPriorityFeePerGas *string   `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559
}

type Log struct {
	Address          string   `json:"address" validate:"required"`
	Topics           []string `json:"topics" validate:"required"`
	Data             string   `json:"data" validate:"required"`
	BlockNumber      uint64   `json:"blockNumber" validate:"-"`
	TransactionHash  string   `json:"transactionHash" validate:"required"`
	TransactionIndex uint64   `json:"transactionIndex" validate:"-"`
	BlockHash        string   `json:"blockHash" validate:"required"`
	LogIndex         uint64   `json:"logIndex" validate:"-"`
	Removed          bool     `json:"removed" validate:"-"`
}

type Rewards struct {
	BurntFee *big.Int            `json:"burntFee" validate:"required"`
	Kgf      *big.Int            `json:"kgf" validate:"required"`
	Kir      *big.Int            `json:"kir" validate:"required"`
	Minted   *big.Int            `json:"minted" validate:"required"`
	Proposer *big.Int            `json:"proposer" validate:"required"`
	Rewards  map[string]*big.Int `json:"rewards" validate:"required"`
	Stakers  *big.Int            `json:"stakers" validate:"required"`
	TotalFee *big.Int            `json:"totalFee" validate:"required"`
}
