package evmc

import (
	"testing"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mainnetURL은 노트북에서만 접근 가능한 이더리움 메인넷 RPC 주소.
// 이 테스트들은 해당 RPC에 접근 가능한 환경에서만 실행된다.
const mainnetURL = "http://61.111.3.69:18014"

// Ethereum mainnet block 20,000,000 기준 레퍼런스 값
const (
	mainnetRefBlockNumber   = uint64(20_000_000)
	mainnetRefBlockHash     = "0xd24fd73f794058a3807db926d8898c6481e902b7edb91ce0d479d6760f276183"
	mainnetRefEip1559TxHash = "0xbb4b3fc2b746877dce70862850602f1d19bd890ab4db47e6b7ee1da1fe578a0d"
	mainnetRefBlobTxHash    = "0x0ff07f37baa7fa26bb7de3d3fc63002bf0acf3295bdab7f67c108c0d1a3bff15"
	mainnetRefUSDT          = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	mainnetRefMinerAddr     = "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
)

func mainnetClient(t *testing.T) *Evmc {
	t.Helper()
	c, err := New(mainnetURL)
	require.NoError(t, err, "mainnet client creation failed (is %s accessible?)", mainnetURL)
	return c
}

// ─── Block 조회 ────────────────────────────────────────────────────────────────

func Test_Mainnet_Eth_BlockNumber(t *testing.T) {
	c := mainnetClient(t)
	n, err := c.Eth().BlockNumber()
	require.NoError(t, err)
	assert.Greater(t, n, mainnetRefBlockNumber, "block number should be past block 20M")
}

func Test_Mainnet_Eth_GetBlockByNumber(t *testing.T) {
	c := mainnetClient(t)
	block, err := c.Eth().GetBlockByNumber(mainnetRefBlockNumber)
	require.NoError(t, err)

	assert.Equal(t, mainnetRefBlockNumber, block.Number)
	assert.Equal(t, mainnetRefBlockHash, block.Hash)
	assert.Equal(t, mainnetRefMinerAddr, block.Miner)
	assert.Equal(t, uint64(1717281407), block.Timestamp)
	assert.Len(t, block.Transactions, 134)
	// EIP-4844 필드 포함 확인
	require.NotNil(t, block.BlobGasUsed)
	assert.Equal(t, "0x20000", *block.BlobGasUsed)
}

func Test_Mainnet_Eth_GetBlockIncTxByNumber(t *testing.T) {
	c := mainnetClient(t)
	block, err := c.Eth().GetBlockIncTxByNumber(mainnetRefBlockNumber)
	require.NoError(t, err)

	assert.Equal(t, mainnetRefBlockNumber, block.Number)
	assert.Equal(t, mainnetRefBlockHash, block.Hash)
	require.Len(t, block.Transactions, 134)
	// 첫 번째 tx가 EIP-1559
	assert.Equal(t, "0x2", block.Transactions[0].Type)
	assert.Equal(t, mainnetRefEip1559TxHash, block.Transactions[0].Hash)
}

func Test_Mainnet_Eth_GetBlockByHash(t *testing.T) {
	c := mainnetClient(t)
	block, err := c.Eth().GetBlockByHash(mainnetRefBlockHash)
	require.NoError(t, err)

	assert.Equal(t, mainnetRefBlockNumber, block.Number)
	assert.Equal(t, mainnetRefBlockHash, block.Hash)
	assert.Len(t, block.Transactions, 134)
}

func Test_Mainnet_Eth_GetBlockIncTxByHash(t *testing.T) {
	c := mainnetClient(t)
	block, err := c.Eth().GetBlockIncTxByHash(mainnetRefBlockHash)
	require.NoError(t, err)

	assert.Equal(t, mainnetRefBlockNumber, block.Number)
	assert.Equal(t, mainnetRefBlockHash, block.Hash)
	require.Len(t, block.Transactions, 134)
	assert.Equal(t, "0x2", block.Transactions[0].Type)
}

func Test_Mainnet_Eth_GetBlockByTag(t *testing.T) {
	c := mainnetClient(t)
	block, err := c.Eth().GetBlockByTag(evmctypes.Latest)
	require.NoError(t, err)
	assert.Greater(t, block.Number, mainnetRefBlockNumber)
}

func Test_Mainnet_Eth_BlockTransactionCountByHash(t *testing.T) {
	c := mainnetClient(t)
	count, err := c.Eth().BlockTransactionCountByHash(mainnetRefBlockHash)
	require.NoError(t, err)
	assert.Equal(t, uint64(134), count)
}

func Test_Mainnet_Eth_BlockTransactionCountByNumber(t *testing.T) {
	c := mainnetClient(t)
	count, err := c.Eth().BlockTransactionCountByNumber(mainnetRefBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, uint64(134), count)
}

func Test_Mainnet_Eth_UncleCountByBlockHash(t *testing.T) {
	c := mainnetClient(t)
	count, err := c.Eth().UncleCountByBlockHash(mainnetRefBlockHash)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), count)
}

func Test_Mainnet_Eth_UncleCountByBlockNumber(t *testing.T) {
	c := mainnetClient(t)
	count, err := c.Eth().UncleCountByBlockNumber(mainnetRefBlockNumber)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), count)
}

func Test_Mainnet_Eth_GetBlockReceipts(t *testing.T) {
	c := mainnetClient(t)
	receipts, err := c.Eth().GetBlockReceipts(mainnetRefBlockNumber)
	require.NoError(t, err)

	assert.Len(t, receipts, 134)
	// 첫 번째: EIP-1559 타입
	assert.Equal(t, "0x2", receipts[0].Type)
	assert.Equal(t, mainnetRefEip1559TxHash, receipts[0].TransactionHash)
	assert.Equal(t, mainnetRefBlockNumber, receipts[0].BlockNumber)
	// blob tx receipt
	var blobReceipt *evmctypes.Receipt
	for _, r := range receipts {
		if r.Type == "0x3" {
			blobReceipt = r
			break
		}
	}
	require.NotNil(t, blobReceipt, "blob tx receipt not found")
	assert.Equal(t, mainnetRefBlobTxHash, blobReceipt.TransactionHash)
	require.NotNil(t, blobReceipt.BlobGasUsed)
	assert.Equal(t, "0x20000", *blobReceipt.BlobGasUsed)
	require.NotNil(t, blobReceipt.BlobGasPrice)
	assert.Equal(t, "0x1", *blobReceipt.BlobGasPrice)
}

func Test_Mainnet_Eth_GetBlockRange(t *testing.T) {
	c := mainnetClient(t)
	from := mainnetRefBlockNumber
	to := mainnetRefBlockNumber + 2
	blocks, err := c.Eth().GetBlockRange(from, to)
	require.NoError(t, err)

	assert.Len(t, blocks, 3)
	assert.Equal(t, from, blocks[0].Number)
	assert.Equal(t, to, blocks[2].Number)
}

func Test_Mainnet_Eth_GetBlockIncTxRange(t *testing.T) {
	c := mainnetClient(t)
	from := mainnetRefBlockNumber
	to := mainnetRefBlockNumber + 1
	blocks, err := c.Eth().GetBlockIncTxRange(from, to)
	require.NoError(t, err)

	assert.Len(t, blocks, 2)
	assert.Equal(t, from, blocks[0].Number)
	require.NotEmpty(t, blocks[0].Transactions)
	assert.Equal(t, mainnetRefEip1559TxHash, blocks[0].Transactions[0].Hash)
}

// ─── Transaction 조회 ──────────────────────────────────────────────────────────

func Test_Mainnet_Eth_GetTransactionByHash_EIP1559(t *testing.T) {
	c := mainnetClient(t)
	tx, err := c.Eth().GetTransactionByHash(mainnetRefEip1559TxHash)
	require.NoError(t, err)

	assert.Equal(t, "0x2", tx.Type)
	assert.Equal(t, mainnetRefEip1559TxHash, tx.Hash)
	assert.Equal(t, "0xae2fc483527b8ef99eb5d9b44875f005ba1fae13", tx.From)
	assert.Equal(t, mainnetRefBlockNumber, tx.BlockNumber)
	assert.Equal(t, mainnetRefBlockHash, tx.BlockHash)
	require.NotNil(t, tx.ChainID)
	assert.Equal(t, "0x1", *tx.ChainID)
	require.NotNil(t, tx.MaxFeePerGas)
	require.NotNil(t, tx.MaxPriorityFeePerGas)
	assert.Nil(t, tx.MaxFeePerBlobGas)
	assert.Empty(t, tx.BlobVersionedHashes)
}

func Test_Mainnet_Eth_GetTransactionByHash_BlobType3(t *testing.T) {
	c := mainnetClient(t)
	tx, err := c.Eth().GetTransactionByHash(mainnetRefBlobTxHash)
	require.NoError(t, err)

	assert.Equal(t, "0x3", tx.Type)
	assert.Equal(t, mainnetRefBlobTxHash, tx.Hash)
	assert.Equal(t, mainnetRefBlockNumber, tx.BlockNumber)
	require.NotNil(t, tx.MaxFeePerBlobGas)
	assert.Equal(t, "0x3b9aca00", *tx.MaxFeePerBlobGas)
	require.Len(t, tx.BlobVersionedHashes, 1)
	assert.Equal(t,
		"0x017ba4bd9c166498865a3d08618e333ee84812941b5c3a356971b4a6ffffa574",
		tx.BlobVersionedHashes[0])
}

func Test_Mainnet_Eth_GetTransactionByBlockHashAndIndex(t *testing.T) {
	c := mainnetClient(t)
	tx, err := c.Eth().GetTransactionByBlockHashAndIndex(mainnetRefBlockHash, 0)
	require.NoError(t, err)

	assert.Equal(t, mainnetRefEip1559TxHash, tx.Hash)
	assert.Equal(t, uint64(0), tx.TransactionIndex)
	assert.Equal(t, mainnetRefBlockHash, tx.BlockHash)
}

func Test_Mainnet_Eth_GetTransactionByBlockNumberAndIndex(t *testing.T) {
	c := mainnetClient(t)
	tx, err := c.Eth().GetTransactionByBlockNumberAndIndex(mainnetRefBlockNumber, 0)
	require.NoError(t, err)

	assert.Equal(t, mainnetRefEip1559TxHash, tx.Hash)
	assert.Equal(t, uint64(0), tx.TransactionIndex)
}

func Test_Mainnet_Eth_GetTransactionReceipt_EIP1559(t *testing.T) {
	c := mainnetClient(t)
	receipt, err := c.Eth().GetTransactionReceipt(mainnetRefEip1559TxHash)
	require.NoError(t, err)

	assert.Equal(t, "0x2", receipt.Type)
	assert.Equal(t, mainnetRefEip1559TxHash, receipt.TransactionHash)
	assert.Equal(t, mainnetRefBlockNumber, receipt.BlockNumber)
	assert.Equal(t, mainnetRefBlockHash, receipt.BlockHash)
	assert.Equal(t, "0xae2fc483527b8ef99eb5d9b44875f005ba1fae13", receipt.From)
	require.NotNil(t, receipt.Status)
	assert.Equal(t, "0x1", *receipt.Status)
	assert.Nil(t, receipt.BlobGasUsed)
	assert.Nil(t, receipt.BlobGasPrice)
}

func Test_Mainnet_Eth_GetTransactionReceipt_Blob(t *testing.T) {
	c := mainnetClient(t)
	receipt, err := c.Eth().GetTransactionReceipt(mainnetRefBlobTxHash)
	require.NoError(t, err)

	assert.Equal(t, "0x3", receipt.Type)
	assert.Equal(t, mainnetRefBlobTxHash, receipt.TransactionHash)
	require.NotNil(t, receipt.Status)
	assert.Equal(t, "0x1", *receipt.Status)
	require.NotNil(t, receipt.BlobGasUsed)
	assert.Equal(t, "0x20000", *receipt.BlobGasUsed)
	require.NotNil(t, receipt.BlobGasPrice)
	assert.Equal(t, "0x1", *receipt.BlobGasPrice)
}

// ─── Account / State 조회 ──────────────────────────────────────────────────────

func Test_Mainnet_Eth_GetBalance(t *testing.T) {
	c := mainnetClient(t)
	// 블록 20M 시점의 miner 잔액
	blockTag := evmctypes.FormatNumber(mainnetRefBlockNumber)
	balance, err := c.Eth().GetBalance(mainnetRefMinerAddr, blockTag)
	require.NoError(t, err)
	// 0 초과 정수여야 함
	assert.True(t, balance.IsPositive(), "balance should be positive")
	assert.False(t, balance.IsNegative())
}

func Test_Mainnet_Eth_GetTransactionCount(t *testing.T) {
	c := mainnetClient(t)
	blockTag := evmctypes.FormatNumber(mainnetRefBlockNumber)
	nonce, err := c.Eth().GetTransactionCount(mainnetRefMinerAddr, blockTag)
	require.NoError(t, err)
	assert.Greater(t, nonce, uint64(0))
}

func Test_Mainnet_Eth_GetCode(t *testing.T) {
	c := mainnetClient(t)
	blockTag := evmctypes.FormatNumber(mainnetRefBlockNumber)
	code, err := c.Eth().GetCode(mainnetRefUSDT, blockTag)
	require.NoError(t, err)
	// USDT는 컨트랙트이므로 코드가 존재해야 함
	assert.NotEmpty(t, code)
	assert.NotEqual(t, "0x", code)
}

func Test_Mainnet_Eth_GetStorageAt(t *testing.T) {
	c := mainnetClient(t)
	blockTag := evmctypes.FormatNumber(mainnetRefBlockNumber)
	// USDT slot 0: owner 관련
	value, err := c.Eth().GetStorageAt(mainnetRefUSDT, "0x0", blockTag)
	require.NoError(t, err)
	// 32바이트 hex 문자열 (0x + 64자)
	assert.Len(t, value, 66)
	assert.Equal(t, "0x", value[:2])
}

func Test_Mainnet_Eth_GetProof(t *testing.T) {
	c := mainnetClient(t)
	storageKey := "0x0000000000000000000000000000000000000000000000000000000000000001"
	proof, err := c.Eth().GetProof(mainnetRefUSDT, []string{storageKey}, evmctypes.Latest)
	require.NoError(t, err)

	// USDT 주소 (소문자)
	assert.Equal(t, "0xdac17f958d2ee523a2206206994597c13d831ec7", proof.Address)
	assert.NotEmpty(t, proof.AccountProof)
	assert.NotEmpty(t, proof.Balance)
	assert.NotEmpty(t, proof.CodeHash)
	assert.NotEmpty(t, proof.StorageHash)
	require.Len(t, proof.StorageProof, 1)
	assert.Equal(t, storageKey, proof.StorageProof[0].Key)
	assert.NotEmpty(t, proof.StorageProof[0].Proof)
}

// ─── Gas / Fee 조회 ────────────────────────────────────────────────────────────

func Test_Mainnet_Eth_GasPrice(t *testing.T) {
	c := mainnetClient(t)
	price, err := c.Eth().GasPrice()
	require.NoError(t, err)
	// 0 초과
	assert.True(t, price.IsPositive(), "gas price should be positive: %s", price)
}

func Test_Mainnet_Eth_MaxPriorityFeePerGas(t *testing.T) {
	c := mainnetClient(t)
	fee, err := c.Eth().MaxPriorityFeePerGas()
	require.NoError(t, err)
	assert.False(t, fee.IsNegative())
}

func Test_Mainnet_Eth_BlobBaseFee(t *testing.T) {
	c := mainnetClient(t)
	fee, err := c.Eth().BlobBaseFee()
	require.NoError(t, err)
	assert.False(t, fee.IsNegative())
}

func Test_Mainnet_Eth_FeeHistory(t *testing.T) {
	c := mainnetClient(t)
	fh, err := c.Eth().FeeHistory(4, evmctypes.FormatNumber(mainnetRefBlockNumber), []float64{25, 75})
	require.NoError(t, err)

	assert.Equal(t, uint64(19_999_997), fh.OldestBlock)
	// blockCount+1개의 baseFeePerGas (마지막은 예측값)
	assert.Len(t, fh.BaseFeePerGas, 5)
	assert.Len(t, fh.BaseFeePerBlobGas, 5)
	assert.Len(t, fh.GasUsedRatio, 4)
	assert.Len(t, fh.BlobGasUsedRatio, 4)
	// percentile reward: 25th, 75th
	assert.Len(t, fh.Reward, 4)
	for _, r := range fh.Reward {
		assert.Len(t, r, 2, "each reward entry should have 2 percentiles")
	}
}

// ─── Chain 정보 ────────────────────────────────────────────────────────────────

func Test_Mainnet_Eth_ChainID(t *testing.T) {
	c := mainnetClient(t)
	id, err := c.Eth().ChainID()
	require.NoError(t, err)
	assert.Equal(t, uint64(1), id, "Ethereum mainnet chainId should be 1")
}

func Test_Mainnet_Eth_Syncing(t *testing.T) {
	c := mainnetClient(t)
	syncing, detail, err := c.Eth().Syncing()
	require.NoError(t, err)
	// 동기화 완료된 노드이므로 false여야 함
	assert.False(t, syncing, "node should be fully synced")
	assert.Nil(t, detail)
}

// ─── Log 조회 ─────────────────────────────────────────────────────────────────

func Test_Mainnet_Eth_GetLogs(t *testing.T) {
	c := mainnetClient(t)
	addr := mainnetRefUSDT
	filter := &evmctypes.LogFilter{
		FromBlock: &[]uint64{mainnetRefBlockNumber}[0],
		ToBlock:   &[]uint64{mainnetRefBlockNumber}[0],
		Address:   &addr,
	}
	logs, err := c.Eth().GetLogs(filter)
	require.NoError(t, err)

	assert.Len(t, logs, 12)
	for _, l := range logs {
		assert.Equal(t, "0xdac17f958d2ee523a2206206994597c13d831ec7", l.Address)
		assert.Equal(t, mainnetRefBlockNumber, l.BlockNumber)
		assert.Equal(t, mainnetRefBlockHash, l.BlockHash)
		assert.NotEmpty(t, l.Topics)
	}
}

func Test_Mainnet_Eth_GetLogsByBlockNumber(t *testing.T) {
	c := mainnetClient(t)
	logs, err := c.Eth().GetLogsByBlockNumber(mainnetRefBlockNumber)
	require.NoError(t, err)

	// 블록 전체 로그
	assert.NotEmpty(t, logs)
	for _, l := range logs {
		assert.Equal(t, mainnetRefBlockNumber, l.BlockNumber)
		assert.Equal(t, mainnetRefBlockHash, l.BlockHash)
	}
}

func Test_Mainnet_Eth_GetLogsByBlockHash(t *testing.T) {
	c := mainnetClient(t)
	logs, err := c.Eth().GetLogsByBlockHash(mainnetRefBlockHash)
	require.NoError(t, err)

	assert.NotEmpty(t, logs)
	for _, l := range logs {
		assert.Equal(t, mainnetRefBlockHash, l.BlockHash)
		assert.Equal(t, mainnetRefBlockNumber, l.BlockNumber)
	}
}

// ─── eth_call / estimateGas / createAccessList ────────────────────────────────

func Test_Mainnet_Eth_Call(t *testing.T) {
	c := mainnetClient(t)
	// USDT totalSupply() 호출 (selector: 0x18160ddd)
	tx := &Tx{
		To:   mainnetRefUSDT,
		Data: "0x18160ddd",
	}
	result, err := c.Eth().Call(tx, evmctypes.FormatNumber(mainnetRefBlockNumber))
	require.NoError(t, err)

	// uint256 hex 결과 (0x + 64자)
	assert.NotEmpty(t, result)
	assert.Equal(t, "0x", result[:2])
}

func Test_Mainnet_Eth_EstimateGas(t *testing.T) {
	c := mainnetClient(t)
	// 단순 ETH 전송 gas 추정
	tx := &Tx{
		From:  mainnetRefMinerAddr,
		To:    "0x000000000000000000000000000000000000dead",
		Value: decimal.NewFromInt(1),
	}
	gas, err := c.Eth().EstimateGas(tx)
	require.NoError(t, err)
	// ETH 전송은 21000 gas
	assert.Equal(t, uint64(21000), gas)
}

func Test_Mainnet_Eth_CreateAccessList(t *testing.T) {
	c := mainnetClient(t)
	// USDT totalSupply() 에 대한 access list 생성
	tx := &Tx{
		From:     mainnetRefMinerAddr,
		To:       mainnetRefUSDT,
		Data:     "0x18160ddd",
		GasLimit: 100_000,
	}
	resp, err := c.Eth().CreateAccessList(tx)
	require.NoError(t, err)

	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.GasUsed)
	// access list는 nil이거나 빈 배열일 수 있음
}
