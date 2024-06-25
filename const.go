package evmc

const (
	ZeroValue   = "0x0"
	ZeroData    = "0x"
	ZeroAddress = "0x0000000000000000000000000000000000000000"
	ZeroHash    = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

type ClientName string

const (
	Geth   ClientName = "geth"
	Erigon ClientName = "erigon"
	Bor    ClientName = "bor"
	Besu   ClientName = "besu"
	Nitro  ClientName = "nitro"
)

func (c ClientName) String() string {
	return string(c)
}
