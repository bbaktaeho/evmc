package evmc

// TODO: omitempty settings
// TODO: valid tag settings

type WithTxs interface{ []string | []*Transaction }

type Block[T WithTxs] struct {
	ParentHash       string      `json:"parentHash"`
	Sha3Uncles       string      `json:"sha3Uncles"`
	Miner            string      `json:"miner"`
	StateRoot        string      `json:"stateRoot"`
	TransactionsRoot string      `json:"transactionsRoot"`
	ReceiptsRoot     string      `json:"receiptsRoot"`
	LogsBloom        string      `json:"logsBloom"`
	Difficulty       string      `json:"difficulty"`
	Number           string      `json:"number"`
	GasLimit         string      `json:"gasLimit"`
	GasUsed          string      `json:"gasUsed"`
	Timestamp        string      `json:"timestamp"`
	ExtraData        string      `json:"extraData"`
	MixHash          string      `json:"mixHash"`
	Nonce            string      `json:"nonce"`
	Hash             string      `json:"hash"`
	Size             string      `json:"size"`
	TotalDifficulty  string      `json:"totalDifficulty"`
	Transactions     T           `json:"transactions"`
	Uncles           []string    `json:"uncles"`
	UncleBlocks      *[]Block[T] `json:"uncleBlocks,omitempty"`

	BaseFeePerGas *string `json:"baseFeePerGas,omitempty"` // EIP-1559

	WithdrawalsRoot *string       `json:"withdrawalsRoot,omitempty"` // EIP-4895
	Withdrawals     *[]Withdrawal `json:"withdrawals,omitempty"`     // EIP-4895

	BlobGasUsed   *string `json:"blobGasUsed,omitempty"`   // EIP-4844
	ExcessBlobGas *string `json:"excessBlobGas,omitempty"` // EIP-4844

	ParentBeaconBlockRoot *string `json:"parentBeaconBlockRoot,omitempty"` // EIP-4788

	*additionalArbitrumBlock // Arbitrum
}

type additionalArbitrumBlock struct {
	L1BlockNumber *string `json:"l1BlockNumber,omitempty"`
	SendCount     *string `json:"sendCount,omitempty"`
	SendRoot      *string `json:"sendRoot,omitempty"`
}

type Withdrawal struct {
	Index          string `json:"index"`
	ValidatorIndex string `json:"validatorIndex"`
	Address        string `json:"address"`
	Amount         string `json:"amount"`
}

type Transaction struct {
	BlockHash        string  `json:"blockHash"`
	BlockNumber      string  `json:"blockNumber"`
	From             string  `json:"from"`
	Gas              string  `json:"gas"`
	GasPrice         string  `json:"gasPrice"`
	Hash             string  `json:"hash"`
	Input            string  `json:"input"`
	Nonce            string  `json:"nonce"`
	To               string  `json:"to"`
	TransactionIndex string  `json:"transactionIndex"`
	Value            string  `json:"value"`
	Type             string  `json:"type"`
	V                *string `json:"v"`
	R                string  `json:"r"`
	S                string  `json:"s"`

	ChainID *string `json:"chainId,omitempty"` // EIP-155

	AccessList *[]Access `json:"accessList,omitempty"` // EIP-2930

	MaxFeePerGas         *string `json:"maxFeePerGas,omitempty"`         // EIP-1559
	MaxPriorityFeePerGas *string `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559

	MaxFeePerBlobGas    *string   `json:"maxFeePerBlobGas,omitempty"`    // EIP-4844
	BlobVersionedHashes *[]string `json:"blobVersionedHashes,omitempty"` // EIP-4844

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
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type Receipt struct {
	Type              string `json:"type"`
	Root              string `json:"root"`
	Status            string `json:"status"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	LogsBloom         string `json:"logsBloom"`
	Logs              []*Log `json:"logs"`
	TransactionHash   string `json:"transactionHash"`
	ContractAddress   string `json:"contractAddress"`
	GasUsed           string `json:"gasUsed"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	TransactionIndex  string `json:"transactionIndex"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	From              string `json:"from"`
	To                string `json:"to"`

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
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

type callTracer interface {
	*CallTracer | []*FlatCallTracer
}

type DebugTrace[T callTracer] struct {
	TxHash string `json:"txHash"`
	Result T      `json:"result"`
}

type CallTracer struct {
	From         string       `json:"from"`
	Gas          string       `json:"gas"`
	GasUsed      string       `json:"gasUsed"`
	To           string       `json:"to,omitempty"`
	Input        string       `json:"input"`
	Calls        []CallTracer `json:"calls,omitempty"`
	Output       string       `json:"output,omitempty"`
	Type         string       `json:"type"`
	Value        string       `json:"value,omitempty"`
	Error        string       `json:"error,omitempty"`
	RevertReason string       `json:"revertReason,omitempty"`
}

type FlatCallTracer struct {
	Action struct {
		Author         string `json:"author,omitempty"`
		RewardType     string `json:"rewardType,omitempty"`
		Address        string `json:"address,omitempty"`
		Balance        string `json:"balance,omitempty"`
		CreationMethod string `json:"creationMethod,omitempty"`
		RefundAddress  string `json:"refundAddress,omitempty"`
		CallType       string `json:"callType,omitempty"`
		From           string `json:"from,omitempty"`
		Gas            string `json:"gas,omitempty"`
		Input          string `json:"input,omitempty"`
		To             string `json:"to,omitempty"`
		Init           string `json:"init,omitempty"`
		Value          string `json:"value,omitempty"`
	} `json:"action"`
	BlockHash   string `json:"blockHash"`
	BlockNumber uint64 `json:"blockNumber"`
	Error       string `json:"error,omitempty"`
	Result      *struct {
		Address string `json:"address,omitempty"`
		Code    string `json:"code,omitempty"`
		GasUsed string `json:"gasUsed,omitempty"`
		Output  string `json:"output,omitempty"`
	} `json:"result,omitempty"`
	Subtraces           uint64   `json:"subtraces"`
	TraceAddress        []uint64 `json:"traceAddress"`
	TransactionHash     string   `json:"transactionHash"`
	TransactionPosition uint64   `json:"transactionPosition"`
	Type                string   `json:"type"`
}

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
