package evmc

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_contract_mock_Query(t *testing.T) {
	mock := newMockRPCServer(t)
	mock.on("eth_call", func(params json.RawMessage) interface{} {
		return "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000"
	})

	client := testEvmc(mock.url())
	resp, err := client.Contract().Query(context.Background(), &evmctypes.QueryParams{
		To:       "0xcontract",
		Data:     "0x18160ddd",
		NumOrTag: evmctypes.Latest,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "0xcontract", resp.To)
	assert.Equal(t, "0x18160ddd", resp.Data)
	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000", resp.Result)
}

func Test_contract_mock_BatchQueries(t *testing.T) {
	mock := newMockRPCServer(t)
	callIndex := 0
	results := []string{
		"0x0000000000000000000000000000000000000000000000000de0b6b3a7640000",
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000000000000000000000000000000",
	}
	mock.on("eth_call", func(params json.RawMessage) interface{} {
		result := results[callIndex%len(results)]
		callIndex++
		return result
	})

	client := testEvmc(mock.url())
	queries := []*evmctypes.QueryParams{
		{To: "0xcontract1", Data: "0x18160ddd", NumOrTag: evmctypes.Latest},
		{To: "0xcontract2", Data: "0x70a08231", NumOrTag: evmctypes.Latest},
		{To: "0xcontract3", Data: "0x06fdde03", NumOrTag: evmctypes.Latest},
	}

	resps, err := client.Contract().BatchQueries(queries)
	require.NoError(t, err)
	require.Len(t, resps, 3)
	assert.Equal(t, "0xcontract1", resps[0].To)
	assert.Equal(t, "0xcontract2", resps[1].To)
	assert.Equal(t, "0xcontract3", resps[2].To)
	assert.Equal(t, results[0], resps[0].Result)
}

func Test_contract_mock_BatchQueries_emptyList(t *testing.T) {
	mock := newMockRPCServer(t)
	client := testEvmc(mock.url())

	resps, err := client.Contract().BatchQueries([]*evmctypes.QueryParams{})
	require.NoError(t, err)
	assert.Empty(t, resps)
}

func Test_contract_mock_Query_withBlockNumber(t *testing.T) {
	mock := newMockRPCServer(t)
	mock.on("eth_call", func(params json.RawMessage) interface{} {
		var args []json.RawMessage
		json.Unmarshal(params, &args)
		// 두 번째 파라미터가 블록 태그/번호여야 함
		assert.Len(t, args, 2)
		return "0x0000000000000000000000000000000000000000000000000000000000000001"
	})

	client := testEvmc(mock.url())
	resp, err := client.Contract().Query(context.Background(), &evmctypes.QueryParams{
		To:       "0xcontract",
		Data:     "0x18160ddd",
		NumOrTag: evmctypes.FormatNumber(12345),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
}
