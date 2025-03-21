package evmc

import (
	"context"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/ethereum/go-ethereum/rpc"
)

type contract struct {
	c caller
}

func (c *contract) Query(ctx context.Context, queryParams *evmctypes.QueryParams) (*evmctypes.QueryResp, error) {
	var (
		result = new(string)
		params = []interface{}{queryParams, evmctypes.ParseBlockAndTag(queryParams.NumOrTag)}
	)
	if err := c.c.call(ctx, result, EthCall, params...); err != nil {
		return nil, err
	}
	return &evmctypes.QueryResp{
		To:     queryParams.To,
		Data:   queryParams.Data,
		Result: *result,
	}, nil
}

// deprecated
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
		results[i] = &evmctypes.QueryResp{To: batchQueryParams[i].To, Data: batchQueryParams[i].Data}
		numOrTag := evmctypes.ParseBlockAndTag(batchQueryParams[i].NumOrTag)
		elements[i] = rpc.BatchElem{
			Method: EthCall.String(),
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

func (c *contract) BatchQueries(batchQueryParams []*evmctypes.QueryParams) ([]*evmctypes.QueryResp, error) {
	return c.BatchQueriesWithContext(context.Background(), batchQueryParams)
}

func (c *contract) BatchQueriesWithContext(
	ctx context.Context,
	batchQueryParams []*evmctypes.QueryParams,
) ([]*evmctypes.QueryResp, error) {
	var (
		size     = len(batchQueryParams)
		elements = make([]rpc.BatchElem, size)
		results  = make([]*evmctypes.QueryResp, size)
	)
	for i := range elements {
		results[i] = &evmctypes.QueryResp{To: batchQueryParams[i].To, Data: batchQueryParams[i].Data}
		numOrTag := batchQueryParams[i].NumOrTag.String()
		elements[i] = rpc.BatchElem{
			Method: EthCall.String(),
			Args: []interface{}{
				batchQueryParams[i],
				numOrTag,
			},
			Result: &results[i].Result,
		}
	}
	if err := c.c.BatchCallWithContext(ctx, elements, -1); err != nil {
		return nil, err
	}
	for i, el := range elements {
		results[i].Error = el.Error
	}
	return results, nil
}
