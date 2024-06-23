package evmc

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ethNamespace struct {
	c caller
}

func (e *ethNamespace) ChainID() (uint64, error) {
	return e.chainID(context.Background())
}

func (e *ethNamespace) ChainIDWithContext(ctx context.Context) (uint64, error) {
	return e.chainID(ctx)
}

func (e *ethNamespace) chainID(ctx context.Context) (uint64, error) {
	result := new(string)
	if err := e.c.call(ctx, result, ethChainID); err != nil {
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
