package main

import (
	"fmt"

	"github.com/bbaktaeho/evmc"
	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/shopspring/decimal"
)

func main() {
	client, err := evmc.New("https://ethereum-mainnet.nodit.io/<api-key>")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Rarible (ERC-1155)
	contractAddr := "0xd07dc4262BCDbf85190C01c996b4C06a461d2430"
	holder := "0xc5e08104c19DAfd00Fe40737490Da9552Db5bFe5"
	tokenID := decimal.NewFromInt(53776)

	// BalanceOf
	balance, err := client.ERC1155().BalanceOf(contractAddr, holder, tokenID, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("BalanceOf:", balance.String())

	// BalanceOfBatch
	owners := []string{holder, holder}
	ids := []decimal.Decimal{tokenID, decimal.NewFromInt(53777)}
	balances, err := client.ERC1155().BalanceOfBatch(contractAddr, owners, ids, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	for i, b := range balances {
		fmt.Printf("BalanceOfBatch[%d] (id=%s): %s\n", i, ids[i].String(), b.String())
	}

	// IsApprovedForAll
	isApproved, err := client.ERC1155().IsApprovedForAll(
		contractAddr,
		holder,
		"0x0000000000000000000000000000000000000001",
		evmctypes.Latest,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("IsApprovedForAll:", isApproved)

	// URI
	uri, err := client.ERC1155().URI(contractAddr, tokenID, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("URI:", uri)
}
