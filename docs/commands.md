# Commands

## Test

```bash
# Run all tests (requires live RPC endpoints)
go test ./...

# Unit tests only (no network required)
go test ./evmctypes/... ./evmcutils/... ./evmcsoltypes/...

# Run a single test
go test -v -run Test_ethNamespace_BlockNumber ./...
```

## Build

```bash
go build ./...
```

## Lint

```bash
# Requires golangci-lint
golangci-lint run
```
