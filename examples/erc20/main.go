package main

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/bbaktaeho/evmc"
	"github.com/shopspring/decimal"
)

func viewExample(client *evmc.Evmc) {
	name, err := client.ERC20().Name("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", evmc.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println(name)

	symbol, err := client.ERC20().Symbol("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", evmc.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println(symbol)

	totalSupply, err := client.ERC20().TotalSupply("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", evmc.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println(totalSupply.String())

	decimals, err := client.ERC20().Decimals("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", evmc.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println(decimals.String())

	balance1, err := client.ERC20().BalanceOf("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0a3f6849f78076aefaDf113F5BED87720274dDC0", evmc.Latest)
	if err != nil {
		panic(err)
	}
	balance2, err := client.ERC20().BalanceOf("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0a3f6849f78076aefaDf113F5BED87720274dDC0", evmc.FormatNumber(18000000))
	if err != nil {
		panic(err)
	}
	fmt.Println(balance1.String(), balance2.String())

	allowance, err := client.ERC20().Allowance("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0a3f6849f78076aefaDf113F5BED87720274dDC0", "0x0a3f6849f78076aefaDf113F5BED87720274dDC0", evmc.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println(allowance.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	_, err = client.ERC20().NameWithContext(ctx, "0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", evmc.Latest)
	if err != nil {
		// context deadline exceeded
		fmt.Println(err.Error())
	}
}

func transferExample(client *evmc.Evmc) {
	var (
		pk                   = "<primary key>"
		tokenAddress         = "0x97AC1b933AA2B4aB14b5f7Cbe9270A2279eb3F21"
		recipient            = "0x2887229b2328470DDA115ba0A41c66485e069DD2"
		amount               = decimal.NewFromBigInt(big.NewInt(123), 18)
		nonce                uint64
		maxFeePerGas         decimal.Decimal
		maxPriorityFeePerGas decimal.Decimal
	)
	wallet, err := evmc.NewWallet(pk)
	if err != nil {
		panic(err)
	}

	n, err := client.Eth().GetTransactionCount(wallet.Address(), evmc.Pending)
	if err != nil {
		panic(err)
	}
	nonce = n

	block, err := client.Eth().GetBlockByTag(evmc.Latest)
	if err != nil {
		panic(err)
	}
	maxFeePerGas = block.NextBaseFee()

	mfpg, err := client.Eth().MaxPriorityFeePerGas()
	if err != nil {
		panic(err)
	}
	maxPriorityFeePerGas = mfpg

	tx := &evmc.Tx{
		To:                   tokenAddress,
		Nonce:                nonce,
		GasLimit:             60000,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
	}
	txHash, err := client.ERC20().Transfer(tx, wallet, recipient, amount)
	if err != nil {
		panic(err)
	}
	fmt.Println(txHash)
}

func main() {
	// set url to connect to blockchain node
	client, err := evmc.New("https://ethereum-mainnet.nodit.io/<api-key>")
	if err != nil {
		panic(err)
	}

	// viewExample(client)
	transferExample(client)
}
