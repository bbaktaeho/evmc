package evmc

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
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

	ChainID uint64 `json:"chainId"` // EIP-155

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

func (t *Tx) valid() error {
	if t.ChainID < 0 {
		return ErrChainIDLessThanZero
	}
	if t.To == "" {
		return ErrToRequired
	}
	if t.GasLimit == 0 {
		return ErrTxGasLimitZero
	}
	if t.GasPrice.IsNegative() {
		return ErrTxGasPriceLessThanZero
	}
	if t.MaxFeePerGas.IsNegative() {
		return ErrTxMaxFeePerGasLessThanZero
	}
	if t.MaxPriorityFeePerGas.IsNegative() {
		return ErrTxMaxPriorityFeePerGasLessThanZero
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
	if !tx.GasPrice.IsZero() {
		return NewLegacyTx(tx)
	}
	return NewDynamicFeeTx(tx)
}

func NewLegacyTx(tx *Tx) (*SendingTx, error) {
	if err := tx.valid(); err != nil {
		return nil, err
	}
	data, toAddress, err := parseDataAndToAddress(tx.Data, tx.To)
	if err != nil {
		return nil, err
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
	if err := tx.valid(); err != nil {
		return nil, err
	}
	data, toAddress, err := parseDataAndToAddress(tx.Data, tx.To)
	if err != nil {
		return nil, err
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
		ChainID:    decimal.NewFromUint64(tx.ChainID).BigInt(),
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
	if err := tx.valid(); err != nil {
		return nil, err
	}
	data, toAddress, err := parseDataAndToAddress(tx.Data, tx.To)
	if err != nil {
		return nil, err
	}
	dynamicFeeTx := &types.DynamicFeeTx{
		ChainID:   decimal.NewFromUint64(tx.ChainID).BigInt(),
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

type SetCodeAuthorization struct {
	ChainID uint64
	Address string
	Nonce   uint64
}

type SignedSetCodeAuthorization struct {
	SetCodeAuthorization
	V uint8
	R *big.Int
	S *big.Int
}

var hasherPool = sync.Pool{
	New: func() any { return sha3.NewLegacyKeccak256() },
}

func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}

func SignSetCode(wallet *Wallet, auth SetCodeAuthorization) (*SignedSetCodeAuthorization, error) {
	if !common.IsHexAddress(auth.Address) {
		return nil, errors.New("invalid SetCodeAuthorization address")
	}
	var (
		sha = hasherPool.Get().(crypto.KeccakState)
		h   common.Hash
	)
	defer hasherPool.Put(sha)

	sha.Reset()
	sha.Write([]byte{0x05})
	rlp.Encode(sha, []any{*uint256.NewInt(auth.ChainID), common.HexToAddress(auth.Address), auth.Nonce})
	sha.Read(h[:])

	sig, err := crypto.Sign(h[:], wallet.pk)
	if err != nil {
		return nil, err
	}
	r, s, _ := decodeSignature(sig)
	return &SignedSetCodeAuthorization{
		SetCodeAuthorization: auth,
		V:                    sig[64],
		R:                    r,
		S:                    s,
	}, nil
}

func NewSetCodeTx(tx *Tx, signedAuthList []SignedSetCodeAuthorization) (*SendingTx, error) {
	if err := tx.valid(); err != nil {
		return nil, err
	}
	data, toAddress, err := parseDataAndToAddress(tx.Data, tx.To)
	if err != nil {
		return nil, err
	}
	authorizations := make([]types.SetCodeAuthorization, len(signedAuthList))
	for i, auth := range signedAuthList {
		if !common.IsHexAddress(auth.Address) {
			return nil, errors.New("invalid SetCodeAuthorization address")
		}
		authorizations[i] = types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(auth.ChainID),
			Address: common.HexToAddress(auth.Address),
			Nonce:   auth.Nonce,
			V:       auth.V,
			R:       *uint256.MustFromBig(auth.R),
			S:       *uint256.MustFromBig(auth.S),
		}
	}
	setCodeTx := &types.SetCodeTx{
		ChainID:   uint256.NewInt(tx.ChainID),
		Nonce:     tx.Nonce,
		To:        toAddress,
		Value:     uint256.MustFromBig(tx.Value.BigInt()),
		Gas:       tx.GasLimit,
		GasFeeCap: uint256.MustFromBig(tx.MaxFeePerGas.BigInt()),
		Data:      data,
		GasTipCap: uint256.MustFromBig(tx.MaxPriorityFeePerGas.BigInt()),
		AuthList:  authorizations,
	}
	return &SendingTx{txData: setCodeTx}, nil
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
