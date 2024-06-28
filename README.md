# evmc

### install

- go version

  required Go version (v1.22 or later)

- install

  ```bash
  go get github.com/bbaktaeho/evmc
  ```

### Usage

```go
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
