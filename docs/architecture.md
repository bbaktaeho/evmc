# Architecture

## Core Design

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

## Constructors

```go
evmc.New(httpURL, opts...)                 // HTTP/HTTPS
evmc.NewWithContext(ctx, httpURL, opts...) // HTTP/HTTPS with context
evmc.NewWebsocket(ctx, wsURL, opts...)    // WS/WSS
```

Options (`With*` functions): `WithConnPool`, `WithReqTimeout`, `WithMaxBatchItems`, `WithMaxBatchSize`, `WithBatchCallWorkers`, `WithWsReadBufferSize`, `WithWsWriteBufferSize`, `WithWsMessageSizeLimit`

## Internal Interfaces (evmc.go)

Namespaces depend on these interfaces, all implemented by `*Evmc`:

- `caller` – single (`call`) and batch (`batchCall`, `BatchCallWithContext`) RPC calls
- `subscriber` – WebSocket event subscriptions
- `clientInfo` – chain ID and node version
- `transactionSender` – raw transaction sending

`BatchCallWithContext` splits elements into chunks of `maxBatchItems` and fans them out across `workers` goroutines. Default: 100 items/batch, 3 workers.

## Key Files

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
| `erc20.go` | Auto-generated ABI bindings – DO NOT EDIT |
| `evmctypes/evmctypes.go` | Core EVM types (Block, Tx, Receipt, Log, etc.) |
| `evmctypes/*_unmarshaling.go` | Custom JSON unmarshaling for each type |
