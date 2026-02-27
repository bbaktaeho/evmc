# Development Guide

새로운 기능을 개발하기 전에 반드시 RPC Spec을 먼저 확인하고, 아래 워크플로우 순서대로 진행합니다.

---

## 1. RPC Spec 확인

구현하려는 메서드의 요청/응답 스키마를 공식 스펙에서 먼저 확인합니다.

### Ethereum (eth, debug)

| Resource | URL |
|----------|-----|
| Execution-APIs Repository | https://github.com/ethereum/execution-apis/tree/main/src |
| eth namespace spec | https://github.com/ethereum/execution-apis/tree/main/src/eth |
| debug namespace spec | https://github.com/ethereum/execution-apis/tree/main/src/debug |
| JSON schemas | https://github.com/ethereum/execution-apis/tree/main/src/schemas |

### Kaia

| Resource | URL |
|----------|-----|
| Kaia RPC Specs | https://github.com/kaiachain/kaia-sdk/tree/dev/web3rpc/rpc-specs |

### Implementation Status

현재 어떤 메서드가 구현되어 있는지는 [ethereum-jsonrpc-list.md](ethereum-jsonrpc-list.md)에서 확인합니다.

---

## 2. Development Workflow

### 2-1. Adding a New RPC Method

| 순서 | 파일 | 작업 |
|:----:|------|------|
| 1 | `procedure.go` | `Procedure` 상수 추가 |
| 2 | `<namespace>_namespace.go` | 메서드 구현 |
| 3 | `evmctypes/evmctypes.go` | 반환 타입 struct 정의 (필요 시) |
| 4 | `evmctypes/*_unmarshaling.go` | hex 필드 `UnmarshalJSON` 구현 (필요 시) |
| 5 | `<namespace>_namespace_test.go` | 통합 테스트 작성 ([Testing](testing.md) 참고) |
| 6 | `docs/ethereum-jsonrpc-list.md` | 구현 상태 체크 표시 |

**Step 1 – Procedure 상수**

```go
EthGetProof Procedure = "eth_getProof"
```

**Step 2 – 메서드 구현**

```go
func (e *ethNamespace) GetProof(address string, keys []string, block interface{}) (*evmctypes.Proof, error) {
    result := new(evmctypes.Proof)
    if err := e.c.call(context.Background(), result, EthGetProof, address, keys, block); err != nil {
        return nil, err
    }
    return result, nil
}
```

### 2-2. Adding a New Namespace

| 순서 | 파일 | 작업 |
|:----:|------|------|
| 1 | `<name>_namespace.go` | 네임스페이스 struct 생성 (필요한 인터페이스 임베드) |
| 2 | `evmc.go` | `Evmc` struct에 필드 추가 |
| 3 | `evmc.go` | `newClient()`에서 초기화 |
| 4 | `evmc.go` | getter 메서드 추가 |
| 5 | `procedure.go` | `Procedure` 상수 추가 |
| 6 | `<name>_namespace_test.go` | 통합 테스트 작성 |

### 2-3. Adding a New Contract Interface (ERC-*)

| 순서 | 파일 | 작업 |
|:----:|------|------|
| 1 | `abi/*.json` | ABI JSON 추가 |
| 2 | `<name>.go` | abigen으로 바인딩 생성 또는 수동 구현 |
| 3 | `evmc.go` | `Evmc` struct에 필드 추가 + getter 메서드 |
| 4 | `<name>_test.go` | 통합 테스트 작성 |

### 2-4. Adding a New Type

| 순서 | 파일 | 작업 |
|:----:|------|------|
| 1 | `evmctypes/evmctypes.go` | struct 정의 |
| 2 | `evmctypes/<name>_unmarshaling.go` | `UnmarshalJSON` 구현 (hex 필드가 있을 때) |
| 3 | `evmctypes/<name>_test.go` | 언마샬링 테스트 |

---

## 3. Code Conventions

- **Receiver name**: struct 이름의 첫 글자 소문자 (`e` for ethNamespace, `d` for debugNamespace)
- **Error wrapping**: `fmt.Errorf("methodName: %w", err)`
- **Context**: public API는 `context.Context`를 전달하고, 그 외에는 `context.Background()` 사용
- **Hex decoding**: RPC 응답의 hex 문자열은 별도 `*_unmarshaling.go` 파일에 `UnmarshalJSON`으로 디코딩. `evmctypes/block_unmarshaling.go` 참고
- **Block parameter**: `"latest"`, `"earliest"`, `"pending"` 또는 hex number string. `evmcutils` 헬퍼 사용
