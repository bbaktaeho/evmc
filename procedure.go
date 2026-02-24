package evmc

type subscription string

const (
	newHeads               subscription = "newHeads"
	newPendingTransactions subscription = "newPendingTransactions"
	logs                   subscription = "logs"
)

type Procedure string

func (p Procedure) String() string {
	return string(p)
}

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
	EthGetTransactionByHash                    Procedure = "eth_getTransactionByHash"
	EthGetTransactionByBlockHashAndIndex       Procedure = "eth_getTransactionByBlockHashAndIndex"
	EthGetTransactionByBlockNumberAndIndex     Procedure = "eth_getTransactionByBlockNumberAndIndex"
	EthGetUncleCountByBlockHash                Procedure = "eth_getUncleCountByBlockHash"
	EthGetUncleCountByBlockNumber              Procedure = "eth_getUncleCountByBlockNumber"
	EthGetProof                                Procedure = "eth_getProof"
	EthSimulateV1                              Procedure = "eth_simulateV1"
	EthGetStorageAt                            Procedure = "eth_getStorageAt"
	EthGetLogs                          Procedure = "eth_getLogs"
	EthGetTransactionCount              Procedure = "eth_getTransactionCount"
	EthGetTransactionReceiptsByBlock    Procedure = "eth_getTransactionReceiptsByBlock" // bor

	DebugTraceBlockByNumber Procedure = "debug_traceBlockByNumber"
	DebugTraceBlockByHash   Procedure = "debug_traceBlockByHash"
	DebugTraceTransaction   Procedure = "debug_traceTransaction"
	DebugTraceCall          Procedure = "debug_traceCall"

	DebugGetRawHeader      Procedure = "debug_getRawHeader"
	DebugGetRawBlock       Procedure = "debug_getRawBlock"
	DebugGetRawTransaction Procedure = "debug_getRawTransaction"
	DebugGetRawReceipts    Procedure = "debug_getRawReceipts"
	DebugGetBadBlocks      Procedure = "debug_getBadBlocks"

	OtsGetContractCreator Procedure = "ots_getContractCreator" // erigon
	TraceBlock            Procedure = "trace_block"            // erigon

	// arb_trace methods on the Arbitrum One chain should be called on blocks prior to 22207815
	ArbitraceBlock Procedure = "arbtrace_block" // arbitrum
)

const (
	KaiaBlockNumber                            Procedure = "kaia_blockNumber"
	KaiaGetBlockByHash                         Procedure = "kaia_getBlockByHash"
	KaiaGetBlockByNumber                       Procedure = "kaia_getBlockByNumber"
	KaiaGetBlockReceipts                       Procedure = "kaia_getBlockReceipts"
	KaiaGetBlockTransactionCountByHash         Procedure = "kaia_getBlockTransactionCountByHash"
	KaiaGetBlockTransactionCountByNumber       Procedure = "kaia_getBlockTransactionCountByNumber"
	KaiaGetBlockWithConsensusInfoByHash        Procedure = "kaia_getBlockWithConsensusInfoByHash"
	KaiaGetBlockWithConsensusInfoByNumber      Procedure = "kaia_getBlockWithConsensusInfoByNumber"
	KaiaGetBlockWithConsensusInfoByNumberRange Procedure = "kaia_getBlockWithConsensusInfoByNumberRange"
	KaiaGetCommittee                           Procedure = "kaia_getCommittee"
	KaiaGetCommitteeSize                       Procedure = "kaia_getCommitteeSize"
	KaiaGetCouncil                             Procedure = "kaia_getCouncil"
	KaiaGetCouncilSize                         Procedure = "kaia_getCouncilSize"
	KaiaGetHeaderByHash                        Procedure = "kaia_getHeaderByHash"
	KaiaGetHeaderByNumber                      Procedure = "kaia_getHeaderByNumber"
	KaiaGetRewards                             Procedure = "kaia_getRewards"
	KaiaGetStorageAt                           Procedure = "kaia_getStorageAt"
	KaiaSyncing                                Procedure = "kaia_syncing"
	KaiaCall                                   Procedure = "kaia_call"
	KaiaCreateAccessList                       Procedure = "kaia_createAccessList"
	KaiaEstimateComputationCost                Procedure = "kaia_estimateComputationCost"
	KaiaEstimateGas                            Procedure = "kaia_estimateGas"
	KaiaGetDecodedAnchoringTransactionByHash   Procedure = "kaia_getDecodedAnchoringTransactionByHash"
	KaiaGetRawTransactionByBlockHashAndIndex   Procedure = "kaia_getRawTransactionByBlockHashAndIndex"
	KaiaGetRawTransactionByBlockNumberAndIndex Procedure = "kaia_getRawTransactionByBlockNumberAndIndex"
	KaiaGetRawTransactionByHash                Procedure = "kaia_getRawTransactionByHash"
	KaiaGetTransactionByBlockHashAndIndex      Procedure = "kaia_getTransactionByBlockHashAndIndex"
	KaiaGetTransactionByBlockNumberAndIndex    Procedure = "kaia_getTransactionByBlockNumberAndIndex"
	KaiaGetTransactionByHash                   Procedure = "kaia_getTransactionByHash"
	KaiaGetTransactionBySenderTxHash           Procedure = "kaia_getTransactionBySenderTxHash"
	KaiaGetTransactionReceipt                  Procedure = "kaia_getTransactionReceipt"
	KaiaGetTransactionReceiptBySenderTxHash    Procedure = "kaia_getTransactionReceiptBySenderTxHash"
	KaiaPendingTransactions                    Procedure = "kaia_pendingTransactions"
	KaiaResend                                 Procedure = "kaia_resend"
	KaiaSendRawTransaction                     Procedure = "kaia_sendRawTransaction"
	KaiaSendTransaction                        Procedure = "kaia_sendTransaction"
	KaiaSendTransactionAsFeePayer              Procedure = "kaia_sendTransactionAsFeePayer"
	KaiaSignTransaction                        Procedure = "kaia_signTransaction"
	KaiaSignTransactionAsFeePayer              Procedure = "kaia_signTransactionAsFeePayer"
)
