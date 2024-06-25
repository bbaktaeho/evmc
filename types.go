package evmc

type incTxs interface{ []string | []*Transaction }

type Block[T incTxs] struct {
	Number           string     `json:"number" validate:"required"`
	Hash             string     `json:"hash" validate:"required"`
	ParentHash       string     `json:"parentHash" validate:"required"`
	Nonce            string     `json:"nonce" validate:"required"`
	MixHash          string     `json:"mixHash" validate:"required"`
	Sha3Uncles       string     `json:"sha3Uncles" validate:"required"`
	LogsBloom        string     `json:"logsBloom" validate:"required"`
	StateRoot        string     `json:"stateRoot" validate:"required"`
	Miner            string     `json:"miner" validate:"required"`
	Difficulty       string     `json:"difficulty" validate:"required"`
	ExtraData        string     `json:"extraData" validate:"required"`
	GasLimit         string     `json:"gasLimit" validate:"required"`
	GasUsed          string     `json:"gasUsed" validate:"required"`
	Timestamp        string     `json:"timestamp" validate:"required"`
	TransactionsRoot string     `json:"transactionsRoot" validate:"required"`
	ReceiptsRoot     string     `json:"receiptsRoot" validate:"required"`
	Transactions     T          `json:"transactions" validate:"required"`
	TotalDifficulty  string     `json:"totalDifficulty" validate:"required"`
	Size             string     `json:"size"`
	Uncles           []string   `json:"uncles"`
	UncleBlocks      []Block[T] `json:"uncleBlocks,omitempty"`

	BaseFeePerGas         *string      `json:"baseFeePerGas,omitempty"`         // EIP-1559
	WithdrawalsRoot       *string      `json:"withdrawalsRoot,omitempty"`       // EIP-4895
	Withdrawals           []Withdrawal `json:"withdrawals,omitempty"`           // EIP-4895
	BlobGasUsed           *string      `json:"blobGasUsed,omitempty"`           // EIP-4844
	ExcessBlobGas         *string      `json:"excessBlobGas,omitempty"`         // EIP-4844
	ParentBeaconBlockRoot *string      `json:"parentBeaconBlockRoot,omitempty"` // EIP-4788

	*additionalArbitrumBlock // Arbitrum
}

type additionalArbitrumBlock struct {
	L1BlockNumber *string `json:"l1BlockNumber,omitempty"`
	SendCount     *string `json:"sendCount,omitempty"`
	SendRoot      *string `json:"sendRoot,omitempty"`
}

type Withdrawal struct {
	Index          string `json:"index" validate:"required"`
	ValidatorIndex string `json:"validatorIndex" validate:"required"`
	Address        string `json:"address" validate:"required"`
	Amount         string `json:"amount" validate:"required"`
}

type Transaction struct {
	BlockHash        string  `json:"blockHash" validate:"required"`
	BlockNumber      string  `json:"blockNumber" validate:"required"`
	From             string  `json:"from" validate:"required"`
	Gas              string  `json:"gas" validate:"required"`
	GasPrice         string  `json:"gasPrice" validate:"required"`
	Hash             string  `json:"hash" validate:"required"`
	Input            string  `json:"input" validate:"required"`
	Nonce            string  `json:"nonce" validate:"required"`
	To               string  `json:"to" validate:"required"`
	TransactionIndex string  `json:"transactionIndex" validate:"required"`
	Value            string  `json:"value" validate:"required"`
	Type             string  `json:"type" validate:"required"`
	V                string  `json:"v" validate:"required"`
	R                string  `json:"r" validate:"required"`
	S                string  `json:"s" validate:"required"`
	YParity          *string `json:"yParity,omitempty"`

	ChainID              *string   `json:"chainId,omitempty"`              // EIP-155
	AccessList           []*Access `json:"accessList,omitempty"`           // EIP-2930
	MaxFeePerGas         *string   `json:"maxFeePerGas,omitempty"`         // EIP-1559
	MaxPriorityFeePerGas *string   `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559
	MaxFeePerBlobGas     *string   `json:"maxFeePerBlobGas,omitempty"`     // EIP-4844
	BlobVersionedHashes  []string  `json:"blobVersionedHashes,omitempty"`  // EIP-4844

	*additionalArbitrumTx // Arbitrum
	*additionalOptimismTx // Optimism
}

type additionalArbitrumTx struct {
	RequestID           *string `json:"requestId,omitempty"`
	TicketID            *string `json:"ticketId,omitempty"`
	MaxRefund           *string `json:"maxRefund,omitempty"`
	SubmissionFeeRefund *string `json:"submissionFeeRefund,omitempty"`
	RefundTo            *string `json:"refundTo,omitempty"`
	L1BaseFee           *string `json:"l1BaseFee,omitempty"`
	L1BlockNumber       *string `json:"l1BlockNumber,omitempty"`
	DepositValue        *string `json:"depositValue,omitempty"`
	RetryTo             *string `json:"retryTo,omitempty"`
	RetryValue          *string `json:"retryValue,omitempty"`
	RetryData           *string `json:"retryData,omitempty"`
	Beneficiary         *string `json:"beneficiary,omitempty"`
	MaxSubmissionFee    *string `json:"maxSubmissionFee,omitempty"`
	EffectiveGasPrice   *string `json:"effectiveGasPrice,omitempty"`
}

type additionalOptimismTx struct {
	QueueOrigin           *string `json:"queueOrigin,omitempty"`
	L1TxOrigin            *string `json:"l1TxOrigin,omitempty"`
	L1BlockTimestamp      *string `json:"l1Timestamp,omitempty"`
	L1BlockNumber         *string `json:"l1BlockNumber,omitempty"`
	Index                 *string `json:"index,omitempty"`
	QueueIndex            *string `json:"queueIndex,omitempty"`
	SourceHash            *string `json:"sourceHash,omitempty"`
	Mint                  *string `json:"mint,omitempty"`
	IsSystemTx            *bool   `json:"isSystemTx,omitempty"`
	DepositReceiptVersion *string `json:"depositReceiptVersion,omitempty"`
}

type Access struct {
	Address     string   `json:"address" validate:"required"`
	StorageKeys []string `json:"storageKeys" validate:"required"`
}

type Receipt struct {
	BlockHash         string  `json:"blockHash" validate:"required"`
	BlockNumber       string  `json:"blockNumber" validate:"required"`
	TransactionHash   string  `json:"transactionHash" validate:"required"`
	TransactionIndex  string  `json:"transactionIndex" validate:"required"`
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

	*additionalArbitrumReceipt // Arbitrum
	*additionalOptimismReceipt // Optimism
}

type additionalArbitrumReceipt struct {
	GasUsedForL1  *string `json:"gasUsedForL1,omitempty"`
	L1BlockNumber *string `json:"l1BlockNumber,omitempty"`
}

type additionalOptimismReceipt struct {
	L1GasPrice            *string `json:"l1GasPrice,omitempty"`
	L1GasUsed             *string `json:"l1GasUsed,omitempty"`
	L1FeeScalar           *string `json:"l1FeeScalar,omitempty"`
	L1Fee                 *string `json:"l1Fee,omitempty"`
	DepositNonce          *string `json:"depositNonce,omitempty"`
	DepositReceiptVersion *string `json:"depositReceiptVersion,omitempty"`
	L1BlobBaseFee         *string `json:"l1BlobBaseFee,omitempty"`
	L1BaseFeeScalar       *string `json:"l1BaseFeeScalar,omitempty"`
	L1BlobBaseFeeScalar   *string `json:"l1BlobBaseFeeScalar,omitempty"`
}

type Log struct {
	Address          string   `json:"address" validate:"required"`
	Topics           []string `json:"topics" validate:"required"`
	Data             string   `json:"data" validate:"required"`
	BlockNumber      string   `json:"blockNumber" validate:"required"`
	TransactionHash  string   `json:"transactionHash" validate:"required"`
	TransactionIndex string   `json:"transactionIndex" validate:"required"`
	BlockHash        string   `json:"blockHash" validate:"required"`
	LogIndex         string   `json:"logIndex" validate:"required"`
	Removed          bool     `json:"removed" validate:"-"`
}

type LogFilter struct {
	BlockHash *string  `json:"blockHash,omitempty"`
	FromBlock *uint64  `json:"fromBlock,omitempty"`
	ToBlock   *uint64  `json:"toBlock,omitempty"`
	Address   *string  `json:"address,omitempty"`
	Topics    []string `json:"topics,omitempty"`
}

type callTracer interface {
	*CallFrame | []*FlatCallFrame
}

type DebugCallTracer[T callTracer] struct {
	TxHash string  `json:"txHash"`
	Result T       `json:"result"`
	Error  *string `json:"error,omitempty"`
}

type CallFrame struct {
	From         string      `json:"from" validate:"required"`
	Gas          string      `json:"gas" validate:"required"`
	GasUsed      string      `json:"gasUsed" validate:"required"`
	To           *string     `json:"to,omitempty"`
	Input        *string     `json:"input"`
	Output       *string     `json:"output,omitempty"`
	Error        *string     `json:"error,omitempty"`
	RevertReason *string     `json:"revertReason,omitempty"`
	Calls        []CallFrame `json:"calls,omitempty"`
	Logs         []CallLog   `json:"logs,omitempty"`
	Value        *string     `json:"value,omitempty"`
	Type         string      `json:"type"`
}

type CallLog struct {
	Address  string   `json:"address" validate:"required"`
	Topics   []string `json:"topics" validate:"required"`
	Data     string   `json:"data" validate:"required"`
	Position string   `json:"position" validate:"required"`
}

type FlatCallFrame struct {
	Action struct {
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
}

type ContractCreator struct {
	TransactionHash string `json:"hash"`
	Creator         string `json:"creator"`
}
