package main

import (
	"fmt"

	"github.com/bbaktaeho/evmc"
	"github.com/bbaktaeho/evmc/evmctypes"
)

func main() {
	client, err := evmc.New("http://61.111.3.69:18014")
	if err != nil {
		panic(err)
	}

	r, err := client.ERC20().BalanceOf(
		"0xdac17f958d2ee523a2206206994597c13d831ec7",
		"0x74e7fd0b532f88cf8cc50922f7a8f51e3f320fa7",
		evmctypes.Latest,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(r.String())
}
