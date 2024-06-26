package evmc

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
)

// Tx is a structure that contains transaction information
// to be sent to the blockchain network by EOA.
type Tx struct {
	Nonce    uint64
	To       string
	Data     string
	Value    decimal.Decimal
	GasPrice decimal.Decimal // legacy
	GasLimit uint64

	AccessList []struct {
		Address     string
		StorageKeys []string
	} // EIP-2930
	MaxPriorityFeePerGas decimal.Decimal // EIP-1559
	MaxFeePerGas         decimal.Decimal // EIP-1559
}

type SendingTx struct {
	txData types.TxData
}

func NewSendingTx(tx *Tx) (*SendingTx, error) {
	return NewDynamicFeeTx(tx)
}

func NewLegacyTx(tx *Tx) (*SendingTx, error) {
	data, toAddress, err := parseDataAndToAddress(tx.Data, tx.To)
	if err != nil {
		return nil, err
	}
	if tx.GasPrice.IsZero() {
		return nil, errors.New("gas price is zero")
	}
	legacyTx := &types.LegacyTx{
		Nonce:    tx.Nonce,
		To:       &toAddress,
		Value:    tx.Value.BigInt(),
		Gas:      tx.GasLimit,
		GasPrice: tx.GasPrice.BigInt(),
		Data:     data,
	}
	return &SendingTx{txData: legacyTx}, nil
}

func NewAccessListTx(tx *Tx) (*SendingTx, error) {
	data, toAddress, err := parseDataAndToAddress(tx.Data, tx.To)
	if err != nil {
		return nil, err
	}
	if tx.GasPrice.IsZero() {
		return nil, errors.New("gas price is zero")
	}
	accessList := make([]types.AccessTuple, len(tx.AccessList))
	for i, access := range tx.AccessList {
		storageKeys := make([]common.Hash, len(access.StorageKeys))
		for j, key := range access.StorageKeys {
			storageKeys[j] = common.HexToHash(key)
		}
		accessList[i] = types.AccessTuple{
			Address:     common.HexToAddress(access.Address),
			StorageKeys: storageKeys,
		}
	}
	accessListTx := &types.AccessListTx{
		Nonce:      tx.Nonce,
		To:         &toAddress,
		Value:      tx.Value.BigInt(),
		Gas:        tx.GasLimit,
		GasPrice:   tx.GasPrice.BigInt(),
		Data:       data,
		AccessList: accessList,
	}
	return &SendingTx{txData: accessListTx}, nil
}

// TODO: auto setting maxPriorityFeePerGas and maxFeePerGas
func NewDynamicFeeTx(tx *Tx) (*SendingTx, error) {
	data, toAddress, err := parseDataAndToAddress(tx.Data, tx.To)
	if err != nil {
		return nil, err
	}
	if tx.MaxFeePerGas.IsZero() {
		return nil, errors.New("maxFeePerGas is zero")
	}
	dynamicFeeTx := &types.DynamicFeeTx{
		Nonce:     tx.Nonce,
		To:        &toAddress,
		Value:     tx.Value.BigInt(),
		Gas:       tx.GasLimit,
		GasTipCap: tx.MaxPriorityFeePerGas.BigInt(),
		GasFeeCap: tx.MaxFeePerGas.BigInt(),
		Data:      data,
	}
	return &SendingTx{txData: dynamicFeeTx}, nil
}

func parseDataAndToAddress(data string, to string) ([]byte, common.Address, error) {
	var d []byte
	if data != "" {
		decodedData, err := hexutil.Decode(data)
		if err != nil {
			return nil, common.Address{}, err
		}
		d = decodedData
	}
	if !common.IsHexAddress(to) {
		return nil, common.Address{}, errors.New("invalid to address")
	}
	return d, common.HexToAddress(to), nil
}
