package evmc

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ethNamespace struct {
	c caller
}

func (e *ethNamespace) GetChainID() (uint64, error) {
	result := new(string)
	if err := e.c.call(context.Background(), result, ethChainID); err != nil {
		return 0, err
	}
	id, err := hexutil.DecodeUint64(*result)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (e *ethNamespace) GetblockNumber() uint64 {
	return 0
}
