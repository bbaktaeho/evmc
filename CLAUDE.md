# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

evmc (Ethereum Virtual Machine Client) is a Go library that wraps go-ethereum's `rpc` package to provide a user-friendly EVM-compatible blockchain RPC client. It supports standard Ethereum JSON-RPC namespaces plus debug/trace for data analysis and standard token interfaces.

**Status:** Work In Progress (WIP)

## Commands

```bash
# 전체 테스트 실행 (live RPC 엔드포인트 필요)
go test ./...

# 네트워크 불필요한 유닛 테스트만 실행
go test ./evmctypes/... ./evmcutils/... ./evmcsoltypes/...

# 단일 테스트 실행
go test -v -run Test_ethNamespace_BlockNumber ./...

# 빌드 확인
go build ./...

# 린트 (golangci-lint 설치된 경우)
golangci-lint run
```

## Architecture

### Core Design

The `Evmc` struct (`evmc.go`) is the single entry point. It holds the underlying `rpc.Client` and exposes namespaces via getter methods:

```
client.Eth()      -> *ethNamespace      (eth_* methods)
client.Web3()     -> *web3Namespace     (web3_* methods)
client.Debug()    -> *debugNamespace    (debug_* trace methods)
client.Kaia()     -> *kaiaNamespace     (kaia_* methods)
client.Contract() -> *contract          (raw contract calls)
client.ERC20()    -> *erc20Contract     (ERC-20 token methods)
client.ERC721()   -> *erc721Contract    (ERC-721 token methods)
client.ERC1155()  -> *erc1155Contract   (ERC-1155 token methods)
```

### Constructors

```go
evmc.New(httpURL, opts...)                    // HTTP/HTTPS
evmc.NewWithContext(ctx, httpURL, opts...)    // HTTP/HTTPS with context
evmc.NewWebsocket(ctx, wsURL, opts...)        // WS/WSS
```

Options (`With*` functions): `WithConnPool`, `WithReqTimeout`, `WithMaxBatchItems`, `WithMaxBatchSize`, `WithBatchCallWorkers`, `WithWsReadBufferSize`, `WithWsWriteBufferSize`, `WithWsMessageSizeLimit`

### Internal Interfaces (evmc.go)

Namespaces depend on these interfaces, all implemented by `*Evmc`:
- `caller` - single (`call`) and batch (`batchCall`, `BatchCallWithContext`) RPC calls
- `subscriber` - WebSocket event subscriptions
- `clientInfo` - chain ID and node version
- `transactionSender` - raw transaction sending

`BatchCallWithContext` splits elements into chunks of `maxBatchItems` and fans them out across `workers` goroutines. Default: 100 items/batch, 3 workers.

### Key Files

| File | Role |
|------|------|
| `evmc.go` | `Evmc` struct, constructors, `BatchCallWithContext`, internal interfaces |
| `procedure.go` | All RPC method name constants (`Procedure` type) |
| `options.go` | Client configuration (connection pool, timeouts, batch limits) |
| `eth_namespace.go` | All `eth_*` RPC method implementations |
| `debug_namespace.go` | All `debug_*` trace RPC implementations |
| `kaia_namespace.go` | All `kaia_*` RPC method implementations |
| `sending_transaction.go` | Transaction building and signing helpers |
| `wallet.go` | Key management, address derivation |
| `erc20.go` | Auto-generated ABI bindings - DO NOT EDIT |
| `evmctypes/evmctypes.go` | Core EVM types (Block, Tx, Receipt, Log, etc.) |
| `evmctypes/*_unmarshaling.go` | Custom JSON unmarshaling for each type |

## Adding a New RPC Method

Follow this exact sequence:

1. **`procedure.go`** - Add a `Procedure` constant for the RPC method name
   ```go
   EthGetProof Procedure = "eth_getProof"
   ```

2. **`eth_namespace.go`** (or relevant namespace file) - Implement the method
   ```go
   func (e *ethNamespace) GetProof(address string, keys []string, block interface{}) (*evmctypes.Proof, error) {
       result := new(evmctypes.Proof)
       if err := e.c.call(context.Background(), result, EthGetProof, address, keys, block); err != nil {
           return nil, err
       }
       return result, nil
   }
   ```

3. **`evmctypes/evmctypes.go`** - Add the return type struct if needed

4. **`evmctypes/*_unmarshaling.go`** - Add custom JSON unmarshaling if the RPC response uses hex-encoded fields

5. **`*_namespace_test.go`** - Add an integration test (see Testing section below)

6. **`docs/ethereum-jsonrpc-list.md`** - Mark the method as implemented

## Adding a New Namespace

1. Create `<name>_namespace.go` with a struct that embeds the required interfaces
2. Add the namespace field to `Evmc` struct in `evmc.go`
3. Initialize it in `newClient()` in `evmc.go`
4. Add a getter method on `Evmc`
5. Add Procedure constants to `procedure.go`

## Testing

### Important: Tests Require Live RPC Nodes

All tests are **integration tests that make real network calls**. There are no mocks.

- Tests using hardcoded public endpoints (BSC, etc.) may pass in CI
- Tests using private/local endpoints (e.g., `http://61.111.3.69:18014`) will fail in most environments
- Tests with empty URL string (`testEvmc("")`) will always panic - these are placeholders

### Test Pattern

Use table-driven tests with testify:

```go
func Test_ethNamespace_MethodName(t *testing.T) {
    t.Parallel()
    tests := []struct {
        name string
        url  string
        // inputs and expected values
    }{
        {name: "mainnet", url: "https://..."},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := testEvmc(tt.url)
            // call and assert
        })
    }
}
```

## Code Conventions

- Receiver name: single letter matching struct name (`e` for ethNamespace, `d` for debugNamespace)
- Error wrapping: `fmt.Errorf("methodName: %w", err)`
- Use `context.Background()` for non-context-aware callers; pass context through for public APIs
- RPC responses use hex strings for numbers - define custom `UnmarshalJSON` in a separate `*_unmarshaling.go` file. See `evmctypes/block_unmarshaling.go` for the pattern.
- Block parameters accept `"latest"`, `"earliest"`, `"pending"`, or a hex number string. Use `evmcutils` helpers for conversion.

## Known Limitations and WIP Items

- WebSocket RPC: partially implemented, no reconnect/backoff
- ERC-721 and ERC-1155 contract methods: not fully implemented
- Trace namespace (`trace_block`, `arbtrace_block`): stub only in `trace_namespace.go`
- Ots namespace (Otterscan): not implemented
- `erc20.go` is auto-generated by abigen - regenerate if ABI changes, never edit manually

## Do Not

- Edit `erc20.go` directly (auto-generated)
- Add observability/OpenTelemetry instrumentation (not in scope)
- Add global variables or package-level state
- Skip the Procedure constant step when adding RPC methods
