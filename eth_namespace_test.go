package evmc

import (
	"testing"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func testEvmc() *Evmc {
	rpcURL := "https://ethereum-mainnet.nodit.io/<api-key>"
	client, err := New(rpcURL)
	if err != nil {
		panic(err)
	}
	return client
}

// v1.14.5
// 2.59.0

func Test_ethNamespace(t *testing.T) {
	var (
		client          = testEvmc()
		testBlockNumber = uint64(20445681)
		testBlockHash   = "0xbd3c4cf74090e6a08285ae016df4c268220cb14fef1653de38348b5745956838"
		txCount         = 1163
	)
	t.Parallel()

	t.Log(client.nodeVersion)

	t.Run("eth_blobBaseFee", func(t *testing.T) {
		// TODO: skip this test if the node does not support the blobBaseFee method
		// t.SkipNow()
		baseFee, err := client.eth.BlobBaseFee()
		if err != nil {
			t.Error(err)
		} else {
			t.Log(baseFee)
			assert.True(t, true)
		}
	})

	t.Run("eth_chainId", func(t *testing.T) {
		chainID, err := client.eth.ChainID()
		if err != nil {
			t.Error(err)
		} else {
			assert.NotZero(t, chainID)
		}
	})

	t.Run("eth_blockNumber", func(t *testing.T) {
		latestBlockNumber, err := client.eth.BlockNumber()
		if err != nil {
			t.Error(err)
		} else {
			assert.NotZero(t, latestBlockNumber)
		}
	})

	t.Run("eth_getBlockByNumber", func(t *testing.T) {
		block, err := client.eth.GetBlockByNumber(testBlockNumber)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, testBlockNumber, block.Number)
			assert.Equal(t, txCount, len(block.Transactions))
		}
		block, err = client.eth.GetBlockByTag(evmctypes.Safe)
		if err != nil {
			t.Error(err)
		} else {
			assert.NotZero(t, block.Number)
		}
		blockIncTx, err := client.eth.GetBlockIncTxByNumber(testBlockNumber)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, testBlockNumber, blockIncTx.Number)
			assert.Equal(t, txCount, len(blockIncTx.Transactions))
		}
	})

	t.Run("eth_getBlockByHash", func(t *testing.T) {
		block, err := client.eth.GetBlockByHash(testBlockHash)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, testBlockHash, block.Hash)
			assert.Equal(t, txCount, len(block.Transactions))
		}
	})

	t.Run("eth_feeHistory", func(t *testing.T) {
		// get the fee history
		blockCount := uint64(1024)
		feeHistory, err := client.eth.FeeHistory(blockCount, evmctypes.FormatNumber(testBlockNumber), nil)
		if err != nil {
			t.Error(err)
		} else {
			wantOldestBlockNumber := testBlockNumber - blockCount + 1
			assert.Equal(t, wantOldestBlockNumber, feeHistory.OldestBlock)
		}
	})

	t.Run("eth_getBlockTranasactionCountByNumber", func(t *testing.T) {
		blockTxCount, err := client.eth.BlockTransactionCountByNumber(testBlockNumber)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, int(blockTxCount), txCount)
		}
	})

	t.Run("eth_getBlockTranasactionCountByHash", func(t *testing.T) {
		blockTxCount, err := client.eth.BlockTransactionCountByHash(testBlockHash)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, int(blockTxCount), txCount)
		}
	})

	t.Run("eth_getStorageAt", func(t *testing.T) {
		var (
			addr1 = "0x4Fabb145d64652a948d72533023f6E7A623C7C53"
			addr2 = "0xa3728dfd8b471e5ead8b6f08cc37d139b23fe5a9"
			key1  = "0x7050c9e0f4ca769c69bd3a8ef740bc37934f8e2c036e5a723fd8ee048ed3f8c3"
			key2  = "0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"
		)
		storage1, err := client.eth.GetStorageAt(addr1, key1, evmctypes.Latest)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, "0x0000000000000000000000002a3f1a37c04f82aa274f5353834b2d002db91015", storage1)
		}
		storage2, err := client.eth.GetStorageAt(addr2, key2, evmctypes.Latest)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, "0x0000000000000000000000008b142accdc5d4233abef308974aba3566e9a6824", storage2)
		}
	})

	t.Run("eth_getCode", func(t *testing.T) {
		addr := "0x8b142accdc5d4233abef308974aba3566e9a6824"
		code, err := client.eth.GetCode(addr, evmctypes.Latest)
		if err != nil {
			t.Error(err)
		} else {
			assert.NotEmpty(t, code)
		}
	})

	t.Run("eth_getTransactionByHash", func(t *testing.T) {
		txHash := "0xd31e242b8cbc5d32a4db2d528c0b9695e2a022100fd4f992dde8fc48ce262348"
		tx, err := client.eth.GetTransactionByHash(txHash)
		if err != nil {
			t.Error(err)
		} else {
			assert.NotEmpty(t, tx)
		}
	})

	t.Run("eth_getTransactionReceipt", func(t *testing.T) {
		txHash := "0xd31e242b8cbc5d32a4db2d528c0b9695e2a022100fd4f992dde8fc48ce262348"
		receipt, err := client.eth.GetTransactionReceipt(txHash)
		if err != nil {
			t.Error(err)
		} else {
			assert.NotEmpty(t, receipt)
		}
	})

	t.Run("eth_getTransactionCount", func(t *testing.T) {
		addr := "0xA627202EC3A92e647da14A54F7973478b7cFAe4c"
		nonce, err := client.eth.GetTransactionCount(addr, evmctypes.Pending)
		if err != nil {
			t.Error(err)
		} else {
			assert.NotZero(t, nonce)
		}
	})

	t.Run("eth_getBalance", func(t *testing.T) {
		addr := "0xA627202EC3A92e647da14A54F7973478b7cFAe4c"
		balance, err := client.eth.GetBalance(addr, evmctypes.Latest)
		if err != nil {
			t.Error(err)
		} else {
			assert.NotZero(t, balance)
		}
	})

	t.Run("eth_getLogs", func(t *testing.T) {
		var (
			fromBlock uint64 = 20411245
			toBlock   uint64 = 20411256
			addr             = "0xA3728dFd8B471e5EAD8b6f08cc37d139B23Fe5A9"
		)
		logs, err := client.eth.GetLogs(&evmctypes.LogFilter{
			FromBlock: &fromBlock,
			ToBlock:   &toBlock,
			Address:   &addr,
		})
		if err != nil {
			t.Error(err)
		} else {
			assert.NotEmpty(t, logs)
			assert.Equal(t, 2, len(logs))
		}
	})

	t.Run("eth_getBlockReceipts", func(t *testing.T) {
		blockReceipts, err := client.eth.GetBlockReceipts(testBlockNumber)
		if err != nil {
			t.Error(err)
		} else {
			assert.Equal(t, txCount, len(blockReceipts))
			assert.NotEmpty(t, blockReceipts)
		}
	})

	t.Run("eth_gasPrice", func(t *testing.T) {
		gasPrice, err := client.eth.GasPrice()
		if err != nil {
			t.Error(err)
		} else {
			assert.NotZero(t, gasPrice)
		}
	})

	t.Run("eth_MaxPriorityFeePerGas", func(t *testing.T) {
		maxPriorityFeePerGas, err := client.eth.MaxPriorityFeePerGas()
		if err != nil {
			t.Error(err)
		} else {
			assert.NotZero(t, maxPriorityFeePerGas)
		}
	})

	t.Run("eth_estimateGas", func(t *testing.T) {
		tx := &Tx{
			From:  ZeroAddress,
			To:    ZeroAddress,
			Nonce: 0,
			Data:  "0x",
			Value: decimal.Zero,
		}
		gas, err := client.eth.EstimateGas(tx)
		if err != nil {
			t.Error(err)
		} else {
			assert.NotZero(t, gas)
		}
	})
}
