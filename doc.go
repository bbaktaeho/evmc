// Package evmc provides a user-friendly EVM-compatible blockchain RPC client
// built on top of go-ethereum's rpc package.
//
// It exposes standard Ethereum JSON-RPC namespaces (eth, web3, debug) as well as
// chain-specific namespaces (kaia) and standard token interfaces (ERC-20, ERC-721, ERC-1155).
//
// # Quick Start
//
// Create an HTTP client and query the latest block number:
//
//	client, err := evmc.New("https://your-rpc-endpoint")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	blockNumber, err := client.Eth().BlockNumber()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("latest block:", blockNumber)
//
// # Namespaces
//
// Each JSON-RPC namespace is accessed through a dedicated getter on the [Evmc] client:
//
//   - [Evmc.Eth] – standard eth_* methods (blocks, transactions, receipts, logs)
//   - [Evmc.Web3] – web3_* utility methods (client version)
//   - [Evmc.Debug] – debug_* trace methods (traceTransaction, traceBlockByNumber)
//   - [Evmc.Kaia] – kaia_* methods for the Kaia blockchain
//   - [Evmc.Contract] – raw smart contract calls
//   - [Evmc.ERC20] – ERC-20 token standard methods
//
// # Client Options
//
// Use [Options] functions to configure the client:
//
//	client, err := evmc.New(url,
//	    evmc.WithConnPool(20),
//	    evmc.WithReqTimeout(30 * time.Second),
//	    evmc.WithMaxBatchItems(200),
//	)
//
// # Batch Calls
//
// For high-throughput scenarios, use [Evmc.BatchCallWithContext] to send
// multiple RPC requests in parallel across worker goroutines:
//
//	elements := []rpc.BatchElem{ /* ... */ }
//	err := client.BatchCallWithContext(ctx, elements, 5)
package evmc
