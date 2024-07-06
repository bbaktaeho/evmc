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
	From     string          `json:"from"`
	Nonce    uint64          `json:"nonce"`
	To       string          `json:"to"`
	Data     string          `json:"data"`
	Value    decimal.Decimal `json:"value"`
	GasPrice decimal.Decimal `json:"gasPrice"` // legacy
	GasLimit uint64          `json:"gasLimit"`

	AccessList []struct {
		Address     string   `json:"address"`
		StorageKeys []string `json:"storageKeys"`
	} `json:"accessList"` // EIP-2930
	MaxPriorityFeePerGas decimal.Decimal `json:"maxPriorityFeePerGas"` // EIP-1559
	MaxFeePerGas         decimal.Decimal `json:"maxFeePerGas"`         // EIP-1559
}

func (t *Tx) parseCallMsg() (map[string]interface{}, error) {
	accessList := make([]map[string]interface{}, len(t.AccessList))
	for i, access := range t.AccessList {
		storageKeys := make([]string, len(access.StorageKeys))
		for j, key := range access.StorageKeys {
			storageKeys[j] = key
		}
		accessList[i] = map[string]interface{}{
			"address":     access.Address,
			"storageKeys": storageKeys,
		}
	}
	if t.From == "" {
		// return nil, ErrFromRequired
	}
	if t.To == "" {
		return nil, ErrToRequired
	}
	return map[string]interface{}{
		"from":       t.From,
		"nonce":      hexutil.EncodeUint64(t.Nonce),
		"to":         t.To,
		"data":       t.Data,
		"value":      hexutil.EncodeBig(t.Value.BigInt()),
		"gasLimit":   hexutil.EncodeUint64(t.GasLimit),
		"accessList": accessList,
	}, nil
}

func (t *Tx) checkSendingTx() error {
	if t.To == "" {
		return ErrToRequired
	}
	if t.GasLimit == 0 {
		return ErrTxGasLimitZero
	}
	if t.GasPrice.IsNegative() {
		return ErrTxGasPriceZero
	}
	if t.MaxFeePerGas.IsNegative() {
		return ErrTxMaxFeePerGasZero
	}
	if t.MaxPriorityFeePerGas.IsNegative() {
		return ErrTxMaxPriorityFeePerGasZero
	}
	if t.Value.IsNegative() {
		return ErrTxValueLessThanZero
	}
	return nil
}

type SendingTx struct {
	txData types.TxData
}

func NewSendingTx(tx *Tx) (*SendingTx, error) {
	if err := tx.checkSendingTx(); err != nil {
		return nil, err
	}
	if !tx.GasPrice.IsZero() {
		return NewLegacyTx(tx)
	}
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
