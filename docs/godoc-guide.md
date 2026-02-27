# GoDoc & pkg.go.dev Guide

이 프로젝트의 코드 주석은 [pkg.go.dev](https://pkg.go.dev) 문서 페이지에 자동으로 렌더링됩니다.
문서 품질을 유지하려면 아래 규칙을 따릅니다.

## 핵심 개념

Go는 코드 주석을 공식 API 문서로 사용합니다. 별도의 문서 생성 도구 없이, 주석만으로 세 가지가 만들어집니다:

| 요소 | 위치 | pkg.go.dev 렌더링 |
|------|------|-------------------|
| **Package comment** | `doc.go`의 `package` 선언 위 | 패키지 Overview 페이지 |
| **Symbol comment** | exported 함수/타입/상수 바로 위 | API Reference 섹션 |
| **Example function** | `_test.go`의 `func Example*()` | 해당 심볼 아래 실행 가능 예제 |

## 1. Package Comment (doc.go)

각 패키지에 `doc.go` 파일을 두고, `package` 선언 바로 위에 패키지 설명을 작성합니다.

```go
// Package evmctypes defines the core data types returned by EVM-compatible
// blockchain RPC methods.
package evmctypes
```

규칙:
- 첫 문장은 `Package <name>` 으로 시작
- 완전한 영문 문장 사용
- 코드 예시를 넣으려면 한 탭 들여쓰기

## 2. Symbol Comment

exported 심볼 바로 위에 빈 줄 없이 주석을 작성합니다.

```go
// New creates a new Evmc client connected to the given HTTP/HTTPS RPC endpoint.
func New(httpURL string, opts ...Options) (*Evmc, error) {
```

규칙:
- 첫 단어는 심볼 이름과 동일 (`New creates...`, `WithConnPool sets...`)
- 한 줄이면 충분, 복잡하면 여러 줄
- 다른 심볼을 참조할 때는 `[Evmc.Eth]` 형식 사용 (doc link)

## 3. Example Function

`_test.go` 파일에 `Example` 접두사로 함수를 작성하면, pkg.go.dev 문서에 실행 가능한 예제로 표시됩니다.

```go
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
```

### 네이밍 규칙

| 패턴 | 문서화 대상 |
|------|------------|
| `func Example()` | 패키지 전체 |
| `func ExampleNew()` | `New` 함수 |
| `func ExampleEvmc_Eth()` | `Evmc` 타입의 `Eth` 메서드 |
| `func ExampleNew_withOptions()` | `New` 함수 (두 번째 예제, 접미사로 구분) |

### Output 주석

- `// Output:` 주석을 넣으면 `go test`가 실제 출력과 비교하여 검증
- 주석이 없으면 컴파일만 되고 실행은 안 됨 (네트워크 필요한 예제에 적합)

## 4. 로컬에서 확인하기

```bash
# go doc 으로 패키지 문서 확인
go doc github.com/bbaktaeho/evmc

# 특정 심볼 확인
go doc github.com/bbaktaeho/evmc New

# pkgsite 로컬 서버로 pkg.go.dev 와 동일한 UI 확인
go install golang.org/x/pkgsite/cmd/pkgsite@latest
pkgsite -http=:8080
# http://localhost:8080/github.com/bbaktaeho/evmc 에서 확인
```

## References

- [Go Doc Comments](https://go.dev/doc/comment) - 공식 주석 문법 가이드
- [Testable Examples in Go](https://go.dev/blog/examples) - Example 함수 상세 설명
- [Godoc: documenting Go code](https://go.dev/blog/godoc) - GoDoc 개요
