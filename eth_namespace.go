package evmc

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// TODO: get uncle block
// TODO: batch call

type ethNamespace struct {
	c caller
}

func (e *ethNamespace) ChainID() (uint64, error) {
	return e.chainID(context.Background())
}

func (e *ethNamespace) ChainIDWithContext(ctx context.Context) (uint64, error) {
	return e.chainID(ctx)
}

func (e *ethNamespace) chainID(ctx context.Context) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethChainID); err != nil {
		return 0, err
	}
	id, err := hexutil.DecodeUint64(*result)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (e *ethNamespace) GetStorageAt(address, position string, numOrTag interface{}) (string, error) {
	return e.getStorageAt(context.Background(), address, position, numOrTag)
}

func (e *ethNamespace) GetStorageAtWithContext(
	ctx context.Context,
	address string,
	position string,
	numOrTag interface{},
) (string, error) {
	return e.getStorageAt(ctx, address, position, numOrTag)
}

func (e *ethNamespace) getStorageAt(
	ctx context.Context,
	address string,
	position string,
	numOrTag interface{},
) (string, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numOrTag)
	if err != nil {
		return "", err
	}
	if err := e.c.call(
		ctx,
		result,
		ethGetStorageAt,
		address,
		position,
		parsedNumOrTag,
	); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *ethNamespace) GetBlockNumber() (uint64, error) {
	return e.getBlockNumber(context.Background())
}

func (e *ethNamespace) GetBlockNumberWithContext(ctx context.Context) (uint64, error) {
	return e.getBlockNumber(ctx)
}

func (e *ethNamespace) getBlockNumber(ctx context.Context) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethBlockNumber); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) GetCode(address string, numOrTag interface{}) (string, error) {
	return e.getCode(context.Background(), address, numOrTag)
}

func (e *ethNamespace) GetCodeWithContext(ctx context.Context, address string, numOrTag interface{}) (string, error) {
	return e.getCode(ctx, address, numOrTag)
}

func (e *ethNamespace) getCode(
	ctx context.Context,
	address string,
	numOrTag interface{},
) (string, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numOrTag)
	if err != nil {
		return "", err
	}
	if err := e.c.call(ctx, result, ethGetCode, address, parsedNumOrTag); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *ethNamespace) GetBlockByNumber(number uint64) (*Block[[]string], error) {
	block := new(Block[[]string])
	if err := e.getBlockByNumber(context.Background(), block, number, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByNumberWithContext(ctx context.Context, number uint64) (*Block[[]string], error) {
	block := new(Block[[]string])
	if err := e.getBlockByNumber(ctx, block, number, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByNumberIncTx(number uint64) (*Block[[]*Transaction], error) {
	return e.getBlockByNumberIncTx(context.Background(), number)
}

func (e *ethNamespace) GetBlockByNumberIncTxWithContext(
	ctx context.Context,
	number uint64,
) (*Block[[]*Transaction], error) {
	return e.getBlockByNumberIncTx(ctx, number)
}

func (e *ethNamespace) getBlockByNumberIncTx(ctx context.Context, number uint64) (*Block[[]*Transaction], error) {
	block := new(Block[[]*Transaction])
	if err := e.getBlockByNumber(ctx, block, number, true); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) getBlockByNumber(
	ctx context.Context,
	result interface{},
	number uint64,
	incTx bool,
) error {
	params := []interface{}{hexutil.EncodeUint64(number), incTx}
	if err := e.c.call(ctx, result, ethGetBlockByNumber, params...); err != nil {
		return err
	}
	return nil
}

func (e *ethNamespace) GetBlockByHash(hash string) (*Block[[]string], error) {
	block := new(Block[[]string])
	if err := e.getBlockByHash(context.Background(), block, hash, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByHashWithContext(ctx context.Context, hash string) (*Block[[]string], error) {
	block := new(Block[[]string])
	if err := e.getBlockByHash(ctx, block, hash, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByHashIncTx(hash string) (*Block[[]*Transaction], error) {
	return e.getBlockByHashIncTx(context.Background(), hash)
}

func (e *ethNamespace) GetBlockByHashIncTxWithContext(
	ctx context.Context,
	hash string,
) (*Block[[]*Transaction], error) {
	return e.getBlockByHashIncTx(ctx, hash)
}

func (e *ethNamespace) getBlockByHashIncTx(ctx context.Context, hash string) (*Block[[]*Transaction], error) {
	block := new(Block[[]*Transaction])
	if err := e.getBlockByHash(ctx, block, hash, true); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) getBlockByHash(ctx context.Context, result interface{}, hash string, incTx bool) error {
	params := []interface{}{hash, incTx}
	if err := e.c.call(ctx, result, ethGetBlockByHash, params...); err != nil {
		return err
	}
	return nil
}

func (e *ethNamespace) GetTransaction(hash string) (*Transaction, error) {
	return e.getTransaction(context.Background(), hash)
}

func (e *ethNamespace) GetTransactionWithContext(ctx context.Context, hash string) (*Transaction, error) {
	return e.getTransaction(ctx, hash)
}

func (e *ethNamespace) getTransaction(ctx context.Context, hash string) (*Transaction, error) {
	tx := new(Transaction)
	if err := e.c.call(ctx, tx, ethGetTransaction, hash); err != nil {
		return nil, err
	}
	return tx, nil
}

func (e *ethNamespace) GetReceipt(hash string) (*Receipt, error) {
	return e.getReceipt(context.Background(), hash)
}

func (e *ethNamespace) GetReceiptWithContext(ctx context.Context, hash string) (*Receipt, error) {
	return e.getReceipt(ctx, hash)
}

func (e *ethNamespace) getReceipt(ctx context.Context, hash string) (*Receipt, error) {
	receipt := new(Receipt)
	if err := e.c.call(ctx, receipt, ethGetReceipt, hash); err != nil {
		return nil, err
	}
	return receipt, nil
}
