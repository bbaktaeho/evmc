package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bbaktaeho/evmc"
)

func main() {
	// set url to connect to blockchain node
	client, err := evmc.New("http://localhost:8545")
	if err != nil {
		panic(err)
	}

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

	balance, err := client.ERC20().BalanceOf("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0a3f6849f78076aefaDf113F5BED87720274dDC0", evmc.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println(balance.String())

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
