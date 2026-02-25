package evmctypes

import (
	"encoding/json"

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
	UncleBlocks  []*Block       `json:"uncleBlocks,omitempty"`
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

	ChainID              *string          `json:"chainId,omitempty"`              // EIP-155
	AccessList           []*Access        `json:"accessList,omitempty"`           // EIP-2930
	MaxFeePerGas         *string          `json:"maxFeePerGas,omitempty"`         // EIP-1559
	MaxPriorityFeePerGas *string          `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559
	MaxFeePerBlobGas     *string          `json:"maxFeePerBlobGas,omitempty"`     // EIP-4844
	BlobVersionedHashes  []string         `json:"blobVersionedHashes,omitempty"`  // EIP-4844
	AuthorizationList    []*Authorization `json:"authorizationList,omitempty"`    // EIP-7702

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

type Authorization struct {
	ChainID string `json:"chainId" validate:"required"`
	Nonce   string `json:"nonce" validate:"required"`
	Address string `json:"address" validate:"required"`
	YParity string `json:"yParity" validate:"required"`
	R       string `json:"r" validate:"required"`
	S       string `json:"s" validate:"required"`
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

	BlockTimestamp uint64 `json:"blockTimestamp" validate:"-"`
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

type Syncing struct {
	syncing
}

type SubLog struct {
	Address string   `json:"address,omitempty"`
	Topics  []string `json:"topics,omitempty"`
}

type defaultTraceResult struct {
	TxHash string      `json:"txHash"`
	Error  interface{} `json:"error"` // geth is string, erigon is map[string]interface{}
}

type TraceResult struct {
	defaultTraceResult
	Result interface{} `json:"result,omitempty"`
}

type CallTracer struct {
	defaultTraceResult
	Result *CallFrame `json:"result,omitempty"`
}

type FlatCallTracer struct {
	defaultTraceResult
	Result []*FlatCallFrame `json:"result,omitempty"`
}

type PrestateTracer struct {
	defaultTraceResult
	Result json.RawMessage `json:"result,omitempty"`
}

type CustomTraceResult struct {
	defaultTraceResult
	Result json.RawMessage `json:"result,omitempty"`
}

type PrestateResult struct {
	json.RawMessage
}

func (p *PrestateResult) ParseFrame() (PrestateFrame, error) {
	if len(p.RawMessage) == 0 || string(p.RawMessage) == "null" {
		return nil, nil
	}
	var frame PrestateFrame
	if err := json.Unmarshal(p.RawMessage, &frame); err != nil {
		return nil, err
	}
	return frame, nil
}

func (p *PrestateResult) ParseDiffFrame() (*PrestateDiffFrame, error) {
	if len(p.RawMessage) == 0 || string(p.RawMessage) == "null" {
		return nil, nil
	}
	var diffFrame PrestateDiffFrame
	if err := json.Unmarshal(p.RawMessage, &diffFrame); err != nil {
		return nil, err
	}
	return &diffFrame, nil
}

// ParseFrames parses the result as PrestateFrame
func (p *PrestateTracer) ParseFrames() (PrestateFrame, error) {
	if len(p.Result) == 0 || string(p.Result) == "null" {
		return nil, nil
	}
	var frame PrestateFrame
	if err := json.Unmarshal(p.Result, &frame); err != nil {
		return nil, err
	}
	return frame, nil
}

func (p *PrestateTracer) ParseDiffFrames() (*PrestateDiffFrame, error) {
	if len(p.Result) == 0 || string(p.Result) == "null" {
		return nil, nil
	}
	var diffFrame PrestateDiffFrame
	if err := json.Unmarshal(p.Result, &diffFrame); err != nil {
		return nil, err
	}
	diffFrame.TxHash = p.TxHash
	return &diffFrame, nil
}

// TODO: WIP
// DefaultFrame is a response in default tracer.
//
//nolint:unused
type DefaultFrame struct {
	Depth   int         `json:"depth"`
	Gas     uint64      `json:"gas"`
	GasCost uint64      `json:"gasCost"`
	Memory  []string    `json:"memory"`
	OP      string      `json:"op"`
	PC      uint64      `json:"pc"`
	Stack   []string    `json:"stack"`
	Error   interface{} `json:"error"`
	Storage interface{} `json:"storage"`
}

// nolint:unused
type FourByteFrame map[string]uint64

// Arbitrum
type EVMTransfer struct {
	From    *string `json:"from"`
	To      *string `json:"to"`
	Value   string  `json:"value"`
	Purpose string  `json:"purpose"`
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

	// Arbitrum
	AfterEVMTransfers  []*EVMTransfer `json:"afterEVMTransfers,omitempty"`
	BeforeEVMTransfers []*EVMTransfer `json:"beforeEVMTransfers,omitempty"`
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
	BlockNumber uint64  `json:"blockNumber" validate:"required"`
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

	// Arbitrum
	AfterEVMTransfers  []*EVMTransfer `json:"afterEVMTransfers,omitempty"`
	BeforeEVMTransfers []*EVMTransfer `json:"beforeEVMTransfers,omitempty"`
}

type PrestateFrame map[string]*PrestateAccount

type PrestateDiffFrame struct {
	TxHash string        `json:"txHash,omitempty"` // Only by BlockByNumber
	Pre    PrestateFrame `json:"pre,omitempty"`
	Post   PrestateFrame `json:"post,omitempty"`
}

type PrestateAccount struct {
	Balance  *decimal.Decimal  `json:"balance,omitempty"`
	Code     string            `json:"code,omitempty"`
	CodeHash string            `json:"codeHash,omitempty"`
	Nonce    uint64            `json:"nonce,omitempty"`
	Storage  map[string]string `json:"storage,omitempty"`
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

// SimulateBlockOverride defines fields to override in simulated block headers.
type SimulateBlockOverride struct {
	BlockNumber  *string `json:"number,omitempty"`
	Time         *string `json:"time,omitempty"`
	Gas          *string `json:"gas,omitempty"`
	FeeRecipient *string `json:"feeRecipient,omitempty"`
	PrevRandao   *string `json:"prevRandao,omitempty"`
	BaseFeePerGas *string `json:"baseFeePerGas,omitempty"`
	BlobBaseFee  *string `json:"blobBaseFee,omitempty"`
}

// SimulateStateOverride defines per-account state overrides for simulation.
type SimulateStateOverride map[string]*SimulateAccountOverride

// SimulateAccountOverride overrides a single account's state for simulation.
type SimulateAccountOverride struct {
	Balance   *string           `json:"balance,omitempty"`
	Nonce     *uint64           `json:"nonce,omitempty"`
	Code      *string           `json:"code,omitempty"`
	State     map[string]string `json:"state,omitempty"`
	StateDiff map[string]string `json:"stateDiff,omitempty"`
}

// SimulateCall represents a single call to simulate within a block.
type SimulateCall struct {
	From                 *string `json:"from,omitempty"`
	To                   *string `json:"to,omitempty"`
	Gas                  *string `json:"gas,omitempty"`
	GasPrice             *string `json:"gasPrice,omitempty"`
	MaxFeePerGas         *string `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *string `json:"maxPriorityFeePerGas,omitempty"`
	Value                *string `json:"value,omitempty"`
	Input                *string `json:"input,omitempty"`
	Nonce                *string `json:"nonce,omitempty"`
}

// BlockStateCall defines calls and overrides for a single simulated block.
type BlockStateCall struct {
	BlockOverrides *SimulateBlockOverride `json:"blockOverrides,omitempty"`
	StateOverrides SimulateStateOverride  `json:"stateOverrides,omitempty"`
	Calls          []*SimulateCall        `json:"calls"`
}

// SimulatePayload is the input for eth_simulateV1.
type SimulatePayload struct {
	BlockStateCalls       []*BlockStateCall `json:"blockStateCalls"`
	TraceTransfers        bool              `json:"traceTransfers,omitempty"`
	Validation            bool              `json:"validation,omitempty"`
	ReturnFullTransactions bool             `json:"returnFullTransactions,omitempty"`
}

// CallError holds error information from a simulated call.
type CallError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// SimulateCallResult holds the result of a single simulated call.
type SimulateCallResult struct {
	ReturnValue string     `json:"returnData"`
	Logs        []*Log     `json:"logs"`
	GasUsed     uint64     `json:"gasUsed"`
	Status      uint64     `json:"status"`
	Error       *CallError `json:"error,omitempty"`
}

// SimulateBlockResult holds the result of a single simulated block.
type SimulateBlockResult struct {
	block
	// FeeRecipient is the block producer address (maps to the "miner" JSON field).
	FeeRecipient string `json:"miner,omitempty"`
	// BaseFeePerGas shadows block.BaseFeePerGas to expose as a plain string.
	BaseFeePerGas string `json:"baseFeePerGas,omitempty"`
	Calls         []*SimulateCallResult `json:"calls"`
	Transactions  interface{}           `json:"transactions"`
	Withdrawals   []*Withdrawal         `json:"withdrawals,omitempty"`
}

type StorageProof struct {
	Key   string   `json:"key"`
	Value string   `json:"value"`
	Proof []string `json:"proof"`
}

type AccountProof struct {
	Address      string          `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      string          `json:"balance"`
	CodeHash     string          `json:"codeHash"`
	Nonce        string          `json:"nonce"`
	StorageHash  string          `json:"storageHash"`
	StorageProof []*StorageProof `json:"storageProof"`
}

type BadBlock struct {
	Block *Block `json:"block"`
	Hash  string `json:"hash"`
	RLP   string `json:"rlp"`
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
