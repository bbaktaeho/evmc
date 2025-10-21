package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bbaktaeho/evmc"
	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/shopspring/decimal"
)

func main() {
	// Set URL to connect to blockchain node
	client, err := evmc.New("https://ethereum-mainnet.nodit.io/<api-key>")
	if err != nil {
		panic(err)
	}

	debugTraceBlockByNumberCallTracer(client)
	debugTraceBlockByNumberFlatCallTracer(client)
	debugTraceBlockByNumberPrestateTracer_default(client)
	debugTraceBlockByNumberPrestateTracer_diff(client)
	debugTraceTransactionCallTracer(client)
	debugTraceTransactionFlatCallTracer(client)
	debugTraceTransactionPrestateTracer_default(client)
	debugTraceTransactionPrestateTracer_diff(client)
}

func debugTraceBlockByNumberCallTracer(client *evmc.Evmc) {
	// test block
	blockNumber := uint64(3000000)

	result, err := client.Debug().TraceBlockByNumber_callTracer(blockNumber, 10*time.Second, nil, nil)
	if err != nil {
		fmt.Printf("Error tracing block: %v\n\n", err)
		return
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling result: %v\n\n", err)
		return
	}
	fmt.Println(string(b))
}

func debugTraceBlockByNumberFlatCallTracer(client *evmc.Evmc) {
	// test block
	blockNumber := uint64(3000000)

	result, err := client.Debug().TraceBlockByNumber_flatCallTracer(blockNumber, 10*time.Second, nil, nil)
	if err != nil {
		fmt.Printf("Error tracing block: %v\n\n", err)
		return
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling result: %v\n\n", err)
		return
	}
	fmt.Println(string(b))
}

func debugTraceTransactionCallTracer(client *evmc.Evmc) {
	// test transaction
	txHash := "0xb95ab9484280074f7b8c6a3cf5ffe2bf0c39168433adcdedc1aacd10d994d95a"

	result, err := client.Debug().TraceTransaction_callTracer(txHash, 10*time.Second, nil, nil)
	if err != nil {
		fmt.Printf("Error tracing transaction: %v\n\n", err)
		return
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling result: %v\n\n", err)
		return
	}
	fmt.Println(string(b))
}
func debugTraceTransactionFlatCallTracer(client *evmc.Evmc) {
	// test transaction
	txHash := "0xb95ab9484280074f7b8c6a3cf5ffe2bf0c39168433adcdedc1aacd10d994d95a"

	result, err := client.Debug().TraceTransaction_flatCallTracer(txHash, 10*time.Second, nil, nil)
	if err != nil {
		fmt.Printf("Error tracing transaction: %v\n\n", err)
		return
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling result: %v\n\n", err)
		return
	}
	fmt.Println(string(b))
}

func debugTraceBlockByNumberPrestateTracer_default(client *evmc.Evmc) {
	// test block
	blockNumber := uint64(3000000)

	result, err := client.Debug().TraceBlockByNumber_prestateTracer(blockNumber, 10*time.Second, nil, nil)
	if err != nil {
		fmt.Printf("Error tracing block: %v\n\n", err)
		return
	}

	frames := make([]evmctypes.PrestateFrame, len(result))
	for i, result := range result {
		frame, err := result.ParseFrames()
		if err != nil {
			fmt.Printf("Error parsing prestate frame: %v\n\n", err)
			continue
		}
		frames[i] = frame
	}
	b, err := json.MarshalIndent(frames, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling frames: %v\n\n", err)
		return
	}
	fmt.Println(string(b))
}

func debugTraceBlockByNumberPrestateTracer_diff(client *evmc.Evmc) {
	// test block
	blockNumber := uint64(3000000)

	result, err := client.Debug().TraceBlockByNumber_prestateTracer(blockNumber, 10*time.Second, nil, &evmc.PrestateTracerConfig{
		DiffMode:       true,
		DisableCode:    true,
		DisableStorage: true,
		IncludeEmpty:   false,
	})
	if err != nil {
		fmt.Printf("Error tracing block: %v\n\n", err)
		return
	}
	frames := make([]*evmctypes.PrestateDiffFrame, len(result))
	for i, result := range result {
		frame, err := result.ParseDiffFrames()
		if err != nil {
			fmt.Printf("Error parsing prestate frame: %v\n\n", err)
			continue
		}
		frames[i] = frame
	}
	dedupNativeHolders := make(map[string]*decimal.Decimal)
	for _, frame := range frames {
		for address, account := range frame.Post {
			if account.Balance == nil {
				continue
			}
			dedupNativeHolders[address] = account.Balance
		}
	}
	for address, balance := range dedupNativeHolders {
		fmt.Printf("Address: %s, Balance: %s\n", address, balance.String())
	}
}

func debugTraceTransactionPrestateTracer_default(client *evmc.Evmc) {
	// test transaction
	txHash := "0xb95ab9484280074f7b8c6a3cf5ffe2bf0c39168433adcdedc1aacd10d994d95a"

	result, err := client.Debug().TraceTransaction_prestateTracer(txHash, 10*time.Second, nil, nil)
	if err != nil {
		fmt.Printf("Error tracing transaction: %v\n\n", err)
		return
	}
	prestateFrame, err := result.ParseFrame()
	if err != nil {
		fmt.Printf("Error parsing prestate frame: %v\n\n", err)
		return
	}
	for address, account := range prestateFrame {
		fmt.Printf("  Address: %s\n", address)
		fmt.Printf("  Balance: %s\n", account.Balance.String())
		fmt.Printf("  Nonce: %d\n", account.Nonce)
		fmt.Printf("  Code: %s\n", account.Code)
		fmt.Printf("  Code Hash: %s\n", account.CodeHash)
		fmt.Printf("  Storage: %v\n", account.Storage)
	}
}

func debugTraceTransactionPrestateTracer_diff(client *evmc.Evmc) {
	// test transaction
	txHash := "0xb95ab9484280074f7b8c6a3cf5ffe2bf0c39168433adcdedc1aacd10d994d95a"

	result, err := client.Debug().TraceTransaction_prestateTracer(txHash, 10*time.Second, nil, &evmc.PrestateTracerConfig{
		DiffMode:       true,
		DisableCode:    true,
		DisableStorage: true,
	})
	if err != nil {
		fmt.Printf("Error tracing transaction: %v\n\n", err)
		return
	}
	frame, err := result.ParseDiffFrame()
	if err != nil {
		fmt.Printf("Error parsing prestate frame: %v\n\n", err)
		return
	}
	for address, account := range frame.Post {
		if account.Balance == nil {
			continue
		}
		preAccount := frame.Pre[address]
		if preAccount.Balance == nil {
			continue
		}
		diff := account.Balance.Sub(*preAccount.Balance)
		fmt.Printf("  Address: %s\n", address)
		fmt.Printf("  Diff Balance: %s\n", diff.String())
	}
}
