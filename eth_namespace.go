package evmc

import (
	"context"
	"errors"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

// TODO: get uncle block
// TODO: batch call
// TODO: describe custom functions
// TODO: subscription pending transaction details

type ethNamespace struct {
	info clientInfo
	c    caller
	s    subscriber
}

func (e *ethNamespace) SubscribeNewHeads(
	ctx context.Context,
	ch chan<- *evmctypes.Header,
) (evmctypes.Subscription, error) {
	return e.subscribe(ctx, ch, newHeads)
}

func (e *ethNamespace) SubscribeNewPendingTransactions(
	ctx context.Context,
	ch chan<- string,
) (evmctypes.Subscription, error) {
	return e.subscribe(ctx, ch, newPendingTransactions)
}

func (e *ethNamespace) SubscribeLogs(
	ctx context.Context,
	ch chan<- *evmctypes.Log,
	params *evmctypes.SubLog,
) (evmctypes.Subscription, error) {
	if params == nil {
		return e.s.subscribe(ctx, "eth", ch, logs)
	}
	return e.s.subscribe(ctx, "eth", ch, logs, params)
}

func (e *ethNamespace) subscribe(ctx context.Context, ch interface{}, args ...interface{}) (evmctypes.Subscription, error) {
	if !e.info.IsWebsocket() {
		return nil, ErrWebsocketRequired
	}
	return e.s.subscribe(ctx, "eth", ch, args...)
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

func (e *ethNamespace) GetStorageAt(address, position string, blockAndTag BlockAndTag) (string, error) {
	return e.getStorageAt(context.Background(), address, position, blockAndTag)
}

func (e *ethNamespace) GetStorageAtWithContext(
	ctx context.Context,
	address string,
	position string,
	blockAndTag BlockAndTag,
) (string, error) {
	return e.getStorageAt(ctx, address, position, blockAndTag)
}

func (e *ethNamespace) getStorageAt(
	ctx context.Context,
	address string,
	position string,
	numOrTag BlockAndTag,
) (string, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(numOrTag)
	if err := e.c.call(
		ctx,
		result,
		ethGetStorageAt,
		address,
		position,
		parsedBT,
	); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *ethNamespace) BlockNumber() (uint64, error) {
	return e.blockNumber(context.Background())
}

func (e *ethNamespace) BlockNumberWithContext(ctx context.Context) (uint64, error) {
	return e.blockNumber(ctx)
}

func (e *ethNamespace) blockNumber(ctx context.Context) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethBlockNumber); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) GetCode(address string, blockAndTag BlockAndTag) (string, error) {
	return e.getCode(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetCodeWithContext(ctx context.Context, address string, blockAndTag BlockAndTag) (string, error) {
	return e.getCode(ctx, address, blockAndTag)
}

func (e *ethNamespace) getCode(
	ctx context.Context,
	address string,
	blockAndTag BlockAndTag,
) (string, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.call(ctx, result, ethGetCode, address, parsedBT); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *ethNamespace) GetBlockByTag(tag BlockAndTag) (*evmctypes.Block, error) {
	return e.getBlockByTag(context.Background(), tag)
}

func (e *ethNamespace) GetBlockByTagWithContext(ctx context.Context, tag BlockAndTag) (*evmctypes.Block, error) {
	return e.getBlockByTag(ctx, tag)
}

func (e *ethNamespace) getBlockByTag(ctx context.Context, tag BlockAndTag) (*evmctypes.Block, error) {
	block := new(evmctypes.Block)
	if err := e.getBlockByNumber(ctx, block, tag, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockIncTxByTag(tag BlockAndTag) (*evmctypes.BlockIncTx, error) {
	return e.getBlockIncTxByTag(context.Background(), tag)
}

func (e *ethNamespace) GetBlockByIncTxTagWithContext(ctx context.Context, tag BlockAndTag) (*evmctypes.BlockIncTx, error) {
	return e.getBlockIncTxByTag(ctx, tag)
}

func (e *ethNamespace) getBlockIncTxByTag(ctx context.Context, tag BlockAndTag) (*evmctypes.BlockIncTx, error) {
	block := new(evmctypes.BlockIncTx)
	if err := e.getBlockByNumber(ctx, block, tag, true); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByNumber(number uint64) (*evmctypes.Block, error) {
	return e.getBlock(context.Background(), number)
}

func (e *ethNamespace) GetBlockByNumberWithContext(ctx context.Context, number uint64) (*evmctypes.Block, error) {
	return e.getBlock(ctx, number)
}

func (e *ethNamespace) getBlock(ctx context.Context, number uint64) (*evmctypes.Block, error) {
	result := new(evmctypes.Block)
	if err := e.getBlockByNumber(ctx, result, FormatNumber(number), false); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *ethNamespace) GetBlockIncTxByNumber(number uint64) (*evmctypes.BlockIncTx, error) {
	return e.getBlockIncTxByNumber(context.Background(), number)
}

func (e *ethNamespace) GetBlockIncTxByNumberWithContext(
	ctx context.Context,
	number uint64,
) (*evmctypes.BlockIncTx, error) {
	return e.getBlockIncTxByNumber(ctx, number)
}

func (e *ethNamespace) getBlockIncTxByNumber(ctx context.Context, number uint64) (*evmctypes.BlockIncTx, error) {
	block := new(evmctypes.BlockIncTx)
	if err := e.getBlockByNumber(ctx, block, FormatNumber(number), true); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) getBlockByNumber(
	ctx context.Context,
	result interface{},
	number BlockAndTag,
	incTx bool,
) error {
	if number == Pending {
		return ErrPendingBlockNotSupported
	}
	parsedBT := parseBlockAndTag(number)
	params := []interface{}{parsedBT, incTx}
	if err := e.c.call(ctx, result, ethGetBlockByNumber, params...); err != nil {
		return err
	}
	return nil
}

func (e *ethNamespace) GetBlockByHash(hash string) (*evmctypes.Block, error) {
	block := new(evmctypes.Block)
	if err := e.getBlockByHash(context.Background(), block, hash, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockByHashWithContext(ctx context.Context, hash string) (*evmctypes.Block, error) {
	block := new(evmctypes.Block)
	if err := e.getBlockByHash(ctx, block, hash, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockIncTxByHash(hash string) (*evmctypes.BlockIncTx, error) {
	return e.getBlockIncTxByHash(context.Background(), hash)
}

func (e *ethNamespace) GetBlockIncTxByHashWithContext(
	ctx context.Context,
	hash string,
) (*evmctypes.BlockIncTx, error) {
	return e.getBlockIncTxByHash(ctx, hash)
}

func (e *ethNamespace) getBlockIncTxByHash(ctx context.Context, hash string) (*evmctypes.BlockIncTx, error) {
	block := new(evmctypes.BlockIncTx)
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

func (e *ethNamespace) GetTransaction(hash string) (*evmctypes.Transaction, error) {
	return e.getTransaction(context.Background(), hash)
}

func (e *ethNamespace) GetTransactionWithContext(ctx context.Context, hash string) (*evmctypes.Transaction, error) {
	return e.getTransaction(ctx, hash)
}

func (e *ethNamespace) getTransaction(ctx context.Context, hash string) (*evmctypes.Transaction, error) {
	tx := new(evmctypes.Transaction)
	if err := e.c.call(ctx, tx, ethGetTransaction, hash); err != nil {
		return nil, err
	}
	return tx, nil
}

func (e *ethNamespace) GetTransactionReceipt(hash string) (*evmctypes.Receipt, error) {
	return e.getTransactionReceipt(context.Background(), hash)
}

func (e *ethNamespace) GetTransactionReceiptWithContext(ctx context.Context, hash string) (*evmctypes.Receipt, error) {
	return e.getTransactionReceipt(ctx, hash)
}

func (e *ethNamespace) getTransactionReceipt(ctx context.Context, hash string) (*evmctypes.Receipt, error) {
	receipt := new(evmctypes.Receipt)
	if err := e.c.call(ctx, receipt, ethGetReceipt, hash); err != nil {
		return nil, err
	}
	return receipt, nil
}

func (e *ethNamespace) GetBalance(address string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	return e.getBalance(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetBalanceWithContext(ctx context.Context, address string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	return e.getBalance(ctx, address, blockAndTag)
}

func (e *ethNamespace) getBalance(ctx context.Context, address string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.call(ctx, result, ethGetBalance, address, parsedBT); err != nil {
		return decimal.Zero, err
	}
	if *result == "" {
		*result = "0x0"
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*result), 0), nil
}

func (e *ethNamespace) GetLogs(filter *evmctypes.LogFilter) ([]*evmctypes.Log, error) {
	return e.getLogs(context.Background(), filter)
}

func (e *ethNamespace) GetLogsWithContext(ctx context.Context, filter *evmctypes.LogFilter) ([]*evmctypes.Log, error) {
	return e.getLogs(ctx, filter)
}

func (e *ethNamespace) GetLogsByBlockNumber(number uint64) ([]*evmctypes.Log, error) {
	return e.getLogsByBlockNumber(context.Background(), number)
}

func (e *ethNamespace) GetLogsByBlockNumberWithContext(ctx context.Context, number uint64) ([]*evmctypes.Log, error) {
	return e.getLogsByBlockNumber(ctx, number)
}

func (e *ethNamespace) getLogsByBlockNumber(ctx context.Context, number uint64) ([]*evmctypes.Log, error) {
	filter := &evmctypes.LogFilter{
		FromBlock: &number,
		ToBlock:   &number,
	}
	return e.getLogs(ctx, filter)
}

func (e *ethNamespace) GetLogsByBlockHash(hash string) ([]*evmctypes.Log, error) {
	return e.getLogsByBlockHash(context.Background(), hash)
}

func (e *ethNamespace) GetLogsByBlockHashWithContext(ctx context.Context, hash string) ([]*evmctypes.Log, error) {
	return e.getLogsByBlockHash(ctx, hash)
}

func (e *ethNamespace) getLogsByBlockHash(ctx context.Context, hash string) ([]*evmctypes.Log, error) {
	filter := &evmctypes.LogFilter{
		BlockHash: &hash,
	}
	return e.getLogs(ctx, filter)
}

// TODO: addresses
func (e *ethNamespace) getLogs(ctx context.Context, filter *evmctypes.LogFilter) ([]*evmctypes.Log, error) {
	logs := new([]*evmctypes.Log)
	params := make(map[string]interface{})
	if filter.BlockHash != nil {
		params["blockHash"] = *filter.BlockHash
	} else if filter.FromBlock != nil && filter.ToBlock != nil {
		params["fromBlock"] = hexutil.EncodeUint64(*filter.FromBlock)
		params["toBlock"] = hexutil.EncodeUint64(*filter.ToBlock)
	} else {
		return nil, errors.New("either block hash or block range must be specified")
	}
	if filter.Address != nil {
		params["address"] = *filter.Address
	}
	if filter.Topics != nil {
		params["topics"] = filter.Topics
	}
	if err := e.c.call(ctx, logs, ethGetLogs, params); err != nil {
		return nil, err
	}
	return *logs, nil
}

func (e *ethNamespace) GetTransactionCount(address string, blockAndTag BlockAndTag) (uint64, error) {
	return e.getTransactionCount(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetTransactionCountWithContext(
	ctx context.Context,
	address string,
	blockAndTag BlockAndTag,
) (uint64, error) {
	return e.getTransactionCount(ctx, address, blockAndTag)
}

func (e *ethNamespace) getTransactionCount(ctx context.Context, address string, blockAndTag BlockAndTag) (uint64, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.call(ctx, result, ethGetTransactionCount, address, parsedBT); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) GetBlockReceipts(number uint64) ([]*evmctypes.Receipt, error) {
	return e.getBlockReceipts(context.Background(), number)
}

func (e *ethNamespace) GetBlockReceiptsWithContext(ctx context.Context, number uint64) ([]*evmctypes.Receipt, error) {
	return e.getBlockReceipts(ctx, number)
}

func (e *ethNamespace) getBlockReceipts(ctx context.Context, number uint64) ([]*evmctypes.Receipt, error) {
	var (
		result        = new([]*evmctypes.Receipt)
		method        = ethGetBlockReceipts
		clientName, _ = e.info.NodeClient()
	)
	if ClientName(clientName) == Bor {
		method = ethGetTransactionReceiptsByBlock
	}
	if err := e.c.call(ctx, result, method, hexutil.EncodeUint64(number)); err != nil {
		return nil, err
	}
	return *result, nil
}

func (e *ethNamespace) GasPrice() (decimal.Decimal, error) {
	return e.gasPrice(context.Background())
}

func (e *ethNamespace) GasPriceWithContext(ctx context.Context) (decimal.Decimal, error) {
	return e.gasPrice(ctx)
}

func (e *ethNamespace) gasPrice(ctx context.Context) (decimal.Decimal, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethGasPrice); err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*result), 0), nil
}

func (e *ethNamespace) MaxPriorityFeePerGas() (decimal.Decimal, error) {
	return e.maxPriorityFeePerGas(context.Background())
}

func (e *ethNamespace) MaxPriorityFeePerGasWithContext(ctx context.Context) (decimal.Decimal, error) {
	return e.maxPriorityFeePerGas(ctx)
}

func (e *ethNamespace) maxPriorityFeePerGas(ctx context.Context) (decimal.Decimal, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethMaxPriorityFeePerGas); err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*result), 0), nil
}

func (e *ethNamespace) SendTransaction(sendingTx *SendingTx, wallet *Wallet) (string, error) {
	return e.sendTransaction(context.Background(), sendingTx, wallet)
}

func (e *ethNamespace) SendTransactionWithContext(
	ctx context.Context,
	sendingTx *SendingTx,
	wallet *Wallet,
) (string, error) {
	return e.sendTransaction(ctx, sendingTx, wallet)
}

func (e *ethNamespace) sendTransaction(
	ctx context.Context,
	sendingTx *SendingTx,
	wallet *Wallet,
) (string, error) {
	hash, rawTx, err := wallet.SignTx(sendingTx, e.info.ChainID())
	if err != nil {
		return "", err
	}
	txHash, err := e.sendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", err
	}
	if hash != txHash {
		return "", errors.New("transaction hash mismatch")
	}
	return txHash, nil
}

func (e *ethNamespace) SendRawTransaction(rawTx string) (string, error) {
	return e.sendRawTransaction(context.Background(), rawTx)
}

func (e *ethNamespace) SendRawTransactionWithContext(ctx context.Context, rawTx string) (string, error) {
	return e.sendRawTransaction(ctx, rawTx)
}

func (e *ethNamespace) sendRawTransaction(ctx context.Context, rawTx string) (string, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethSendRawTransaction, rawTx); err != nil {
		return "", err
	}
	return *result, nil
}
