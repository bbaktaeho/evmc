package evmctypes

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type BlockAndTag string

const (
	Pending   BlockAndTag = "pending"
	Earliest  BlockAndTag = "earliest"
	Latest    BlockAndTag = "latest"
	Safe      BlockAndTag = "safe"
	Finalized BlockAndTag = "finalized"
)

func (b BlockAndTag) String() string {
	return string(b)
}

func FormatNumber(number uint64) BlockAndTag {
	return BlockAndTag(hexutil.EncodeUint64(number))
}

func ParseBlockAndTag(blockAndTag interface{}) string {
	if reflect.TypeOf(blockAndTag).Kind() != reflect.String {
		return Latest.String()
	}
	tag, ok := blockAndTag.(BlockAndTag)
	if !ok {
		return blockAndTag.(string)
	}
	return tag.String()
}
