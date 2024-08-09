package evmctypes

import (
	"github.com/shopspring/decimal"
)

type Subscription interface {
	Unsubscribe()
	Err() <-chan error
}

type Header struct {
	block
}

type Block struct {
	block
	Transactions []string      `json:"transactions" validate:""`
	UncleBlocks  []*Block      `json:"uncleBlocks,omitempty"`
	Withdrawals  []*Withdrawal `json:"withdrawals,omitempty"` // EIP-4895
}

type BlockIncTx struct {
	block
	Transactions []*Transaction `json:"transactions" validate:""`
	UncleBlocks  []*BlockIncTx  `json:"uncleBlocks,omitempty"`
	Withdrawals  []*Withdrawal  `json:"withdrawals,omitempty"` // EIP-4895
}

type Withdrawal struct {
	Index          uint64 `json:"index" validate:"-"`
	ValidatorIndex uint64 `json:"validatorIndex" validate:"-"`
	Address        string `json:"address" validate:"required"`
	Amount         uint64 `json:"amount" validate:"required"`
}

type Transaction struct {
	BlockHash        string          `json:"blockHash" validate:"required"`
	BlockNumber      uint64          `json:"blockNumber" validate:"-"`
	From             string          `json:"from" validate:"required"`
	Gas              string          `json:"gas" validate:"required"`
	GasPrice         string          `json:"gasPrice" validate:"required"`
	Hash             string          `json:"hash" validate:"required"`
	Input            string          `json:"input" validate:"required"`
	Nonce            string          `json:"nonce" validate:"required"`
	To               string          `json:"to" validate:"required"`
	TransactionIndex uint64          `json:"transactionIndex" validate:"-"`
	Value            decimal.Decimal `json:"value" validate:"required"`
	Type             string          `json:"type" validate:"required"`
	V                string          `json:"v" validate:"required"`
	R                string          `json:"r" validate:"required"`
	S                string          `json:"s" validate:"required"`
	YParity          *string         `json:"yParity,omitempty"`

	ChainID              *string   `json:"chainId,omitempty"`              // EIP-155
	AccessList           []*Access `json:"accessList,omitempty"`           // EIP-2930
	MaxFeePerGas         *string   `json:"maxFeePerGas,omitempty"`         // EIP-1559
	MaxPriorityFeePerGas *string   `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559
	MaxFeePerBlobGas     *string   `json:"maxFeePerBlobGas,omitempty"`     // EIP-4844
	BlobVersionedHashes  []string  `json:"blobVersionedHashes,omitempty"`  // EIP-4844

	L1BlockNumber *uint64 `json:"l1BlockNumber,omitempty"` // Arbitrum, Optimism

	RequestID           *string `json:"requestId,omitempty"`           // Arbitrum
	TicketID            *string `json:"ticketId,omitempty"`            // Arbitrum
	MaxRefund           *string `json:"maxRefund,omitempty"`           // Arbitrum
	SubmissionFeeRefund *string `json:"submissionFeeRefund,omitempty"` // Arbitrum
	RefundTo            *string `json:"refundTo,omitempty"`            // Arbitrum
	L1BaseFee           *string `json:"l1BaseFee,omitempty"`           // Arbitrum
	DepositValue        *string `json:"depositValue,omitempty"`        // Arbitrum
	RetryTo             *string `json:"retryTo,omitempty"`             // Arbitrum
	RetryValue          *string `json:"retryValue,omitempty"`          // Arbitrum
	RetryData           *string `json:"retryData,omitempty"`           // Arbitrum
	Beneficiary         *string `json:"beneficiary,omitempty"`         // Arbitrum
	MaxSubmissionFee    *string `json:"maxSubmissionFee,omitempty"`    // Arbitrum
	EffectiveGasPrice   *string `json:"effectiveGasPrice,omitempty"`   // Arbitrum

	QueueOrigin           *string `json:"queueOrigin,omitempty"`           // Optimism
	L1TxOrigin            *string `json:"l1TxOrigin,omitempty"`            // Optimism
	L1BlockTimestamp      *string `json:"l1Timestamp,omitempty"`           // Optimism
	Index                 *uint64 `json:"index,omitempty"`                 // Optimism
	QueueIndex            *uint64 `json:"queueIndex,omitempty"`            // Optimism
	SourceHash            *string `json:"sourceHash,omitempty"`            // Optimism
	Mint                  *string `json:"mint,omitempty"`                  // Optimism
	IsSystemTx            *bool   `json:"isSystemTx,omitempty"`            // Optimism
	DepositReceiptVersion *string `json:"depositReceiptVersion,omitempty"` // Optimism
}

type Access struct {
	Address     string   `json:"address" validate:"required"`
	StorageKeys []string `json:"storageKeys" validate:"required"`
}

type AccessListResp struct {
	AccessList []*Access `json:"accessList"`
	Error      *string   `json:"error,omitempty"`
	GasUsed    string    `json:"gasUsed"`
}

type Receipt struct {
	BlockHash         string  `json:"blockHash" validate:"required"`
	BlockNumber       uint64  `json:"blockNumber" validate:"-"`
	TransactionHash   string  `json:"transactionHash" validate:"required"`
	TransactionIndex  uint64  `json:"transactionIndex" validate:"-"`
	From              string  `json:"from" validate:"required"`
	To                string  `json:"to" validate:"required"`
	GasUsed           string  `json:"gasUsed" validate:"required"`
	CumulativeGasUsed string  `json:"cumulativeGasUsed" validate:"required"`
	ContractAddress   *string `json:"contractAddress,omitempty"`
	Logs              []*Log  `json:"logs" validate:"required"`
	Type              string  `json:"type" validate:"required"`
	EffectiveGasPrice string  `json:"effectiveGasPrice" validate:"required"`
	Root              *string `json:"root,omitempty"`
	Status            *string `json:"status,omitempty"`
	LogsBloom         string  `json:"logsBloom"`

	BlobGasPrice *string `json:"blobGasPrice,omitempty"` // EIP-4844
	BlobGasUsed  *string `json:"blobGasUsed,omitempty"`  // EIP-4844

	GasUsedForL1  *string `json:"gasUsedForL1,omitempty"`  // Arbitrum
	L1BlockNumber *uint64 `json:"l1BlockNumber,omitempty"` // Arbitrum

	L1GasPrice            *string `json:"l1GasPrice,omitempty"`            // Optimism
	L1GasUsed             *string `json:"l1GasUsed,omitempty"`             // Optimism
	L1FeeScalar           *string `json:"l1FeeScalar,omitempty"`           // Optimism
	L1Fee                 *string `json:"l1Fee,omitempty"`                 // Optimism
	DepositNonce          *string `json:"depositNonce,omitempty"`          // Optimism
	DepositReceiptVersion *string `json:"depositReceiptVersion,omitempty"` // Optimism
	L1BlobBaseFee         *string `json:"l1BlobBaseFee,omitempty"`         // Optimism
	L1BaseFeeScalar       *string `json:"l1BaseFeeScalar,omitempty"`       // Optimism
	L1BlobBaseFeeScalar   *string `json:"l1BlobBaseFeeScalar,omitempty"`   // Optimism
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

type LogFilter struct {
	BlockHash *string  `json:"blockHash,omitempty"`
	FromBlock *uint64  `json:"fromBlock,omitempty"`
	ToBlock   *uint64  `json:"toBlock,omitempty"`
	Address   *string  `json:"address,omitempty"`
	Topics    []string `json:"topics,omitempty"`
}

type FeeHistory struct {
	feeHistory
}

type SubLog struct {
	Address string   `json:"address,omitempty"`
	Topics  []string `json:"topics,omitempty"`
}

type debugTracer struct {
	TxHash string  `json:"txHash"`
	Error  *string `json:"error,omitempty"`
}

type DebugCallTracer struct {
	debugTracer
	Result *CallFrame `json:"result"`
}

type DebugFlatCallTracer struct {
	debugTracer
	Result *FlatCallFrame `json:"result"`
}

type CallFrame struct {
	From         string           `json:"from" validate:"required"`
	Gas          string           `json:"gas" validate:"required"`
	GasUsed      string           `json:"gasUsed" validate:"required"`
	To           *string          `json:"to,omitempty"`
	Input        *string          `json:"input"`
	Output       *string          `json:"output,omitempty"`
	Error        *string          `json:"error,omitempty"`
	RevertReason *string          `json:"revertReason,omitempty"`
	Calls        []*CallFrame     `json:"calls,omitempty"`
	Logs         []*CallLog       `json:"logs,omitempty"`
	Value        *decimal.Decimal `json:"value,omitempty"`
	Type         string           `json:"type"`
	Index        uint64           `json:"index"` // custom index
}

type CallLog struct {
	Address  string   `json:"address" validate:"required"`
	Topics   []string `json:"topics" validate:"required"`
	Data     string   `json:"data" validate:"required"`
	Position uint64   `json:"position" validate:"-"`
}

type FlatCallFrame struct {
	Action struct {
		Author         *string          `json:"author,omitempty"`
		RewardType     *string          `json:"rewardType,omitempty"`
		Address        *string          `json:"address,omitempty"`
		Balance        *string          `json:"balance,omitempty"`
		CreationMethod *string          `json:"creationMethod,omitempty"`
		RefundAddress  *string          `json:"refundAddress,omitempty"`
		CallType       *string          `json:"callType,omitempty"`
		From           *string          `json:"from,omitempty"`
		Gas            *string          `json:"gas,omitempty"`
		Input          *string          `json:"input,omitempty"`
		To             *string          `json:"to,omitempty"`
		Init           *string          `json:"init,omitempty"`
		Value          *decimal.Decimal `json:"value,omitempty"`
	} `json:"action"`
	BlockHash   string  `json:"blockHash" validate:"required"`
	BlockNumber uint64  `json:"blockNumber" validate:"-"`
	Error       *string `json:"error,omitempty"`
	Result      *struct {
		Address *string `json:"address,omitempty"`
		Code    *string `json:"code,omitempty"`
		GasUsed *string `json:"gasUsed,omitempty"`
		Output  *string `json:"output,omitempty"`
	} `json:"result,omitempty"`
	Subtraces           uint64   `json:"subtraces" validate:"required"`
	TraceAddress        []uint64 `json:"traceAddress" validate:"required"`
	TransactionHash     string   `json:"transactionHash" validate:"required"`
	TransactionPosition uint64   `json:"transactionPosition" validate:"-"`
	Type                string   `json:"type" validate:"required"`
	Index               uint64   `json:"index"` // custom index
}

// TODO: erigon(parity) trace types
type Trace struct {
	Action struct {
		From          string `json:"from"`
		CallType      string `json:"callType"`
		Gas           string `json:"gas"`
		Input         string `json:"input"`
		To            string `json:"to"`
		Value         string `json:"value"`
		Author        string `json:"author"`
		RewardType    string `json:"rewardType"`
		Init          string `json:"init"`
		Address       string `json:"address"`
		RefundAddress string `json:"refundAddress"`
		Balance       string `json:"balance"`
	} `json:"action"`
	BlockHash   string `json:"blockHash"`
	BlockNumber uint64 `json:"blockNumber"`
	Error       string `json:"error"`
	Result      *struct {
		Address string `json:"address"`
		GasUsed string `json:"gasUsed"`
		Output  string `json:"output"`
	} `json:"result"`
	Subtraces           uint64   `json:"subtraces"`
	TraceAddress        []uint64 `json:"traceAddress"`
	TransactionHash     string   `json:"transactionHash"`
	TransactionPosition uint64   `json:"transactionPosition"`
	Type                string   `json:"type"`
	Index               uint64   `json:"index"` // custom index
}

type ContractCreator struct {
	TransactionHash string `json:"hash"`
	Creator         string `json:"creator"`
}

type Balance struct {
	Address string          `json:"address"`
	Value   decimal.Decimal `json:"value"`
}

type QueryParams struct {
	To       string      `json:"to"`
	Data     string      `json:"data"`
	NumOrTag BlockAndTag `json:"-"`
}

type QueryResp struct {
	To     string
	Data   string
	Result string
	Error  error
}
