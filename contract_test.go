package evmc

import (
	"math/rand"
	"testing"

	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/stretchr/testify/assert"
)

func testEvmcForContract(workers, batchItems int) contract {
	rpcURL := "https://ethereum-mainnet.nodit.io/<api-key>"
	client, err := New(rpcURL, WithBatchCallWorkers(workers), WithMaxBatchItems(batchItems))
	if err != nil {
		panic(err)
	}
	return contract{c: client}
}

func Test_contract_BatchQueries(t *testing.T) {
	var (
		workers, batchItems            = 20, 50
		items                          = 1000
		minBlockNumber, maxBlockNumber = 21000000, 22000000
		chain                          = testEvmcForContract(workers, batchItems)
		batchQueryParams               []*evmctypes.QueryParams
	)

	for range items {
		batchQueryParams = append(batchQueryParams, &evmctypes.QueryParams{
			To:       "0xdac17f958d2ee523a2206206994597c13d831ec7",
			Data:     GenerateERC20BalanceOf("0xDBe46A02322e636b92296954637E1D7dB9D5ed26"),
			NumOrTag: evmctypes.FormatNumber(uint64(rand.Intn(maxBlockNumber-minBlockNumber) + minBlockNumber)),
		})
	}
	res, err := chain.BatchQueries(batchQueryParams)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(batchQueryParams), len(res))

	for _, r := range res {
		if r.Error != nil {
			t.Error(r.Error)
			continue
		}
	}
}
