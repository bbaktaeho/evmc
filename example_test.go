package evmc_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bbaktaeho/evmc"
	"github.com/bbaktaeho/evmc/evmctypes"
)

func ExampleNew() {
	client, err := evmc.New("https://your-rpc-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	blockNumber, err := client.Eth().BlockNumber()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("latest block:", blockNumber)
}

func ExampleNew_withOptions() {
	client, err := evmc.New("https://your-rpc-endpoint",
		evmc.WithConnPool(20),
		evmc.WithReqTimeout(30*time.Second),
		evmc.WithMaxBatchItems(200),
		evmc.WithBatchCallWorkers(5),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	_ = client
}

func ExampleNewWebsocket() {
	ctx := context.Background()
	client, err := evmc.NewWebsocket(ctx, "wss://your-rpc-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	_ = client
}

func ExampleEvmc_Eth() {
	client, err := evmc.New("https://your-rpc-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	block, err := client.Eth().GetBlockByTag(evmctypes.Latest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("block number:", block.Number)
}

func ExampleEvmc_ChainID() {
	client, err := evmc.New("https://your-rpc-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	chainID, err := client.ChainID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("chain ID:", chainID)
}

func ExampleEvmc_NodeClient() {
	client, err := evmc.New("https://your-rpc-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	name, version, err := client.NodeClient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("node: %s, version: %s\n", name, version)
}

func ExampleEvmc_Debug() {
	client, err := evmc.New("https://your-rpc-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	traces, err := client.Debug().TraceBlockByNumber(18000000, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("trace count:", len(traces))
}

func ExampleEvmc_ERC20() {
	client, err := evmc.New("https://your-rpc-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	name, err := client.ERC20().Name("0xdAC17F958D2ee523a2206206994597C13D831ec7", evmctypes.Latest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("token name:", name)
}
