package main

import (
	"context"

	"github.com/bbaktaeho/evmc"
)

func main() {
	// set url to connect to blockchain node
	client, err := evmc.New("https://ethereum-mainnet.nodit.io")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	name, err := client.ERC20().Name(ctx, "0xdac17f958d2ee523a2206206994597c13d831ec7", nil)
	if err != nil {
		panic(err)
	}
	println(name)
}
