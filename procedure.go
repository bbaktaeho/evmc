package evmc

type Procedure string

const (
	Web3ClientVersion Procedure = "web3_clientVersion"

	EthCall                          Procedure = "eth_call"
	EthChainID                       Procedure = "eth_chainId"
	EthGetCode                       Procedure = "eth_getCode"
	EthBlockNumber                   Procedure = "eth_blockNumber"
	EthGetBlockByNumber              Procedure = "eth_getBlockByNumber"
	EthGetBlockByHash                Procedure = "eth_getBlockByHash"
	EthGetUncleByBlockNumberAndIndex Procedure = "eth_getUncleByBlockNumberAndIndex"
	EthGetTransaction                Procedure = "eth_getTransactionByHash"
	EthGetReceipt                    Procedure = "eth_getTransactionReceipt"
	EthGetBalance                    Procedure = "eth_getBalance"
	EthGetStorageAt                  Procedure = "eth_getStorageAt"
	EthGetBlockReceipts              Procedure = "eth_getBlockReceipts"
	EthGetLogs                       Procedure = "eth_getLogs"

	EthGetTransactionReceiptsByBlock Procedure = "eth_getTransactionReceiptsByBlock" // bor

	DebugTraceBlockByNumber Procedure = "debug_traceBlockByNumber"
	DebugTraceTransaction   Procedure = "debug_traceTransaction"

	OtsGetContractCreator Procedure = "ots_getContractCreator" // erigon
	TraceBlock            Procedure = "trace_block"            // erigon

	// arb_trace methods on the Arbitrum One chain should be called on blocks prior to 22207815
	ArbitraceBlock Procedure = "arbtrace_block" // arbitrum
)

func (p Procedure) String() string {
	return string(p)
}

type BlockTag string

const (
	Latest    BlockTag = "latest"
	Safe      BlockTag = "safe"
	Finalized BlockTag = "finalized"
)

func (bt BlockTag) String() string {
	switch bt {
	case Latest:
		return "latest"
	case Safe:
		return "safe"
	case Finalized:
		return "finalized"
	default:
		return ""
	}
}
