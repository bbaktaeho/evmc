package evmc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BatchCallWithContext_chunking(t *testing.T) {
	t.Parallel()

	var callCount atomic.Int64
	mock := newMockRPCServer(t)
	mock.on("eth_blockNumber", func(_ json.RawMessage) interface{} {
		callCount.Add(1)
		return "0x10"
	})

	client, err := New(mock.url(), WithMaxBatchItems(3), WithBatchCallWorkers(2))
	require.NoError(t, err)
	defer client.Close()

	elements := make([]rpc.BatchElem, 10)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: "eth_blockNumber",
			Result: new(string),
		}
	}

	err = client.BatchCallWithContext(context.Background(), elements, 0)
	require.NoError(t, err)

	for i, elem := range elements {
		assert.NoError(t, elem.Error, "element %d should not have an error", i)
		assert.Equal(t, "0x10", *elem.Result.(*string), "element %d result mismatch", i)
	}
	assert.Equal(t, int64(10), callCount.Load(), "all 10 requests should be dispatched")
}

func Test_BatchCallWithContext_singleChunk(t *testing.T) {
	t.Parallel()

	mock := newMockRPCServer(t)
	mock.on("eth_chainId", func(_ json.RawMessage) interface{} {
		return "0x1"
	})

	client, err := New(mock.url(), WithMaxBatchItems(100))
	require.NoError(t, err)
	defer client.Close()

	elements := make([]rpc.BatchElem, 5)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: "eth_chainId",
			Result: new(string),
		}
	}

	err = client.BatchCallWithContext(context.Background(), elements, 1)
	require.NoError(t, err)

	for i, elem := range elements {
		assert.NoError(t, elem.Error)
		assert.Equal(t, "0x1", *elem.Result.(*string), "element %d", i)
	}
}

func Test_BatchCallWithContext_errorPropagation(t *testing.T) {
	t.Parallel()

	mock := newMockRPCServer(t)
	// Register nothing – all calls will get "method not found" error.

	client, err := New(mock.url(), WithMaxBatchItems(2), WithBatchCallWorkers(2))
	require.NoError(t, err)
	defer client.Close()

	elements := make([]rpc.BatchElem, 4)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: "nonexistent_method",
			Result: new(string),
		}
	}

	err = client.BatchCallWithContext(context.Background(), elements, 0)
	assert.NoError(t, err, "batch call itself should succeed; per-element errors are in BatchElem.Error")

	for i, elem := range elements {
		assert.Error(t, elem.Error, "element %d should have a per-element error", i)
	}
}

func Test_BatchCallWithContext_defaultWorkers(t *testing.T) {
	t.Parallel()

	mock := newMockRPCServer(t)
	mock.on("eth_blockNumber", func(_ json.RawMessage) interface{} {
		return "0xff"
	})

	client, err := New(mock.url(), WithMaxBatchItems(2), WithBatchCallWorkers(4))
	require.NoError(t, err)
	defer client.Close()

	elements := make([]rpc.BatchElem, 8)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: "eth_blockNumber",
			Result: new(string),
		}
	}

	err = client.BatchCallWithContext(context.Background(), elements, 0)
	require.NoError(t, err)

	for i, elem := range elements {
		assert.NoError(t, elem.Error)
		assert.Equal(t, "0xff", *elem.Result.(*string), "element %d", i)
	}
}

func Test_BatchCallWithContext_contextCancellation(t *testing.T) {
	t.Parallel()

	mock := newMockRPCServer(t)
	mock.on("eth_blockNumber", func(_ json.RawMessage) interface{} {
		return "0x1"
	})

	client, err := New(mock.url(), WithMaxBatchItems(1), WithBatchCallWorkers(1))
	require.NoError(t, err)
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	elements := make([]rpc.BatchElem, 5)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: "eth_blockNumber",
			Result: new(string),
		}
	}

	err = client.BatchCallWithContext(ctx, elements, 1)
	assert.Error(t, err, "cancelled context should cause an error")
}

func Test_BatchCall_wrapper(t *testing.T) {
	t.Parallel()

	mock := newMockRPCServer(t)
	mock.on("eth_gasPrice", func(_ json.RawMessage) interface{} {
		return "0x3b9aca00"
	})

	client, err := New(mock.url(), WithMaxBatchItems(5))
	require.NoError(t, err)
	defer client.Close()

	elements := make([]rpc.BatchElem, 3)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: "eth_gasPrice",
			Result: new(string),
		}
	}

	err = client.BatchCall(elements, 1)
	require.NoError(t, err)

	for i, elem := range elements {
		assert.NoError(t, elem.Error)
		assert.Equal(t, "0x3b9aca00", *elem.Result.(*string), "element %d", i)
	}
}

func Test_BatchCallWithContext_largePayload(t *testing.T) {
	t.Parallel()

	mock := newMockRPCServer(t)
	mock.on("eth_getBlockByNumber", func(params json.RawMessage) interface{} {
		var args []interface{}
		_ = json.Unmarshal(params, &args)
		return blockJSON(fmt.Sprintf("%v", args[0]), "0xhash", true)
	})

	client, err := New(mock.url(), WithMaxBatchItems(10), WithBatchCallWorkers(3))
	require.NoError(t, err)
	defer client.Close()

	const total = 50
	elements := make([]rpc.BatchElem, total)
	for i := range elements {
		elements[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{fmt.Sprintf("0x%x", i+1), true},
			Result: new(json.RawMessage),
		}
	}

	err = client.BatchCallWithContext(context.Background(), elements, 0)
	require.NoError(t, err)

	for i, elem := range elements {
		assert.NoError(t, elem.Error, "element %d should not error", i)
		assert.NotNil(t, elem.Result, "element %d should have a result", i)
	}
}
