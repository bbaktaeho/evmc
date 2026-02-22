package evmctypes

import "github.com/ethereum/go-ethereum/common/hexutil"

type BlockAndTag string

const (
	Pending   BlockAndTag = "pending"
	Earliest  BlockAndTag = "earliest"
	Latest    BlockAndTag = "latest"
	Safe      BlockAndTag = "safe"
	Finalized BlockAndTag = "finalized"
)

func (b BlockAndTag) Uint64() (uint64, error) {
	return hexutil.DecodeUint64(string(b))
}

func (b BlockAndTag) String() string {
	return string(b)
}

// FormatNumber returns a BlockAndTag representing the given block number in hex.
func FormatNumber(number uint64) BlockAndTag {
	return BlockAndTag(hexutil.EncodeUint64(number))
}

// ParseBlockAndTag returns the string form of a block tag or number.
// If v is a BlockAndTag or string, it is returned as-is; otherwise "latest" is returned.
func ParseBlockAndTag(v interface{}) string {
	switch val := v.(type) {
	case BlockAndTag:
		return val.String()
	case string:
		return val
	default:
		return Latest.String()
	}
}
