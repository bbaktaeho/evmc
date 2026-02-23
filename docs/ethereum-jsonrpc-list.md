## Ethereum JSONRPC List

Reference: https://github.com/ethereum/execution-apis/tree/main/src

Legend: [x] implemented, [ ] not implemented

### eth namespace

- [x] eth_blobBaseFee
- [x] eth_blockNumber
- [x] eth_call
- [x] eth_chainId
- [x] eth_createAccessList
- [x] eth_estimateGas
- [x] eth_feeHistory
- [x] eth_gasPrice
- [x] eth_getBalance
- [x] eth_getBlockByHash
- [x] eth_getBlockByNumber
- [x] eth_getBlockReceipts
- [x] eth_getBlockTransactionCountByHash
- [x] eth_getBlockTransactionCountByNumber
- [x] eth_getCode
- [x] eth_getFilterChanges (Procedure constant defined, HTTP filter not implemented)
- [x] eth_getFilterLogs (Procedure constant defined, HTTP filter not implemented)
- [x] eth_getLogs
- [x] eth_getProof
- [x] eth_getStorageAt
- [x] eth_getTransactionByBlockHashAndIndex
- [x] eth_getTransactionByBlockNumberAndIndex
- [x] eth_getTransactionByHash
- [x] eth_getTransactionCount
- [x] eth_getTransactionReceipt
- [x] eth_getUncleCountByBlockHash
- [x] eth_getUncleCountByBlockNumber
- [x] eth_maxPriorityFeePerGas
- [x] eth_newBlockFilter (Procedure constant defined, HTTP filter not implemented)
- [x] eth_newFilter (Procedure constant defined, HTTP filter not implemented)
- [x] eth_newPendingTransactionFilter (Procedure constant defined, HTTP filter not implemented)
- [x] eth_sendRawTransaction
- [x] eth_syncing
- [x] eth_uninstallFilter (Procedure constant defined, HTTP filter not implemented)
- [ ] eth_accounts (local wallet method, out of scope)
- [ ] eth_coinbase (local wallet method, out of scope)
- [ ] eth_sign (local wallet method, out of scope)
- [ ] eth_signTransaction (local wallet method, out of scope)
- [ ] eth_sendTransaction (local wallet method, requires node key management)
- [x] eth_simulateV1

### debug namespace

- [x] debug_getBadBlocks
- [x] debug_getRawBlock
- [x] debug_getRawHeader
- [x] debug_getRawReceipts
- [x] debug_getRawTransaction
- [x] debug_traceBlockByHash
- [x] debug_traceBlockByNumber
- [x] debug_traceTransaction
