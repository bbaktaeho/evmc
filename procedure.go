package evmc

type subscription string

const (
	newHeads               subscription = "newHeads"
	newPendingTransactions subscription = "newPendingTransactions"
	logs                   subscription = "logs"
)

type Procedure string

const (
	Web3ClientVersion Procedure = "web3_clientVersion"

	EthNewBlockFilter              Procedure = "eth_newBlockFilter"
	EthNewPendingTransactionFilter Procedure = "eth_newPendingTransactionFilter"
	EthNewFilter                   Procedure = "eth_newFilter"
	EthUninstallFilter             Procedure = "eth_uninstallFilter"
	EthGetFilterChanges            Procedure = "eth_getFilterChanges"
	EthGetFilterLogs               Procedure = "eth_getFilterLogs"

	EthGetUncleByBlockNumberAndIndex Procedure = "eth_getUncleByBlockNumberAndIndex"

	EthBlobBaseFee                      Procedure = "eth_blobBaseFee"
	EthBlockNumber                      Procedure = "eth_blockNumber"
	EthCall                             Procedure = "eth_call"
	EthChainID                          Procedure = "eth_chainId"
	EthCreateAccessList                 Procedure = "eth_createAccessList"
	EthEstimateGas                      Procedure = "eth_estimateGas"
	EthFeeHistory                       Procedure = "eth_feeHistory"
	EthGasPrice                         Procedure = "eth_gasPrice"
	EthGetBalance                       Procedure = "eth_getBalance"
	EthGetBlockByHash                   Procedure = "eth_getBlockByHash"
	EthGetBlockByNumber                 Procedure = "eth_getBlockByNumber"
	EthGetBlockReceipts                 Procedure = "eth_getBlockReceipts"
	EthGetBlockTransactionCountByHash   Procedure = "eth_getBlockTransactionCountByHash"
	EthGetBlockTransactionCountByNumber Procedure = "eth_getBlockTransactionCountByNumber"
	EthGetCode                          Procedure = "eth_getCode"
	EthGetReceipt                       Procedure = "eth_getTransactionReceipt"
	EthSendRawTransaction               Procedure = "eth_sendRawTransaction"
	EthMaxPriorityFeePerGas             Procedure = "eth_maxPriorityFeePerGas"
	EthSyncing                          Procedure = "eth_syncing"
	EthGetTransactionByHash             Procedure = "eth_getTransactionByHash"
	EthGetStorageAt                     Procedure = "eth_getStorageAt"
	EthGetLogs                          Procedure = "eth_getLogs"
	EthGetTransactionCount              Procedure = "eth_getTransactionCount"
	EthGetTransactionReceiptsByBlock    Procedure = "eth_getTransactionReceiptsByBlock" // bor

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
