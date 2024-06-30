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
	)
	newHeadsSub, err := client.Eth().SubscribeNewHeads(context.Background(), newHeadsCh)
	if err != nil {
		panic(err)
	}
	newPendingTxHashesSub, err := client.Eth().SubscribeNewPendingTransactions(context.Background(), newPendingTxHashesCh)
	if err != nil {
		panic(err)
	}

	timer := time.NewTimer(time.Minute)
	for {
		select {
		case <-timer.C:
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
		}
	}
}
