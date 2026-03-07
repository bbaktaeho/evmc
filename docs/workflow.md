# Development Workflow

## Branch Strategy

```
main        - production
feature/*   - new RPC methods or features
fix/*       - bug fixes
```

## Commit Convention

```
type: subject

type: feat, fix, refactor, docs, test, chore
```

## Adding a New RPC Method

Every new or changed RPC method follows this order:

**1. Spec** — read the official JSON-RPC spec for request/response schema.
Reference links are in [docs/guide.md](guide.md#json-rpc-specs).

**2. Procedure constant**

Add to `procedure.go`:
```go
ProcEthGetBlockByNumber Procedure = "eth_getBlockByNumber"
```

**3. Implement**

Add method to `<namespace>_namespace.go`:
```go
func (e *ethNamespace) GetBlockByNumber(ctx context.Context, number string, fullTx bool) (*evmctypes.Block, error) {
    var result evmctypes.Block
    if err := e.client.CallContext(ctx, &result, string(ProcEthGetBlockByNumber), number, fullTx); err != nil {
        return nil, fmt.Errorf("GetBlockByNumber: %w", err)
    }
    return &result, nil
}
```

**4. Define return type** (if needed)

Add to `evmctypes/evmctypes.go`:
```go
type Block struct {
    Number    *big.Int
    Hash      common.Hash
    // ...
}
```

**5. Implement UnmarshalJSON** (if hex fields exist)

Add to `evmctypes/<namespace>_unmarshaling.go`:
```go
func (b *Block) UnmarshalJSON(data []byte) error {
    // decode hex fields here, not inline in the struct
}
```

**6. Mark implemented**

Update `docs/ethereum-jsonrpc-list.md` — change `[ ]` to `[x]`.

## Testing

### Unit test — `<namespace>_namespace_mock_test.go`

Mock RPC server를 띄워 synthetic JSON payload를 검증한다. 네트워크 불필요.
struct field 매핑과 JSON 언마샬링이 올바른지 확인하는 게 목적이다.

### Golden test — `evmctypes/golden_mainnet_test.go`

`testdata/mainnet/*.json`의 실제 mainnet 응답 스냅샷과 언마샬 결과를 비교한다.
새 타입을 추가했으면 스냅샷을 재생성한다:
```bash
go test ./evmctypes/... -run TestGolden -update
```

### E2E test — `<namespace>_namespace_mainnet_test.go`

실제 RPC 엔드포인트를 호출해 응답을 검증한다. 네트워크 필요.

### Run tests

```bash
# Unit + golden (no network)
go test ./evmctypes/... ./evmcutils/... ./evmcsoltypes/...

# All including E2E
go test ./...

# With race detector
go test -race ./...
```

## PR Process

1. `feature/*` 또는 `fix/*` 브랜치 생성
2. 위 RPC method 체크리스트 순서대로 구현
3. 테스트 실행 확인
4. PR 생성 → 코드 리뷰 → squash merge
5. 브랜치 삭제
