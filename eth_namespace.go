package evmc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

// TODO: get uncle block
// TODO: batch call
// TODO: describe custom functions
// TODO: eth_new... & filter
// TODO: eth_syncing

type ethNamespace struct {
	info clientInfo
	c    caller
	s    subscriber
	ts   transactionSender
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

// BlobBaseFee returns the base fee per gas for the next block in the current block.
func (e *ethNamespace) BlobBaseFee() (decimal.Decimal, error) {
	return e.blobBaseFee(context.Background())
}

// BlobBaseFeeWithContext returns the base fee per gas for the next block in the current block.
func (e *ethNamespace) BlobBaseFeeWithContext(ctx context.Context) (decimal.Decimal, error) {
	return e.blobBaseFee(ctx)
}

func (e *ethNamespace) blobBaseFee(ctx context.Context) (decimal.Decimal, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthBlobBaseFee); err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*result), 0), nil
}

func (e *ethNamespace) CreateAccessList(tx *Tx) (*evmctypes.AccessListResp, error) {
	return e.createAccessList(context.Background(), tx)
}

func (e *ethNamespace) CreateAccessListWithContext(ctx context.Context, tx *Tx) (*evmctypes.AccessListResp, error) {
	return e.createAccessList(ctx, tx)
}

func (e *ethNamespace) createAccessList(ctx context.Context, tx *Tx) (*evmctypes.AccessListResp, error) {
	msg, err := tx.parseCallMsg()
	if err != nil {
		return nil, err
	}
	result := new(evmctypes.AccessListResp)
	if err := e.c.call(ctx, result, EthCreateAccessList, msg); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *ethNamespace) FeeHistory(blockCount uint64, lastBlock evmctypes.BlockAndTag, rewardPercentiles []float64) (*evmctypes.FeeHistory, error) {
	return e.feeHistory(context.Background(), blockCount, lastBlock, rewardPercentiles)
}

func (e *ethNamespace) FeeHistoryWithContext(
	ctx context.Context,
	blockCount uint64,
	lastBlock evmctypes.BlockAndTag,
	rewardPercentiles []float64,
) (*evmctypes.FeeHistory, error) {
	return e.feeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}

func (e *ethNamespace) feeHistory(
	ctx context.Context,
	blockCount uint64,
	lastBlock evmctypes.BlockAndTag,
	rewardPercentiles []float64,
) (*evmctypes.FeeHistory, error) {
	result := new(evmctypes.FeeHistory)
	params := []interface{}{hexutil.EncodeUint64(blockCount), lastBlock.String(), rewardPercentiles}
	if err := e.c.call(ctx, result, EthFeeHistory, params...); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *ethNamespace) BlockTransactionCountByHash(hash string) (uint64, error) {
	return e.blockTransactionCountByHash(context.Background(), hash)
}

func (e *ethNamespace) BlockTransactionCountByHashWithContext(ctx context.Context, hash string) (uint64, error) {
	return e.blockTransactionCountByHash(ctx, hash)
}

func (e *ethNamespace) blockTransactionCountByHash(ctx context.Context, hash string) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthGetBlockTransactionCountByHash, hash); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) BlockTransactionCountByNumber(number uint64) (uint64, error) {
	return e.blockTransactionCountByNumber(context.Background(), number)
}

func (e *ethNamespace) BlockTransactionCountByNumberWithContext(ctx context.Context, number uint64) (uint64, error) {
	return e.blockTransactionCountByNumber(ctx, number)
}

func (e *ethNamespace) blockTransactionCountByNumber(ctx context.Context, number uint64) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthGetBlockTransactionCountByNumber, evmctypes.FormatNumber(number)); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) ChainID() (uint64, error) {
	return e.chainID(context.Background())
}

func (e *ethNamespace) ChainIDWithContext(ctx context.Context) (uint64, error) {
	return e.chainID(ctx)
}

func (e *ethNamespace) chainID(ctx context.Context) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthChainID); err != nil {
		return 0, err
	}
	id, err := hexutil.DecodeUint64(*result)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (e *ethNamespace) GetStorageAt(address, position string, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.getStorageAt(context.Background(), address, position, blockAndTag)
}

func (e *ethNamespace) GetStorageAtWithContext(
	ctx context.Context,
	address string,
	position string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.getStorageAt(ctx, address, position, blockAndTag)
}

func (e *ethNamespace) getStorageAt(
	ctx context.Context,
	address string,
	position string,
	numOrTag evmctypes.BlockAndTag,
) (string, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthGetStorageAt, address, position, numOrTag.String()); err != nil {
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
	if err := e.c.call(ctx, result, EthBlockNumber); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) GetCode(address string, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.getCode(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetCodeWithContext(ctx context.Context, address string, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.getCode(ctx, address, blockAndTag)
}

func (e *ethNamespace) getCode(
	ctx context.Context,
	address string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthGetCode, address, blockAndTag.String()); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *ethNamespace) GetBlockByTag(tag evmctypes.BlockAndTag) (*evmctypes.Block, error) {
	return e.getBlockByTag(context.Background(), tag)
}

func (e *ethNamespace) GetBlockByTagWithContext(ctx context.Context, tag evmctypes.BlockAndTag) (*evmctypes.Block, error) {
	return e.getBlockByTag(ctx, tag)
}

func (e *ethNamespace) getBlockByTag(ctx context.Context, tag evmctypes.BlockAndTag) (*evmctypes.Block, error) {
	block := new(evmctypes.Block)
	if err := e.getBlockByNumber(ctx, block, tag, false); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) GetBlockIncTxByTag(tag evmctypes.BlockAndTag) (*evmctypes.BlockIncTx, error) {
	return e.getBlockIncTxByTag(context.Background(), tag)
}

func (e *ethNamespace) GetBlockByIncTxTagWithContext(ctx context.Context, tag evmctypes.BlockAndTag) (*evmctypes.BlockIncTx, error) {
	return e.getBlockIncTxByTag(ctx, tag)
}

func (e *ethNamespace) getBlockIncTxByTag(ctx context.Context, tag evmctypes.BlockAndTag) (*evmctypes.BlockIncTx, error) {
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
	if err := e.getBlockByNumber(ctx, result, evmctypes.FormatNumber(number), false); err != nil {
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
	if err := e.getBlockByNumber(ctx, block, evmctypes.FormatNumber(number), true); err != nil {
		return nil, err
	}
	return block, nil
}

func (e *ethNamespace) getBlockByNumber(
	ctx context.Context,
	result interface{},
	number evmctypes.BlockAndTag,
	incTx bool,
) error {
	if number == evmctypes.Pending {
		return ErrPendingBlockNotSupported
	}
	params := []interface{}{number.String(), incTx}
	if err := e.c.call(ctx, result, EthGetBlockByNumber, params...); err != nil {
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
	if err := e.c.call(ctx, result, EthGetBlockByHash, params...); err != nil {
		return err
	}
	return nil
}

func (e *ethNamespace) GetTransactionByHash(hash string) (*evmctypes.Transaction, error) {
	return e.getTransactionByHash(context.Background(), hash)
}

func (e *ethNamespace) GetTransactionByHashWithContext(ctx context.Context, hash string) (*evmctypes.Transaction, error) {
	return e.getTransactionByHash(ctx, hash)
}

func (e *ethNamespace) getTransactionByHash(ctx context.Context, hash string) (*evmctypes.Transaction, error) {
	tx := new(evmctypes.Transaction)
	if err := e.c.call(ctx, tx, EthGetTransactionByHash, hash); err != nil {
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
	if err := e.c.call(ctx, receipt, EthGetReceipt, hash); err != nil {
		return nil, err
	}
	return receipt, nil
}

func (e *ethNamespace) GetBalance(address string, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	return e.getBalance(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetBalanceWithContext(ctx context.Context, address string, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	return e.getBalance(ctx, address, blockAndTag)
}

func (e *ethNamespace) getBalance(ctx context.Context, address string, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthGetBalance, address, blockAndTag.String()); err != nil {
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
	if err := e.c.call(ctx, logs, EthGetLogs, params); err != nil {
		return nil, err
	}
	return *logs, nil
}

func (e *ethNamespace) GetTransactionCount(address string, blockAndTag evmctypes.BlockAndTag) (uint64, error) {
	return e.getTransactionCount(context.Background(), address, blockAndTag)
}

func (e *ethNamespace) GetTransactionCountWithContext(
	ctx context.Context,
	address string,
	blockAndTag evmctypes.BlockAndTag,
) (uint64, error) {
	return e.getTransactionCount(ctx, address, blockAndTag)
}

func (e *ethNamespace) getTransactionCount(ctx context.Context, address string, blockAndTag evmctypes.BlockAndTag) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthGetTransactionCount, address, blockAndTag.String()); err != nil {
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
		method        = EthGetBlockReceipts
		clientName, _ = e.info.NodeClient()
	)
	if ClientName(clientName) == Bor {
		method = EthGetTransactionReceiptsByBlock
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
	if err := e.c.call(ctx, result, EthGasPrice); err != nil {
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
	if err := e.c.call(ctx, result, EthMaxPriorityFeePerGas); err != nil {
		return decimal.Zero, err
	}
	return decimal.NewFromBigInt(hexutil.MustDecodeBig(*result), 0), nil
}

func (e *ethNamespace) Syncing() (bool, *evmctypes.Syncing, error) {
	return e.syncing(context.Background())
}

func (e *ethNamespace) SyncingWithContext(ctx context.Context) (bool, *evmctypes.Syncing, error) {
	return e.syncing(ctx)
}

func (e *ethNamespace) syncing(ctx context.Context) (bool, *evmctypes.Syncing, error) {
	var (
		result        interface{}
		resultSyncing = new(evmctypes.Syncing)
	)
	if err := e.c.call(ctx, &result, EthSyncing); err != nil {
		return false, nil, err
	}
	if result == nil {
		return false, nil, fmt.Errorf("syncing is not supported")
	}
	if _, ok := result.(bool); ok {
		return false, nil, nil
	}
	b, err := json.Marshal(result)
	if err != nil {
		return false, nil, err
	}
	if err := json.Unmarshal(b, resultSyncing); err != nil {
		return false, nil, err
	}
	return true, resultSyncing, nil
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
	_, rawTx, err := wallet.SignTx(sendingTx, e.info.ChainID())
	if err != nil {
		return "", err
	}
	txHash, err := e.ts.sendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", err
	}
	return txHash, nil
}

func (e *ethNamespace) SendRawTransaction(rawTx string) (string, error) {
	return e.ts.sendRawTransaction(context.Background(), rawTx)
}

func (e *ethNamespace) SendRawTransactionWithContext(ctx context.Context, rawTx string) (string, error) {
	return e.ts.sendRawTransaction(ctx, rawTx)
}

func (e *ethNamespace) sendRawTransaction(ctx context.Context, rawTx string) (string, error) {
	result := new(string)
	if err := e.c.call(ctx, result, EthSendRawTransaction, rawTx); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *ethNamespace) EstimateGas(tx *Tx) (uint64, error) {
	return e.estimateGas(context.Background(), tx)
}

func (e *ethNamespace) EstimateGasWithContext(ctx context.Context, tx *Tx) (uint64, error) {
	return e.estimateGas(ctx, tx)
}

func (e *ethNamespace) estimateGas(ctx context.Context, tx *Tx) (uint64, error) {
	result := new(string)
	msg, err := tx.parseCallMsg()
	if err != nil {
		return 0, err
	}
	// gas estimation requires from field
	if msg["from"] == "" {
		return 0, ErrFromRequired
	}
	if err := e.c.call(ctx, result, EthEstimateGas, msg); err != nil {
		return 0, err
	}
	return hexutil.MustDecodeUint64(*result), nil
}

func (e *ethNamespace) Call(tx *Tx, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.ethCall(context.Background(), tx, blockAndTag)
}

func (e *ethNamespace) CallWithContext(ctx context.Context, tx *Tx, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.ethCall(ctx, tx, blockAndTag)
}

func (e *ethNamespace) ethCall(ctx context.Context, tx *Tx, blockAndTag evmctypes.BlockAndTag) (string, error) {
	result := new(string)
	msg, err := tx.parseCallMsg()
	if err != nil {
		return "", err
	}
	if err := e.c.call(ctx, result, EthCall, msg, blockAndTag.String()); err != nil {
		return "", err
	}
	return *result, nil
}
