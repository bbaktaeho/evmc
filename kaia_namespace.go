package evmc

import (
	"context"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/bbaktaeho/evmc/evmctypes/kaiatypes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	// "github.com/shopspring/decimal" // 필요시 사용
)

type kaiaNamespace struct {
	c caller
}

// BlockNumber returns the current block number.
func (k *kaiaNamespace) BlockNumber() (uint64, error) {
	return k.BlockNumberWithContext(context.Background())
}

// BlockNumberWithContext returns the current block number.
func (k *kaiaNamespace) BlockNumberWithContext(ctx context.Context) (uint64, error) {
	return k.blockNumber(ctx)
}

func (k *kaiaNamespace) blockNumber(ctx context.Context) (uint64, error) {
	result := new(string)
	if err := k.c.call(ctx, result, KaiaBlockNumber); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

// GetBlockByHash returns information about a block by hash.
func (k *kaiaNamespace) GetBlockByHash(hash string, incTx bool) (*kaiatypes.Block, error) {
	return k.GetBlockByHashWithContext(context.Background(), hash, incTx)
}

// GetBlockByHashWithContext returns information about a block by hash.
func (k *kaiaNamespace) GetBlockByHashWithContext(ctx context.Context, hash string, incTx bool) (*kaiatypes.Block, error) {
	return k.getBlockByHash(ctx, hash, incTx)
}

func (k *kaiaNamespace) getBlockByHash(ctx context.Context, hash string, incTx bool) (*kaiatypes.Block, error) {
	result := new(kaiatypes.Block) // Kaia의 Block 타입이 Ethereum과 동일하다고 가정, 다르면 수정 필요
	params := []interface{}{hash, incTx}
	if err := k.c.call(ctx, result, KaiaGetBlockByHash, params...); err != nil {
		return nil, err
	}
	// Ethereum과 마찬가지로 Uncle block 처리가 필요하면 여기에 추가
	// if len(result.Uncles) > 0 {
	// 	uncleBlocks, err := k.getUncleBlocks(ctx, result.Number, result.Uncles) // getUncleBlocks 구현 필요
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	result.UncleBlocks = uncleBlocks
	// }
	return result, nil
}

// GetBlockByNumber returns information about a block by block number.
func (k *kaiaNamespace) GetBlockByNumber(number uint64, incTx bool) (*kaiatypes.Block, error) {
	return k.GetBlockByNumberWithContext(context.Background(), number, incTx)
}

// GetBlockByNumberWithContext returns information about a block by block number.
func (k *kaiaNamespace) GetBlockByNumberWithContext(ctx context.Context, number uint64, incTx bool) (*kaiatypes.Block, error) {
	return k.getBlockByNumber(ctx, number, incTx)
}

func (k *kaiaNamespace) getBlockByNumber(ctx context.Context, number uint64, incTx bool) (*kaiatypes.Block, error) {
	result := new(kaiatypes.Block) // Kaia의 Block 타입이 Ethereum과 동일하다고 가정, 다르면 수정 필요
	params := []interface{}{evmctypes.FormatNumber(number), incTx}
	if err := k.c.call(ctx, result, KaiaGetBlockByNumber, params...); err != nil {
		return nil, err
	}
	// Ethereum과 마찬가지로 Uncle block 처리가 필요하면 여기에 추가
	return result, nil
}

// GetBlockReceipts returns the receipts for all transactions in a block.
func (k *kaiaNamespace) GetBlockReceipts(blockNumber uint64) ([]*kaiatypes.Receipt, error) {
	return k.GetBlockReceiptsWithContext(context.Background(), blockNumber)
}

// GetBlockReceiptsWithContext returns the receipts for all transactions in a block.
func (k *kaiaNamespace) GetBlockReceiptsWithContext(ctx context.Context, blockNumber uint64) ([]*kaiatypes.Receipt, error) {
	return k.getBlockReceipts(ctx, blockNumber)
}

func (k *kaiaNamespace) getBlockReceipts(ctx context.Context, blockNumber uint64) ([]*kaiatypes.Receipt, error) {
	var result = new([]*kaiatypes.Receipt) // Kaia의 Receipt 타입이 Ethereum과 동일하다고 가정
	if err := k.c.call(ctx, result, KaiaGetBlockReceipts, evmctypes.FormatNumber(blockNumber)); err != nil {
		return nil, err
	}
	return *result, nil
}

// GetTransactionReceipt returns the receipt for a transaction by transaction hash.
func (k *kaiaNamespace) GetTransactionReceipt(hash string) (*kaiatypes.Receipt, error) {
	return k.GetTransactionReceiptWithContext(context.Background(), hash)
}

// GetTransactionReceiptWithContext returns the receipt for a transaction by transaction hash.
func (k *kaiaNamespace) GetTransactionReceiptWithContext(ctx context.Context, hash string) (*kaiatypes.Receipt, error) {
	return k.getTransactionReceipt(ctx, hash)
}

func (k *kaiaNamespace) getTransactionReceipt(ctx context.Context, hash string) (*kaiatypes.Receipt, error) {
	result := new(kaiatypes.Receipt)
	if err := k.c.call(ctx, result, KaiaGetTransactionReceipt, hash); err != nil {
		return nil, err
	}
	return result, nil
}

func (k *kaiaNamespace) GetRewards(blockNumber uint64) (*kaiatypes.Rewards, error) {
	return k.GetRewardsWithContext(context.Background(), blockNumber)
}

func (k *kaiaNamespace) GetRewardsWithContext(ctx context.Context, blockNumber uint64) (*kaiatypes.Rewards, error) {
	return k.getRewards(ctx, blockNumber)
}

// getRewards 공식 스펙: https://github.com/kaiachain/kaia-sdk/blob/dev/web3rpc/rpc-specs/paths/kaia/block/getRewards.yaml
func (k *kaiaNamespace) getRewards(ctx context.Context, blockNumber uint64) (*kaiatypes.Rewards, error) {
	var result kaiatypes.Rewards
	hexBlock := hexutil.EncodeUint64(blockNumber)
	if err := k.c.call(ctx, &result, KaiaGetRewards, hexBlock); err != nil {
		return nil, err
	}
	return &result, nil
}
