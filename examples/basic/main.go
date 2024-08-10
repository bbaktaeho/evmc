package main

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/bbaktaeho/evmc"
	"github.com/bbaktaeho/evmc/evmcsoltypes"
	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/bbaktaeho/evmc/evmcutils"
	"github.com/shopspring/decimal"
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

	isSyncing, syncing, err := client.Eth().Syncing()
	if err != nil {
		panic(err)
	}
	s, err := json.Marshal(syncing)
	if err != nil {
		panic(err)
	}
	fmt.Printf("isSyncing: %v, syncing: %s\n", isSyncing, string(s))
}

func blockAndTagExample(client *evmc.Evmc) {
	LatestBlock, err := client.Eth().GetBlockByTag(evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	safeBlock, err := client.Eth().GetBlockByTag(evmctypes.Safe)
	if err != nil {
		panic(err)
	}
	finalizedBlock, err := client.Eth().GetBlockByTag(evmctypes.Finalized)
	if err != nil {
		panic(err)
	}
	ealiestBlock, err := client.Eth().GetBlockByTag(evmctypes.Earliest)
	if err != nil {
		panic(err)
	}
	fmt.Println("block:", LatestBlock.Number, safeBlock.Number, finalizedBlock.Number, ealiestBlock.Number)

	pendingBalance, err := client.Eth().GetBalance(evmc.ZeroAddress, evmctypes.Pending)
	if err != nil {
		panic(err)
	}
	latestBalance, err := client.Eth().GetBalance(evmc.ZeroAddress, evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	archiveBalance, err := client.Eth().GetBalance(evmc.ZeroAddress, evmctypes.FormatNumber(4000000))
	if err != nil {
		panic(err)
	}
	fmt.Println("balance:", pendingBalance, latestBalance, archiveBalance)

	pendingStorage, err := client.Eth().GetStorageAt("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0", evmctypes.Pending)
	if err != nil {
		panic(err)
	}
	latestStorage, err := client.Eth().GetStorageAt("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0", evmctypes.Latest)
	if err != nil {
		panic(err)
	}
	archibeStorage, err := client.Eth().GetStorageAt("0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2", "0x0", evmctypes.FormatNumber(18000000))
	if err != nil {
		panic(err)
	}
	fmt.Println("storage:", pendingStorage, latestStorage, archibeStorage)
}

func transactionAndReceiptExample(client *evmc.Evmc) {
	transaction, err := client.Eth().GetTransactionByHash("0x0bf219063db8f75ba381c6b67d7f0f40e0e1d2b40f92725324b11e4ad72a5dab")
	if err != nil {
		panic(err)
	}
	fmt.Println("transaction:", transaction.From, transaction.To, transaction.Value, transaction.Nonce)

	receipt, err := client.Eth().GetTransactionReceipt("0x0bf219063db8f75ba381c6b67d7f0f40e0e1d2b40f92725324b11e4ad72a5dab")
	if err != nil {
		panic(err)
	}
	fmt.Println("receipt:", receipt.BlockNumber, receipt.TransactionHash, receipt.Status, receipt.CumulativeGasUsed)

	fromBlock, toBlock := uint64(4000000), uint64(4000001)
	logs, err := client.Eth().GetLogs(&evmctypes.LogFilter{FromBlock: &fromBlock, ToBlock: &toBlock})
	if err != nil {
		panic(err)
	}
	fmt.Printf("from: %d, to: %d, logs: %d\n", fromBlock, toBlock, len(logs))
}

func callExample(client *evmc.Evmc) {
	from := "0xBB402D125aC13e86457D7516069B359365Ee6F7e"

	nonce, err := client.Eth().GetTransactionCount(from, evmctypes.Pending)
	if err != nil {
		panic(err)
	}
	input, _ := evmcutils.GenerateTxInput(
		"transfer(address,uint256)",
		evmcsoltypes.Address("0x0a3f6849f78076aefaDf113F5BED87720274dDC0"),
		evmcsoltypes.Uint256(decimal.NewFromBigInt(big.NewInt(1), 18)),
	)
	tx := &evmc.Tx{
		From:  from,
		To:    "0x97AC1b933AA2B4aB14b5f7Cbe9270A2279eb3F21",
		Nonce: nonce,
		Data:  input,
	}
	gas, err := client.Eth().EstimateGas(tx)
	if err != nil {
		fmt.Println("failed to estimate gas:", err)
	}
	fmt.Println("gas:", gas)
}

func main() {
	// set url to connect to blockchain node
	client, err := evmc.New("https://ethereum-mainnet.nodit.io/<api-key>")
	if err != nil {
		panic(err)
	}

	nodeInfoExample(client)
	blockAndTagExample(client)
	transactionAndReceiptExample(client)
	callExample(client)
}
