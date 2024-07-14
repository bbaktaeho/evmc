package evmc

type subscription string

const (
	newHeads               subscription = "newHeads"
	newPendingTransactions subscription = "newPendingTransactions"
	logs                   subscription = "logs"
)

type procedure string

const (
	web3ClientVersion procedure = "web3_clientVersion"

	ethBlobBaseFee                      procedure = "eth_blobBaseFee"
	ethBlockNumber                      procedure = "eth_blockNumber"
	ethCall                             procedure = "eth_call"
	ethChainID                          procedure = "eth_chainId"
	ethCreateAccessList                 procedure = "eth_createAccessList"
	ethEstimateGas                      procedure = "eth_estimateGas"
	ethFeeHistory                       procedure = "eth_feeHistory"
	ethGasPrice                         procedure = "eth_gasPrice"
	ethGetBalance                       procedure = "eth_getBalance"
	ethGetBlockByHash                   procedure = "eth_getBlockByHash"
	ethGetBlockByNumber                 procedure = "eth_getBlockByNumber"
	ethGetBlockReceipts                 procedure = "eth_getBlockReceipts"
	ethGetBlockTransactionCountByHash   procedure = "eth_getBlockTransactionCountByHash"
	ethGetBlockTransactionCountByNumber procedure = "eth_getBlockTransactionCountByNumber"
	ethGetCode                          procedure = "eth_getCode"
	ethGetFilterChanges                 procedure = "eth_getFilterChanges"
	ethGetReceipt                       procedure = "eth_getTransactionReceipt"
	ethSendRawTransaction               procedure = "eth_sendRawTransaction"
	ethMaxPriorityFeePerGas             procedure = "eth_maxPriorityFeePerGas"
	ethSyncing                          procedure = "eth_syncing"
	ethGetUncleByBlockNumberAndIndex    procedure = "eth_getUncleByBlockNumberAndIndex"
	ethGetTransaction                   procedure = "eth_getTransactionByHash"
	ethGetStorageAt                     procedure = "eth_getStorageAt"
	ethGetLogs                          procedure = "eth_getLogs"
	ethGetTransactionCount              procedure = "eth_getTransactionCount"
	ethGetTransactionReceiptsByBlock    procedure = "eth_getTransactionReceiptsByBlock" // bor

	debugTraceBlockByNumber procedure = "debug_traceBlockByNumber"
	debugTraceTransaction   procedure = "debug_traceTransaction"

	otsGetContractCreator procedure = "ots_getContractCreator" // erigon
	traceBlock            procedure = "trace_block"            // erigon

	// arb_trace methods on the Arbitrum One chain should be called on blocks prior to 22207815
	arbitraceBlock procedure = "arbtrace_block" // arbitrum
)

func (p procedure) String() string {
	return string(p)
}
