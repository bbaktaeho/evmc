# Debug PrestateTracer Example

This example demonstrates how to use the `prestateTracer` to trace Ethereum transactions and inspect account states.

## What is PrestateTracer?

The prestate tracer is a native Go tracer that shows the state of accounts before and/or after a transaction execution. It has two modes:

1. **Standard Mode** (`diffMode: false`): Shows the state of all accounts accessed during transaction execution BEFORE the transaction was executed.

2. **Diff Mode** (`diffMode: true`): Shows both the state BEFORE and AFTER the transaction execution, making it easy to see what changed.

## Features Demonstrated

This example shows three different ways to use the prestate tracer:

### 1. Standard Prestate Mode

- Traces a transaction with `diffMode: false`
- Shows account states before transaction execution
- Displays balance, nonce, code, and storage for each account

### 2. Diff Mode

- Traces a transaction with `diffMode: true`
- Shows both pre-state and post-state
- Makes it easy to see state changes

### 3. Auto-parse Result

- Automatically detects which mode was used
- Compares pre and post states
- Calculates balance changes

## Configuration Options

The `PrestateTracerConfig` supports the following options:

```go
type PrestateTracerConfig struct {
    // If true, this tracer will return state modifications
    DiffMode bool `json:"diffMode"`
    // If true, this tracer will not return the contract code
    DisableCode bool `json:"disableCode"`
    // If true, this tracer will not return the contract storage
    DisableStorage bool `json:"disableStorage"`
    // If true, this tracer will return empty state objects
    IncludeEmpty bool `json:"includeEmpty"`
}
```

## Usage

Replace `<api-key>` in the code with your actual API key:

```go
client, err := evmc.New("https://ethereum-mainnet.nodit.io/<api-key>")
```

Replace the transaction hash with a real transaction you want to trace:

```go
txHash := "0x1234..."
```

Run the example:

```bash
go run main.go
```

## PrestateTracer API

### Methods

- `ParseResult()`: Automatically parses result as either `PrestateFrame` or `PrestateDiffFrame`
- `ParsePrestateFrame()`: Explicitly parses result as `PrestateFrame` (standard mode)
- `ParseDiffFrame()`: Explicitly parses result as `PrestateDiffFrame` (diff mode)
- `IsDiffMode()`: Checks if the result is in diff mode

### Example Usage

```go
// Trace with standard mode
result, err := client.Debug().TraceTransaction_prestateTracer(
    txHash,
    10*time.Second,
    nil,
    &evmc.PrestateTracerConfig{DiffMode: false},
)

// Parse the result
prestateFrame, err := result.ParsePrestateFrame()
for address, account := range *prestateFrame {
    fmt.Printf("Address: %s, Balance: %s\n", address, account.Balance)
}
```

```go
// Trace with diff mode
result, err := client.Debug().TraceTransaction_prestateTracer(
    txHash,
    10*time.Second,
    nil,
    &evmc.PrestateTracerConfig{DiffMode: true},
)

// Parse the result
diffFrame, err := result.ParseDiffFrame()
for address, postAccount := range diffFrame.Post {
    preAccount := diffFrame.Pre[address]
    if preAccount.Balance != nil && postAccount.Balance != nil {
        diff := postAccount.Balance.Sub(*preAccount.Balance)
        fmt.Printf("Address: %s, Balance change: %s\n", address, diff)
    }
}
```

## Use Cases

- **State inspection**: View the state of accounts before transaction execution
- **State change analysis**: Compare account states before and after execution
- **Balance tracking**: Track balance changes for specific addresses
- **Storage analysis**: Inspect storage changes in smart contracts
- **Debugging**: Debug failed transactions by examining state

## Notes

- The prestate tracer only includes accounts that were accessed during transaction execution
- In diff mode, both pre and post states are included
- Large transactions may take longer to trace
- Archive nodes are required to trace old transactions
