package evmc

import (
	"context"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/ethereum/go-ethereum/rpc"
)

type contract struct {
	c caller
}

func (c *contract) Query(ctx context.Context, queryParams *evmctypes.QueryParams) (string, error) {
	var (
		result = new(string)
		params = []interface{}{queryParams, evmctypes.ParseBlockAndTag(queryParams.NumOrTag)}
	)
	if err := c.c.call(ctx, result, ethCall, params...); err != nil {
		return "", err
	}
	return *result, nil
}

func (c *contract) BatchQuery(
	ctx context.Context,
	batchQueryParams []*evmctypes.QueryParams,
) ([]*evmctypes.QueryResp, error) {
	var (
		size     = len(batchQueryParams)
		elements = make([]rpc.BatchElem, size)
		results  = make([]*evmctypes.QueryResp, size)
	)
	for i := range elements {
		numOrTag := evmctypes.ParseBlockAndTag(batchQueryParams[i].NumOrTag)
		elements[i] = rpc.BatchElem{
			Method: ethCall.String(),
			Args: []interface{}{
				batchQueryParams[i],
				numOrTag,
			},
			Result: &results[i].Result,
		}
	}
	if err := c.c.batchCall(ctx, elements); err != nil {
		return nil, err
	}
	for i, el := range elements {
		results[i].Error = el.Error
	}
	return results, nil
}
