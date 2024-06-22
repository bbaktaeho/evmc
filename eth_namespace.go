package evmc

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	ethChainID                       = "eth_chainId"
	ethGetBlockNumber                = "eth_blockNumber"
	ethGetBlockByNumber              = "eth_getBlockByNumber"
	ethCall                          = "eth_call"
	ethGetCode                       = "eth_getCode"
	ethBlockNumber                   = "eth_blockNumber"
	ethGetBlockByHash                = "eth_getBlockByHash"
	ethGetUncleByBlockNumberAndIndex = "eth_getUncleByBlockNumberAndIndex"
	ethGetTransaction                = "eth_getTransactionByHash"
	ethGetReceipt                    = "eth_getTransactionReceipt"
	ethGetBalance                    = "eth_getBalance"
	ethGetStorageAt                  = "eth_getStorageAt"
	ethGetBlockReceipts              = "eth_getBlockReceipts"
	ethGetLogs                       = "eth_getLogs"
	ethGetTransactionReceiptsByBlock = "eth_getTransactionReceiptsByBlock" // bor
)

type ethNamespace struct {
	c caller
}

func (e *ethNamespace) GetChainID() (uint64, error) {
	result := new(string)
	if err := e.c.call(context.Background(), result, ethChainID); err != nil {
		return 0, err
	}
	id, err := hexutil.DecodeUint64(*result)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (e *ethNamespace) GetblockNumber() uint64 {
	return 0
}
