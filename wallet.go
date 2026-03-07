package evmc

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Wallet struct {
	pk      *ecdsa.PrivateKey
	address string
}

func NewWallet(privateKey string) (*Wallet, error) {
	if privateKey[:2] == "0x" {
		privateKey = privateKey[2:]
	}
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		pk:      pk,
		address: crypto.PubkeyToAddress(pk.PublicKey).Hex(),
	}, nil
}

// SignTx signs sendingTx with the wallet's private key for the given chainID
// and returns the transaction hash and the RLP-encoded raw transaction.
func (w *Wallet) SignTx(sendingTx *SendingTx, chainID uint64) (hash, rawTx string, err error) {
	if chainID == 0 {
		return "", "", errors.New("chainID is zero")
	}
	gethTx := types.NewTx(sendingTx.txData)
	signedTx, err := types.SignTx(gethTx, types.LatestSignerForChainID(new(big.Int).SetUint64(chainID)), w.pk)
	if err != nil {
		return "", "", err
	}
	raw, err := signedTx.MarshalBinary()
	if err != nil {
		return "", "", err
	}
	return signedTx.Hash().Hex(), hexutil.Encode(raw), nil
}

func (w *Wallet) Address() string {
	return w.address
}
