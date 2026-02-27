# Testing

테스트 작성 전에 [Development Guide](development-guide.md)의 워크플로우를 먼저 확인하세요.

## Important: Tests Require Live RPC Nodes

All tests are **integration tests that make real network calls**. There are no mocks.

- Tests using hardcoded public endpoints (BSC, etc.) may pass in CI
- Tests using private/local endpoints (e.g., `http://61.111.3.69:18014`) will fail in most environments
- Tests with empty URL string (`testEvmc("")`) will always panic – these are placeholders

## Test Pattern

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

## Commands

```bash
# Run all tests (requires live RPC endpoints)
go test ./...

# Unit tests only (no network required)
go test ./evmctypes/... ./evmcutils/... ./evmcsoltypes/...

# Run a single test
go test -v -run Test_ethNamespace_BlockNumber ./...
```
