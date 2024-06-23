package main

import "github.com/bbaktaeho/evmc"

func main() {
	// set url to connect to blockchain node
	client, err := evmc.New("https://ethereum-mainnet.nodit.io")
	if err != nil {
		panic(err)
	}

	chainID := client.ChainID()
	println(chainID)
	nodeName, nodeVersion := client.NodeClient()
	println(nodeName, nodeVersion)

	chainID, err = client.Eth().ChainID()
	if err != nil {
		panic(err)
	}
	println(chainID)

	cv, err := client.Web3().ClientVersion()
	if err != nil {
		panic(err)
	}
	println(cv)
}
