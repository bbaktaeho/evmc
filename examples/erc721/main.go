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

	// Bored Ape Yacht Club
	contractAddr := "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D"
	tokenID := decimal.NewFromInt(1)

	// Name
	name, err := client.ERC721().Name(contractAddr, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name:", name)

	// Symbol
	symbol, err := client.ERC721().Symbol(contractAddr, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("Symbol:", symbol)

	// BalanceOf
	balance, err := client.ERC721().BalanceOf(contractAddr, "0x46EFbAedc92067E6d60E84ED6395099723252496", evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("BalanceOf:", balance.String())

	// OwnerOf
	owner, err := client.ERC721().OwnerOf(contractAddr, tokenID, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("OwnerOf(tokenID=1):", owner)

	// TokenURI
	tokenURI, err := client.ERC721().TokenURI(contractAddr, tokenID, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("TokenURI(tokenID=1):", tokenURI)

	// GetApproved
	approved, err := client.ERC721().GetApproved(contractAddr, tokenID, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("GetApproved(tokenID=1):", approved)

	// IsApprovedForAll
	isApproved, err := client.ERC721().IsApprovedForAll(
		contractAddr,
		"0x46EFbAedc92067E6d60E84ED6395099723252496",
		"0x0000000000000000000000000000000000000001",
		evmctypes.Latest,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("IsApprovedForAll:", isApproved)

	// TotalSupply (ERC721Enumerable)
	totalSupply, err := client.ERC721().TotalSupply(contractAddr, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("TotalSupply:", totalSupply.String())

	// TokenByIndex (ERC721Enumerable)
	tid, err := client.ERC721().TokenByIndex(contractAddr, decimal.NewFromInt(0), evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	fmt.Println("TokenByIndex(0):", tid.String())

	// TokenOfOwnerByIndex (ERC721Enumerable)
	tid, err = client.ERC721().TokenOfOwnerByIndex(
		contractAddr,
		"0x46EFbAedc92067E6d60E84ED6395099723252496",
		decimal.NewFromInt(0),
		evmctypes.Latest,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("TokenOfOwnerByIndex(0):", tid.String())
}
