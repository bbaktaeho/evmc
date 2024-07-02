package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bbaktaeho/evmc"
	"github.com/bbaktaeho/evmc/evmctypes"
)

func main() {
	// set url to connect to blockchain node

	client, err := evmc.NewWebsocket(context.Background(), "wss://ethereum-mainnet.nodit.io/<api-key>")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	var (
		newHeadsCh           = make(chan *evmctypes.Header, 1)
		newPendingTxHashesCh = make(chan string, 1)
		logsCh               = make(chan *evmctypes.Log, 1)
	)
	newHeadsSub, err := client.Eth().SubscribeNewHeads(context.Background(), newHeadsCh)
	if err != nil {
		panic(err)
	}
	newPendingTxHashesSub, err := client.Eth().SubscribeNewPendingTransactions(context.Background(), newPendingTxHashesCh)
	if err != nil {
		panic(err)
	}
	params := &evmctypes.SubLog{
		Address: "0xdac17f958d2ee523a2206206994597c13d831ec7",
		Topics:  []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"},
	}
	logsSub, err := client.Eth().SubscribeLogs(context.Background(), logsCh, params)
	if err != nil {
		panic(err)
	}

	timer := time.NewTimer(15 * time.Second)
	for {
		select {
		case <-timer.C:
			newHeadsSub.Unsubscribe()
			newPendingTxHashesSub.Unsubscribe()
			logsSub.Unsubscribe()
			fmt.Println("unsubscribed")
			return
		case err := <-newHeadsSub.Err():
			if err != nil {
				panic(err)
			}
		case err := <-newPendingTxHashesSub.Err():
			if err != nil {
				panic(err)
			}
		case header := <-newHeadsCh:
			b, err := json.MarshalIndent(header, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		case txHash := <-newPendingTxHashesCh:
			fmt.Println(txHash)
		case log := <-logsCh:
			b, err := json.MarshalIndent(log, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		}
	}
}
