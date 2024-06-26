package main

import (
	"fmt"

	"github.com/bbaktaeho/evmc"
)

func main() {
	// set url to connect to blockchain node
	client, err := evmc.New("https://ethereum-mainnet.nodit.io")
	if err != nil {
		panic(err)
	}

	chainID := client.ChainID()
	fmt.Println(chainID)
	nodeName, nodeVersion := client.NodeClient()
	fmt.Println(nodeName, nodeVersion)

	chainID, err = client.Eth().ChainID()
	if err != nil {
		panic(err)
	}
	fmt.Println(chainID)

	cv, err := client.Web3().ClientVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println(cv)

	storage, err := client.Eth().GetStorageAt("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(storage)

	latestBlockNumber, err := client.Eth().GetBlockNumber()
	if err != nil {
		panic(err)
	}
	fmt.Println(latestBlockNumber)

	bytecode, err := client.Eth().GetCode("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(bytecode)

	block, err := client.Eth().GetBlockByNumber(latestBlockNumber)
	if err != nil {
		panic(err)
	}
	fmt.Println(block.Number, block.Hash, len(block.Transactions))

	blockIncTx, err := client.Eth().GetBlockByNumberIncTx(latestBlockNumber)
	if err != nil {
		panic(err)
	}
	for _, tx := range blockIncTx.Transactions {
		fmt.Println(tx.Hash, tx.From, tx.To, tx.Value)
	}

	block, err = client.Eth().GetBlockByHash(block.Hash)
	if err != nil {
		panic(err)
	}
	fmt.Println(block.Number, block.Hash, len(block.Transactions))

	blockIncTx, err = client.Eth().GetBlockByHashIncTx(block.Hash)
	if err != nil {
		panic(err)
	}
	for _, tx := range blockIncTx.Transactions {
		fmt.Println(tx.Hash, tx.From, tx.To, tx.Value)
	}

	logs, err := client.Eth().GetLogsByBlockNumber(20141602)
	if err != nil {
		panic(err)
	}
	for _, log := range logs {
		fmt.Println(log.Address, log.BlockHash, log.Data, log.LogIndex, len(log.Topics))
	}

	logs, err = client.Eth().GetLogsByBlockHash("0x5a5310472efc52cc3a7afa2460be8438e88d32d7d59fb06da6bb887a146a2001")
	if err != nil {
		panic(err)
	}
	for _, log := range logs {
		fmt.Println(log.Address, log.BlockHash, log.Data, log.LogIndex, len(log.Topics))
	}

	receipts, err := client.Eth().GetBlockReceipts(20141602)
	if err != nil {
		panic(err)
	}
	for _, receipt := range receipts {
		fmt.Println(receipt.BlockHash, receipt.TransactionHash, receipt.Status, receipt.CumulativeGasUsed)
	}

	gasPrice, err := client.Eth().GasPrice()
	if err != nil {
		panic(err)
	}
	fmt.Println(gasPrice)
}
