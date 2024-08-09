package main

import (
	"fmt"
	"math/big"

	"github.com/bbaktaeho/evmc"
	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/shopspring/decimal"
)

const (
	primaryKey = "<primary key>"
	to         = "<to address>"
)

func main() {
	client, err := evmc.New("https://ethereum-mainnet.nodit.io")
	if err != nil {
		panic(err)
	}

	block, err := client.Eth().GetBlockByTag(evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println(block.Number, block.Hash, len(block.Transactions))

	wallet, err := evmc.NewWallet(primaryKey)
	if err != nil {
		panic(err)
	}

	nonce, err := client.Eth().GetTransactionCount(wallet.Address(), evmctypes.Pending)
	if err != nil {
		panic(err)
	}
	fmt.Println(nonce)

	gasPrice, err := client.Eth().GasPrice()
	if err != nil {
		panic(err)
	}

	tx1 := &evmc.Tx{
		To:       to,
		Nonce:    nonce,
		GasLimit: 21000,
		GasPrice: gasPrice,
		Value:    decimal.NewFromBigInt(big.NewInt(1), 17), // 0.1 ETH
	}
	sendingTx, err := evmc.NewLegacyTx(tx1)
	if err != nil {
		panic(err)
	}

	hash, err := client.Eth().SendTransaction(sendingTx, wallet)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)

	maxGas, err := client.Eth().MaxPriorityFeePerGas()
	if err != nil {
		panic(err)
	}

	nonce, err = client.Eth().GetTransactionCount(wallet.Address(), evmctypes.Pending)
	if err != nil {
		panic(err)
	}
	fmt.Println(nonce)

	tx2 := &evmc.Tx{
		To:                   to,
		Nonce:                nonce,
		GasLimit:             21000,
		Value:                decimal.NewFromBigInt(big.NewInt(1), 17), // 0.1 ETH
		MaxFeePerGas:         block.NextBaseFee(),
		MaxPriorityFeePerGas: maxGas,
	}
	sendingTx, err = evmc.NewSendingTx(tx2)
	if err != nil {
		panic(err)
	}

	hash, err = client.Eth().SendTransaction(sendingTx, wallet)
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)
}
