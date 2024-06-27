package evmc

type procedure string

const (
	web3ClientVersion procedure = "web3_clientVersion"

	ethSendRawTransaction            procedure = "eth_sendRawTransaction"
	ethMaxPriorityFeePerGas          procedure = "eth_maxPriorityFeePerGas"
	ethSyncing                       procedure = "eth_syncing"
	ethGasPrice                      procedure = "eth_gasPrice"
	ethFeeHistory                    procedure = "eth_feeHistory"
	ethCall                          procedure = "eth_call"
	ethChainID                       procedure = "eth_chainId"
	ethGetCode                       procedure = "eth_getCode"
	ethBlockNumber                   procedure = "eth_blockNumber"
	ethGetBlockByNumber              procedure = "eth_getBlockByNumber"
	ethGetBlockByHash                procedure = "eth_getBlockByHash"
	ethGetUncleByBlockNumberAndIndex procedure = "eth_getUncleByBlockNumberAndIndex"
	ethGetTransaction                procedure = "eth_getTransactionByHash"
	ethGetReceipt                    procedure = "eth_getTransactionReceipt"
	ethGetBalance                    procedure = "eth_getBalance"
	ethGetStorageAt                  procedure = "eth_getStorageAt"
	ethGetLogs                       procedure = "eth_getLogs"
	ethGetTransactionCount           procedure = "eth_getTransactionCount"
	ethGetBlockReceipts              procedure = "eth_getBlockReceipts"
	ethGetTransactionReceiptsByBlock procedure = "eth_getTransactionReceiptsByBlock" // bor

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
