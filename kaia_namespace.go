package evmc

// https://github.com/kaiachain/kaia-sdk/tree/dev/web3rpc/rpc-specs

import (
	"context"
	"fmt"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/bbaktaeho/evmc/evmctypes/kaiatypes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type kaiaNamespace struct {
	c caller
}

func (k *kaiaNamespace) GetBlockIncTxRange(from, to uint64) ([]*kaiatypes.BlockIncTx, error) {
	return k.GetBlockIncTxRangeWithContext(context.Background(), from, to)
}

func (k *kaiaNamespace) GetBlockIncTxRangeWithContext(ctx context.Context, from, to uint64) ([]*kaiatypes.BlockIncTx, error) {
	return k.getBlockIncTxRange(ctx, from, to)
}

func (k *kaiaNamespace) getBlockIncTxRange(ctx context.Context, from, to uint64) ([]*kaiatypes.BlockIncTx, error) {
	if from > to {
		return nil, ErrInvalidRange
	}
	var (
		size     = to - from + 1
		results  = make([]*kaiatypes.BlockIncTx, size)
		elements = make([]rpc.BatchElem, size)
	)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: KaiaGetBlockByNumber.String(),
			Args:   []interface{}{evmctypes.FormatNumber(from + uint64(i)), true},
			Result: &results[i],
		}
	}
	if err := k.c.BatchCallWithContext(ctx, elements, -1); err != nil {
		return nil, err
	}
	for i, el := range elements {
		if el.Error != nil {
			return nil, el.Error
		}
		if results[i] == nil || results[i].Hash == "" {
			return nil, fmt.Errorf("block %d not found", from+uint64(i))
		}
	}
	return results, nil
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
func (k *kaiaNamespace) GetBlockByHash(hash string) (*kaiatypes.Block, error) {
	return k.GetBlockByHashWithContext(context.Background(), hash)
}

// GetBlockByHashWithContext returns information about a block by hash.
func (k *kaiaNamespace) GetBlockByHashWithContext(ctx context.Context, hash string) (*kaiatypes.Block, error) {
	block := new(kaiatypes.Block)
	if err := k.getBlockByHash(ctx, block, hash, false); err != nil {
		return nil, err
	}
	return block, nil
}

// GetBlockByHashIncTx returns information about a block by hash, including transactions.
func (k *kaiaNamespace) GetBlockByHashIncTx(hash string) (*kaiatypes.BlockIncTx, error) {
	return k.GetBlockByHashIncTxWithContext(context.Background(), hash)
}

// GetBlockByHashIncTxWithContext returns information about a block by hash, including transactions.
func (k *kaiaNamespace) GetBlockByHashIncTxWithContext(ctx context.Context, hash string) (*kaiatypes.BlockIncTx, error) {
	block := new(kaiatypes.BlockIncTx)
	if err := k.getBlockByHash(ctx, block, hash, true); err != nil {
		return nil, err
	}
	return block, nil
}

func (k *kaiaNamespace) getBlockByHash(ctx context.Context, result interface{}, hash string, incTx bool) error {
	params := []interface{}{hash, incTx}
	if err := k.c.call(ctx, result, KaiaGetBlockByHash, params...); err != nil {
		return err
	}
	return nil
}

// GetBlockByNumber returns information about a block by block number.
func (k *kaiaNamespace) GetBlockByNumber(blockNumber uint64) (*kaiatypes.Block, error) {
	return k.GetBlockByNumberWithContext(context.Background(), blockNumber)
}

// GetBlockByNumberWithContext returns information about a block by block number.
func (k *kaiaNamespace) GetBlockByNumberWithContext(ctx context.Context, blockNumber uint64) (*kaiatypes.Block, error) {
	block := new(kaiatypes.Block)
	if err := k.getBlockByNumber(ctx, block, evmctypes.FormatNumber(blockNumber), false); err != nil {
		return nil, err
	}
	return block, nil
}

// GetBlockByNumberIncTx returns information about a block by block number, including transactions.
func (k *kaiaNamespace) GetBlockByNumberIncTx(blockNumber uint64) (*kaiatypes.BlockIncTx, error) {
	return k.GetBlockByNumberIncTxWithContext(context.Background(), blockNumber)
}

// GetBlockByNumberIncTxWithContext returns information about a block by block number, including transactions.
func (k *kaiaNamespace) GetBlockByNumberIncTxWithContext(ctx context.Context, blockNumber uint64) (*kaiatypes.BlockIncTx, error) {
	block := new(kaiatypes.BlockIncTx)
	if err := k.getBlockByNumber(ctx, block, evmctypes.FormatNumber(blockNumber), true); err != nil {
		return nil, err
	}
	return block, nil
}

func (k *kaiaNamespace) getBlockByNumber(ctx context.Context, result interface{}, number evmctypes.BlockAndTag, incTx bool) error {
	params := []interface{}{number, incTx}
	if err := k.c.call(ctx, result, KaiaGetBlockByNumber, params...); err != nil {
		return err
	}
	return nil
}

func (k *kaiaNamespace) GetTransactionByHash(hash string) (*kaiatypes.Transaction, error) {
	return k.GetTransactionByHashWithContext(context.Background(), hash)
}

func (k *kaiaNamespace) GetTransactionByHashWithContext(ctx context.Context, hash string) (*kaiatypes.Transaction, error) {
	result := new(kaiatypes.Transaction)
	if err := k.c.call(ctx, result, KaiaGetTransactionByHash, hash); err != nil {
		return nil, err
	}
	if result.Hash == "" {
		return nil, fmt.Errorf("transaction %s not found", hash)
	}
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
	var result = new([]*kaiatypes.Receipt)
	if err := k.c.call(ctx, result, KaiaGetBlockReceipts, hexutil.EncodeUint64(blockNumber)); err != nil {
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

// GetRewardsRange returns the rewards for a range of blocks.
func (k *kaiaNamespace) GetRewardsRange(from, to uint64) ([]*kaiatypes.Rewards, error) {
	return k.GetRewardsRangeWithContext(context.Background(), from, to)
}

// GetRewardsRangeWithContext returns the rewards for a range of blocks.
func (k *kaiaNamespace) GetRewardsRangeWithContext(ctx context.Context, from, to uint64) ([]*kaiatypes.Rewards, error) {
	if from > to {
		return nil, ErrInvalidRange
	}
	var (
		size     = to - from + 1
		results  = make([]*kaiatypes.Rewards, size)
		elements = make([]rpc.BatchElem, size)
	)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: KaiaGetRewards.String(),
			Args:   []interface{}{evmctypes.FormatNumber(from + uint64(i))},
			Result: &results[i],
		}
	}
	if err := k.c.BatchCallWithContext(ctx, elements, -1); err != nil {
		return nil, err
	}
	for i, el := range elements {
		if el.Error != nil {
			return nil, el.Error
		}
		if results[i] == nil {
			return nil, fmt.Errorf("rewards for block %d not found", from+uint64(i))
		}
	}
	return results, nil
}

// GetRewards returns the rewards for a block by block number.
func (k *kaiaNamespace) GetRewards(blockNumber uint64) (*kaiatypes.Rewards, error) {
	return k.GetRewardsWithContext(context.Background(), blockNumber)
}

func (k *kaiaNamespace) GetRewardsWithContext(ctx context.Context, blockNumber uint64) (*kaiatypes.Rewards, error) {
	return k.getRewards(ctx, blockNumber)
}

func (k *kaiaNamespace) getRewards(ctx context.Context, blockNumber uint64) (*kaiatypes.Rewards, error) {
	var result kaiatypes.Rewards
	hexBlock := hexutil.EncodeUint64(blockNumber)
	if err := k.c.call(ctx, &result, KaiaGetRewards, hexBlock); err != nil {
		return nil, err
	}
	return &result, nil
}
