package evmc

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEvmcForKaia() *Evmc {
	rpcURL := "https://public-en.node.kaia.io"
	client, err := New(rpcURL)
	if err != nil {
		panic(err)
	}
	return client
}

func newTestKaiaNamespace() *kaiaNamespace {
	return &kaiaNamespace{c: testEvmcForKaia()}
}

func Test_kaiaNamespace_BlockNumber(t *testing.T) {
	ns := newTestKaiaNamespace()
	latest, err := ns.BlockNumber()
	assert.NoError(t, err)
	assert.NotZero(t, latest)
	t.Logf("latest block number: %d", latest)
}

func Test_kaiaNamespace_GetBlockIncTxRange(t *testing.T) {
	ns := newTestKaiaNamespace()
	startBlock := uint64(184994449)
	endBlock := uint64(184994450)
	blocks, err := ns.getBlockIncTxRange(context.Background(), startBlock, endBlock)
	assert.NoError(t, err)
	assert.NotNil(t, blocks)
	// json indent
	json, err := json.MarshalIndent(blocks, "", "  ")
	assert.NoError(t, err)
	t.Logf("blocks: %s", string(json))
}

func Test_kaiaNamespace_GetBlockByHash(t *testing.T) {
	ns := newTestKaiaNamespace()
	hash := "0x0e696ddb2e525893e03a3ee62f7582ffe253e5aeefe2ca569f8ebf0b9622a736"
	block, err := ns.GetBlockByHashIncTx(hash)
	assert.NoError(t, err)
	assert.NotNil(t, block)
	assert.Equal(t, hash, block.Hash)
	assert.Equal(t, len(block.Transactions), 1)
	fmt.Println(block.Number)
	assert.Equal(t, block.Number, uint64(184997666))
	_, err = json.MarshalIndent(block, "", "  ")
	assert.NoError(t, err)

}

func Test_kaiaNamespace_GetBlockByNumber(t *testing.T) {
	ns := newTestKaiaNamespace()
	number := uint64(184994449)
	block, err := ns.GetBlockByNumberIncTx(number)
	assert.NoError(t, err)
	assert.NotNil(t, block)
	// json indent
	data, err := json.MarshalIndent(block, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(data))
}

func Test_kaiaNamespace_GetTransactionByHash(t *testing.T) {
	ns := newTestKaiaNamespace()
	hash := "0xcf72eaf651ac7f19b694654c5af37e1e3a6cd9f93627a6eb425b28bb3627df76"
	tx, err := ns.GetTransactionByHash(hash)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	// json indent
	data, err := json.MarshalIndent(tx, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(data))
}

func Test_kaiaNamespace_GetBlockReceipts(t *testing.T) {
	ns := newTestKaiaNamespace()
	number := uint64(184994449)
	receipts, err := ns.GetBlockReceipts(number)
	assert.NoError(t, err)
	assert.NotNil(t, receipts)
	// data indent
	data, err := json.MarshalIndent(receipts, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(data))
}

func Test_kaiaNamespace_GetTransactionReceipt(t *testing.T) {
	ns := newTestKaiaNamespace()
	hash := "0xb82f1370cc25dacb2880ad315bfff8ff7568e622217ec7c16b3278c1a076aba8" // 테스트용 트랜잭션 해시
	receipt, err := ns.GetTransactionReceipt(hash)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	// json indent
	json, err := json.MarshalIndent(receipt, "", "  ")
	assert.NoError(t, err)
	t.Logf("receipt: %s", string(json))
}

func Test_kaiaNamespace_GetRewards(t *testing.T) {
	ns := newTestKaiaNamespace()
	blockNumber := uint64(184930758)
	resp, err := ns.GetRewards(blockNumber)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	json, err := json.MarshalIndent(resp, "", "  ")
	assert.NoError(t, err)
	t.Logf("rewards: %s", string(json))
}

func Test_kaiaNamespace_GetRewardsRange(t *testing.T) {
	ns := newTestKaiaNamespace()
	startBlock := uint64(184930558)
	endBlock := uint64(184930558)
	resp, err := ns.GetRewardsRange(startBlock, endBlock)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(resp), int(endBlock-startBlock+1))

	b, _ := json.MarshalIndent(resp[0], "", "  ")
	t.Logf("rewards: %s", string(b))
}
