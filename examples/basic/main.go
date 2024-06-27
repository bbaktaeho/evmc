package main

import (
	"fmt"

	"github.com/bbaktaeho/evmc"
	"github.com/bbaktaeho/evmc/evmctypes"
)

func nodeInfoExample(client *evmc.Evmc) {
	chainID, err := client.Eth().ChainID()
	if err != nil {
		panic(err)
	}
	fmt.Println("chain id:", chainID)

	cv, err := client.Web3().ClientVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println("client version:", cv)

	// caching
	chainID = client.ChainID()
	nodeName, nodeVersion := client.NodeClient()
	fmt.Println("chain id:", chainID, "name:", nodeName, "version:", nodeVersion)
}

func blockAndTagExample(client *evmc.Evmc) {
	Latestblock, err := client.Eth().GetBlockByTag(evmc.Latest)
	if err != nil {
		panic(err)
	}
	safeBlock, err := client.Eth().GetBlockByTag(evmc.Safe)
	if err != nil {
		panic(err)
	}
	finalizedBlock, err := client.Eth().GetBlockByTag(evmc.Finalized)
	if err != nil {
		panic(err)
	}
	ealiestBlock, err := client.Eth().GetBlockByTag(evmc.Earliest)
	if err != nil {
		panic(err)
	}
	fmt.Println("block:", Latestblock.Number, safeBlock.Number, finalizedBlock.Number, ealiestBlock.Number)

	pendingBalance, err := client.Eth().GetBalance(evmc.ZeroAddress, evmc.Pending)
	if err != nil {
		panic(err)
	}
	latestBalance, err := client.Eth().GetBalance(evmc.ZeroAddress, evmc.Latest)
	if err != nil {
		panic(err)
	}
	archiveBalance, err := client.Eth().GetBalance(evmc.ZeroAddress, evmc.FormatNumber(18000000))
	if err != nil {
		panic(err)
	}
	fmt.Println("balance:", pendingBalance, latestBalance, archiveBalance)

	pendingStorage, err := client.Eth().GetStorageAt("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0", evmc.Pending)
	if err != nil {
		panic(err)
	}
	latestStorage, err := client.Eth().GetStorageAt("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0", evmc.Latest)
	if err != nil {
		panic(err)
	}
	archibeStorage, err := client.Eth().GetStorageAt("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0", evmc.FormatNumber(18000000))
	if err != nil {
		panic(err)
	}
	fmt.Println("storage:", pendingStorage, latestStorage, archibeStorage)
}

func transactionAndReceiptExample(client *evmc.Evmc) {
	transaction, err := client.Eth().GetTransaction("0x0bf219063db8f75ba381c6b67d7f0f40e0e1d2b40f92725324b11e4ad72a5dab")
	if err != nil {
		panic(err)
	}
	fmt.Println("transaction:", transaction.From, transaction.To, transaction.Value, transaction.Nonce)

	receipt, err := client.Eth().GetTransactionReceipt("0x0bf219063db8f75ba381c6b67d7f0f40e0e1d2b40f92725324b11e4ad72a5dab")
	if err != nil {
		panic(err)
	}
	fmt.Println("receipt:", receipt.BlockNumber, receipt.TransactionHash, receipt.Status, receipt.CumulativeGasUsed)

	fromBlock, toBlock := uint64(18000000), uint64(18000001)
	logs, err := client.Eth().GetLogs(&evmctypes.LogFilter{FromBlock: &fromBlock, ToBlock: &toBlock})
	if err != nil {
		panic(err)
	}
	fmt.Println("logs:", len(logs))
}

func receiptExample(client *evmc.Evmc) {}

func main() {
	// set url to connect to blockchain node
	client, err := evmc.New("http://localhost:8545")
	if err != nil {
		panic(err)
	}

	nodeInfoExample(client)
	blockAndTagExample(client)
	transactionAndReceiptExample(client)
}
