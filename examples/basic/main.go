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

	storage, err := client.Eth().GetStorageAt("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0", nil)
	if err != nil {
		panic(err)
	}
	println(storage)

	latestBlockNumber, err := client.Eth().GetBlockNumber()
	if err != nil {
		panic(err)
	}
	println(latestBlockNumber)

	bytecode, err := client.Eth().GetCode("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", nil)
	if err != nil {
		panic(err)
	}
	println(bytecode)

	block, err := client.Eth().GetBlockByNumber(latestBlockNumber)
	if err != nil {
		panic(err)
	}
	println(block.Number, block.Hash, len(block.Transactions))

	blockIncTx, err := client.Eth().GetBlockByNumberIncTx(latestBlockNumber)
	if err != nil {
		panic(err)
	}
	for _, tx := range blockIncTx.Transactions {
		println(tx.Hash, tx.From, tx.To, tx.Value)
	}

	block, err = client.Eth().GetBlockByHash(block.Hash)
	if err != nil {
		panic(err)
	}
	println(block.Number, block.Hash, len(block.Transactions))

	blockIncTx, err = client.Eth().GetBlockByHashIncTx(block.Hash)
	if err != nil {
		panic(err)
	}
	for _, tx := range blockIncTx.Transactions {
		println(tx.Hash, tx.From, tx.To, tx.Value)
	}
}
