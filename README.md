# evmc (WIP)

evmc is an abbreviation of Ethereum Virtual Machine Client, and I wanted to express it as simply as a Go. But evmc is not simple.

I'm trying to create a more user-friendly client using the rpc package of go-thereum(geth).
It also supports namespace for data analysis such as debug and trace and provides features for standard tokens.

### Install

- go version

  required Go version (v1.22 or later)

- install

  ```bash
  go get github.com/bbaktaeho/evmc
  ```

### Usage

```go
package main

import "github.com/bbaktaeho/evmc"

func main() {
    client, err := evmc.New("http://localhost:8545")
	if err != nil {
		panic(err)
	}

    latest, err := client.Eth().BlockNumber()
    if err != nil {
		panic(err)
	}

    println(latest)
}
```
