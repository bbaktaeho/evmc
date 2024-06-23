package main

import (
	"github.com/bbaktaeho/evmc"
)

func main() {
	// set url to connect to blockchain node
	client, err := evmc.New("https://ethereum-mainnet.nodit.io")
	if err != nil {
		panic(err)
	}

	name, err := client.ERC20().Name("0xdac17f958d2ee523a2206206994597c13d831ec7", nil)
	if err != nil {
		panic(err)
	}
	println(name)

	symbol, err := client.ERC20().Symbol("0xdac17f958d2ee523a2206206994597c13d831ec7", nil)
	if err != nil {
		panic(err)
	}
	println(symbol)

	totalSupply, err := client.ERC20().TotalSupply("0xdac17f958d2ee523a2206206994597c13d831ec7", nil)
	if err != nil {
		panic(err)
	}
	println(totalSupply.String())
}
