package evmctypes

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -update 플래그: 메인넷 RPC를 호출해 골든 파일을 재생성한다.
// 사용 예: go test ./evmctypes/... -run TestGolden -update
var update = flag.Bool("update", false, "regenerate golden files from live mainnet RPC (requires access to mainnet node)")

const (
	goldenDir     = "../testdata/mainnet"
	mainnetRPCURL = "http://61.111.3.69:18014"

	// Ethereum mainnet block 20,000,000 (0x1312d00) 기준 레퍼런스 값
	refBlockNumber   = uint64(20_000_000)
	refBlockHex      = "0x1312d00"
	refBlockHash     = "0xd24fd73f794058a3807db926d8898c6481e902b7edb91ce0d479d6760f276183"
	refBlockMiner    = "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
	refBlockTS       = uint64(1717281407)
	refBlockTxCount  = 134
	refBlockWithdraw = 16

	refEip1559TxHash = "0xbb4b3fc2b746877dce70862850602f1d19bd890ab4db47e6b7ee1da1fe578a0d"
	refEip1559TxFrom = "0xae2fc483527b8ef99eb5d9b44875f005ba1fae13"
	refEip1559TxTo   = "0x6b75d8af000000e20b7a7ddf000ba900b4009a80"

	refBlobTxHash    = "0x0ff07f37baa7fa26bb7de3d3fc63002bf0acf3295bdab7f67c108c0d1a3bff15"
	refBlobTxFrom    = "0x000000633b68f5d8d3a86593ebb815b4663bcbe0"
	refBlobVersioned = "0x017ba4bd9c166498865a3d08618e333ee84812941b5c3a356971b4a6ffffa574"

	refUSDTAddress = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	// getLogs에서 반환된 주소는 소문자
	refUSDTAddressLower = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	// Transfer 이벤트 topic
	refTransferTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	// callTracer와 동일한 구조를 반환하되 input/output 필드를 제거하는 커스텀 JS 트레이서
	customJSTracer = `{
  callstack: [{}],
  enter: function(f) {
    this.callstack.push({
      type: f.getType(),
      from: toHex(f.getFrom()),
      to: toHex(f.getTo()),
      gas: "0x" + f.getGas().toString(16),
      value: "0x0"
    });
  },
  exit: function(r) {
    var c = this.callstack.pop();
    c.gasUsed = "0x" + r.getGasUsed().toString(16);
    var p = this.callstack[this.callstack.length - 1];
    if (!p.calls) p.calls = [];
    p.calls.push(c);
  },
  result: function(ctx) {
    var r = this.callstack[0];
    r.type = ctx.type;
    r.from = toHex(ctx.from);
    r.to = toHex(ctx.to);
    r.gas = "0x" + ctx.gas.toString(16);
    r.gasUsed = "0x" + ctx.gasUsed.toString(16);
    r.value = "0x" + ctx.value.toString(16);
    return r;
  },
  fault: function() {}
}`
)

// loadGolden reads a golden file. Fails if file not found.
func loadGolden(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(goldenDir, name))
	require.NoError(t, err, "golden file not found: %s\nRun: go test ./evmctypes/... -run TestGolden -update", name)
	return data
}

// rpcResult makes a raw JSON-RPC call and returns the result bytes.
func rpcResult(t *testing.T, method string, params ...interface{}) json.RawMessage {
	t.Helper()
	body, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	})
	require.NoError(t, err)
	resp, err := http.Post(mainnetRPCURL, "application/json", bytes.NewReader(body))
	require.NoError(t, err, "cannot reach mainnet RPC at %s", mainnetRPCURL)
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var rpcResp struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(raw, &rpcResp))
	require.Nil(t, rpcResp.Error, "RPC error for %s: %v", method, rpcResp.Error)
	return rpcResp.Result
}

// saveGolden writes pretty-printed JSON to a golden file.
func saveGolden(t *testing.T, name string, data json.RawMessage) {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, json.Indent(&buf, data, "", "  "))
	path := filepath.Join(goldenDir, name)
	require.NoError(t, os.WriteFile(path, buf.Bytes(), 0644))
	t.Logf("updated golden file: %s", path)
}

// TestGolden_Block tests Block and BlockIncTx unmarshaling from mainnet data.
func TestGolden_Block(t *testing.T) {
	if *update {
		saveGolden(t, "eth_getBlockByNumber.json",
			rpcResult(t, "eth_getBlockByNumber", refBlockHex, false))
		saveGolden(t, "eth_getBlockByHash.json",
			rpcResult(t, "eth_getBlockByHash", refBlockHash, false))
		saveGolden(t, "eth_getBlockByNumber_with_tx.json",
			rpcResult(t, "eth_getBlockByNumber", refBlockHex, true))
	}

	t.Run("getBlockByNumber", func(t *testing.T) {
		var block Block
		require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getBlockByNumber.json"), &block))

		assert.Equal(t, refBlockNumber, block.Number)
		assert.Equal(t, refBlockHash, block.Hash)
		assert.Equal(t, refBlockMiner, block.Miner)
		assert.Equal(t, refBlockTS, block.Timestamp)
		assert.Equal(t, "0x1c9c380", block.GasLimit)
		assert.Equal(t, "0xa9371c", block.GasUsed)
		assert.Equal(t, "0x0", block.Difficulty)
		// EIP-1559 필드
		require.NotNil(t, block.BaseFeePerGas)
		assert.Equal(t, "0x12643ff14", *block.BaseFeePerGas)
		// EIP-4844 필드
		require.NotNil(t, block.BlobGasUsed)
		assert.Equal(t, "0x20000", *block.BlobGasUsed)
		require.NotNil(t, block.ExcessBlobGas)
		assert.Equal(t, "0x0", *block.ExcessBlobGas)
		// EIP-4895 필드
		require.NotNil(t, block.WithdrawalsRoot)
		assert.NotEmpty(t, *block.WithdrawalsRoot)
		// EIP-4788 필드
		require.NotNil(t, block.ParentBeaconBlockRoot)
		assert.NotEmpty(t, *block.ParentBeaconBlockRoot)
		// transactions: tx hash 배열
		assert.Len(t, block.Transactions, refBlockTxCount)
		assert.Equal(t, refEip1559TxHash, block.Transactions[0])
		// withdrawals
		assert.Len(t, block.Withdrawals, refBlockWithdraw)
	})

	t.Run("getBlockByHash", func(t *testing.T) {
		var block Block
		require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getBlockByHash.json"), &block))

		assert.Equal(t, refBlockNumber, block.Number)
		assert.Equal(t, refBlockHash, block.Hash)
		assert.Equal(t, refBlockMiner, block.Miner)
		assert.Len(t, block.Transactions, refBlockTxCount)
	})

	t.Run("getBlockByNumber_withFullTx", func(t *testing.T) {
		var block BlockIncTx
		require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getBlockByNumber_with_tx.json"), &block))

		assert.Equal(t, refBlockNumber, block.Number)
		assert.Equal(t, refBlockHash, block.Hash)
		assert.Len(t, block.Transactions, refBlockTxCount)

		// 첫 번째 트랜잭션 검증 (EIP-1559 type 2)
		require.NotEmpty(t, block.Transactions)
		firstTx := block.Transactions[0]
		assert.Equal(t, "0x2", firstTx.Type)
		assert.Equal(t, refEip1559TxHash, firstTx.Hash)
		assert.Equal(t, refEip1559TxFrom, firstTx.From)
		assert.Equal(t, refBlockNumber, firstTx.BlockNumber)
	})
}

// TestGolden_Transaction tests Transaction unmarshaling for EIP-1559 and EIP-4844 (blob) txs.
func TestGolden_Transaction(t *testing.T) {
	if *update {
		saveGolden(t, "eth_getTransactionByHash_eip1559.json",
			rpcResult(t, "eth_getTransactionByHash", refEip1559TxHash))
		saveGolden(t, "eth_getTransactionByHash_blob.json",
			rpcResult(t, "eth_getTransactionByHash", refBlobTxHash))
	}

	t.Run("EIP-1559_type2", func(t *testing.T) {
		var tx Transaction
		require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getTransactionByHash_eip1559.json"), &tx))

		assert.Equal(t, "0x2", tx.Type)
		assert.Equal(t, refEip1559TxHash, tx.Hash)
		assert.Equal(t, refEip1559TxFrom, tx.From)
		assert.Equal(t, refEip1559TxTo, tx.To)
		assert.Equal(t, refBlockNumber, tx.BlockNumber)
		assert.Equal(t, refBlockHash, tx.BlockHash)
		assert.Equal(t, uint64(0), tx.TransactionIndex)
		// EIP-155
		require.NotNil(t, tx.ChainID)
		assert.Equal(t, "0x1", *tx.ChainID)
		// EIP-1559 필드
		require.NotNil(t, tx.MaxFeePerGas)
		require.NotNil(t, tx.MaxPriorityFeePerGas)
		assert.NotEmpty(t, *tx.MaxFeePerGas)
		assert.NotEmpty(t, *tx.MaxPriorityFeePerGas)
		// Blob 필드 없음
		assert.Nil(t, tx.MaxFeePerBlobGas)
		assert.Empty(t, tx.BlobVersionedHashes)
		// Value가 정상 파싱됨
		assert.False(t, tx.Value.IsNegative())
	})

	t.Run("EIP-4844_blob_type3", func(t *testing.T) {
		var tx Transaction
		require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getTransactionByHash_blob.json"), &tx))

		assert.Equal(t, "0x3", tx.Type)
		assert.Equal(t, refBlobTxHash, tx.Hash)
		assert.Equal(t, refBlobTxFrom, tx.From)
		assert.Equal(t, refBlockNumber, tx.BlockNumber)
		// EIP-4844 필드
		require.NotNil(t, tx.MaxFeePerBlobGas)
		assert.Equal(t, "0x3b9aca00", *tx.MaxFeePerBlobGas)
		require.Len(t, tx.BlobVersionedHashes, 1)
		assert.Equal(t, refBlobVersioned, tx.BlobVersionedHashes[0])
		// yParity 존재
		require.NotNil(t, tx.YParity)
	})
}

// TestGolden_Receipt tests Receipt unmarshaling for EIP-1559 and EIP-4844 (blob) receipts.
func TestGolden_Receipt(t *testing.T) {
	if *update {
		saveGolden(t, "eth_getTransactionReceipt_eip1559.json",
			rpcResult(t, "eth_getTransactionReceipt", refEip1559TxHash))
		saveGolden(t, "eth_getTransactionReceipt_blob.json",
			rpcResult(t, "eth_getTransactionReceipt", refBlobTxHash))
	}

	t.Run("EIP-1559_type2", func(t *testing.T) {
		var receipt Receipt
		require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getTransactionReceipt_eip1559.json"), &receipt))

		assert.Equal(t, "0x2", receipt.Type)
		assert.Equal(t, refEip1559TxHash, receipt.TransactionHash)
		assert.Equal(t, refEip1559TxFrom, receipt.From)
		assert.Equal(t, refEip1559TxTo, receipt.To)
		assert.Equal(t, refBlockNumber, receipt.BlockNumber)
		assert.Equal(t, refBlockHash, receipt.BlockHash)
		assert.Equal(t, uint64(0), receipt.TransactionIndex)
		// 성공 상태
		require.NotNil(t, receipt.Status)
		assert.Equal(t, "0x1", *receipt.Status)
		// EIP-1559 effectiveGasPrice
		assert.NotEmpty(t, receipt.EffectiveGasPrice)
		// Blob 필드 없음
		assert.Nil(t, receipt.BlobGasUsed)
		assert.Nil(t, receipt.BlobGasPrice)
	})

	t.Run("EIP-4844_blob_type3", func(t *testing.T) {
		var receipt Receipt
		require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getTransactionReceipt_blob.json"), &receipt))

		assert.Equal(t, "0x3", receipt.Type)
		assert.Equal(t, refBlobTxHash, receipt.TransactionHash)
		assert.Equal(t, refBlockNumber, receipt.BlockNumber)
		// 성공 상태
		require.NotNil(t, receipt.Status)
		assert.Equal(t, "0x1", *receipt.Status)
		// EIP-4844 blob 필드
		require.NotNil(t, receipt.BlobGasUsed)
		assert.Equal(t, "0x20000", *receipt.BlobGasUsed)
		require.NotNil(t, receipt.BlobGasPrice)
		assert.Equal(t, "0x1", *receipt.BlobGasPrice)
		// logs 4개
		assert.Len(t, receipt.Logs, 4)
	})
}

// TestGolden_BlockReceipts tests []*Receipt unmarshaling for all receipts in a block.
func TestGolden_BlockReceipts(t *testing.T) {
	if *update {
		saveGolden(t, "eth_getBlockReceipts.json",
			rpcResult(t, "eth_getBlockReceipts", refBlockHex))
	}

	var receipts []*Receipt
	require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getBlockReceipts.json"), &receipts))

	assert.Len(t, receipts, refBlockTxCount)

	// 첫 번째: EIP-1559 타입
	require.NotEmpty(t, receipts)
	assert.Equal(t, "0x2", receipts[0].Type)
	assert.Equal(t, refEip1559TxHash, receipts[0].TransactionHash)

	// blob tx 영수증 확인 (type 0x3)
	var blobReceipt *Receipt
	for _, r := range receipts {
		if r.Type == "0x3" {
			blobReceipt = r
			break
		}
	}
	require.NotNil(t, blobReceipt, "no blob receipt found in block")
	assert.Equal(t, refBlobTxHash, blobReceipt.TransactionHash)
	require.NotNil(t, blobReceipt.BlobGasUsed)
	assert.Equal(t, "0x20000", *blobReceipt.BlobGasUsed)
}

// TestGolden_Log tests []*Log unmarshaling from eth_getLogs.
func TestGolden_Log(t *testing.T) {
	if *update {
		saveGolden(t, "eth_getLogs.json",
			rpcResult(t, "eth_getLogs", map[string]interface{}{
				"fromBlock": refBlockHex,
				"toBlock":   refBlockHex,
				"address":   refUSDTAddress,
			}))
	}

	var logs []*Log
	require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getLogs.json"), &logs))

	assert.Len(t, logs, 12)

	// 모든 로그가 USDT 컨트랙트에서 발생
	for _, l := range logs {
		assert.Equal(t, refUSDTAddressLower, l.Address)
		assert.Equal(t, refBlockNumber, l.BlockNumber)
		assert.Equal(t, refBlockHash, l.BlockHash)
		assert.Equal(t, refBlockTS, l.BlockTimestamp)
		assert.NotEmpty(t, l.Topics)
		assert.NotEmpty(t, l.TransactionHash)
	}

	// Transfer 이벤트 확인 (첫 번째 로그)
	assert.Equal(t, refTransferTopic, logs[0].Topics[0])
}

// TestGolden_FeeHistory tests FeeHistory unmarshaling including EIP-4844 blob fee fields.
func TestGolden_FeeHistory(t *testing.T) {
	if *update {
		saveGolden(t, "eth_feeHistory.json",
			rpcResult(t, "eth_feeHistory", 4, refBlockHex, []float64{25, 75}))
	}

	var fh FeeHistory
	require.NoError(t, json.Unmarshal(loadGolden(t, "eth_feeHistory.json"), &fh))

	// oldest block: blockCount=4 이전 블록 (19999997)
	assert.Equal(t, uint64(19_999_997), fh.OldestBlock)
	// baseFeePerGas: blockCount+1 개 (현재 블록 + 다음 블록 예측 포함)
	assert.Len(t, fh.BaseFeePerGas, 5)
	assert.NotEmpty(t, fh.BaseFeePerGas[0])
	// EIP-4844 blob base fee
	assert.Len(t, fh.BaseFeePerBlobGas, 5)
	// gasUsedRatio
	assert.Len(t, fh.GasUsedRatio, 4)
	// blobGasUsedRatio
	assert.Len(t, fh.BlobGasUsedRatio, 4)
	// reward: 25th, 75th percentile
	assert.Len(t, fh.Reward, 4)
	for _, r := range fh.Reward {
		assert.Len(t, r, 2)
	}
}

// ─── debug namespace golden tests ────────────────────────────────────────────

// TestGolden_TraceTransaction_callTracer tests CallFrame unmarshaling from debug_traceTransaction with callTracer.
func TestGolden_TraceTransaction_callTracer(t *testing.T) {
	if *update {
		saveGolden(t, "debug_traceTransaction_callTracer.json",
			rpcResult(t, "debug_traceTransaction", refEip1559TxHash, map[string]interface{}{"tracer": "callTracer"}))
		saveGolden(t, "debug_traceTransaction_callTracer_withLog.json",
			rpcResult(t, "debug_traceTransaction", refEip1559TxHash, map[string]interface{}{
				"tracer":       "callTracer",
				"tracerConfig": map[string]interface{}{"withLog": true},
			}))
	}

	t.Run("callTracer", func(t *testing.T) {
		var frame CallFrame
		require.NoError(t, json.Unmarshal(loadGolden(t, "debug_traceTransaction_callTracer.json"), &frame))

		assert.Equal(t, "CALL", frame.Type)
		assert.Equal(t, refEip1559TxFrom, frame.From)
		require.NotNil(t, frame.To)
		assert.Equal(t, refEip1559TxTo, *frame.To)
		assert.NotEmpty(t, frame.Gas)
		assert.NotEmpty(t, frame.GasUsed)
		assert.Equal(t, "0x4e5d3", frame.GasUsed)
		assert.Len(t, frame.Calls, 8)
		// value 파싱
		require.NotNil(t, frame.Value)
		assert.False(t, frame.Value.IsNegative())
	})

	t.Run("callTracer_withLog", func(t *testing.T) {
		var frame CallFrame
		require.NoError(t, json.Unmarshal(loadGolden(t, "debug_traceTransaction_callTracer_withLog.json"), &frame))

		assert.Equal(t, "CALL", frame.Type)
		assert.Equal(t, refEip1559TxFrom, frame.From)
		assert.Len(t, frame.Calls, 8)
		// 하위 call에 로그가 존재
		var totalLogs int
		countLogs(&totalLogs, &frame)
		assert.Equal(t, 16, totalLogs)
	})
}

// countLogs recursively counts logs in a CallFrame tree.
func countLogs(count *int, frame *CallFrame) {
	*count += len(frame.Logs)
	for _, c := range frame.Calls {
		countLogs(count, c)
	}
}

// TestGolden_TraceTransaction_prestateTracer tests PrestateResult unmarshaling.
func TestGolden_TraceTransaction_prestateTracer(t *testing.T) {
	if *update {
		saveGolden(t, "debug_traceTransaction_prestateTracer.json",
			rpcResult(t, "debug_traceTransaction", refEip1559TxHash, map[string]interface{}{"tracer": "prestateTracer"}))
		saveGolden(t, "debug_traceTransaction_prestateTracer_diff.json",
			rpcResult(t, "debug_traceTransaction", refEip1559TxHash, map[string]interface{}{
				"tracer":       "prestateTracer",
				"tracerConfig": map[string]interface{}{"diffMode": true},
			}))
	}

	t.Run("prestateTracer", func(t *testing.T) {
		raw := loadGolden(t, "debug_traceTransaction_prestateTracer.json")
		result := &PrestateResult{RawMessage: raw}
		frame, err := result.ParseFrame()
		require.NoError(t, err)

		assert.Len(t, frame, 12, "should have 12 accounts in prestate")
		// from 계정이 존재해야 함
		fromAccount, ok := frame[refEip1559TxFrom]
		require.True(t, ok, "from account should exist in prestate")
		assert.NotNil(t, fromAccount.Balance)
	})

	t.Run("prestateTracer_diffMode", func(t *testing.T) {
		raw := loadGolden(t, "debug_traceTransaction_prestateTracer_diff.json")
		result := &PrestateResult{RawMessage: raw}
		diffFrame, err := result.ParseDiffFrame()
		require.NoError(t, err)

		require.NotNil(t, diffFrame)
		assert.Len(t, diffFrame.Pre, 11, "pre should have 11 accounts")
		assert.Len(t, diffFrame.Post, 11, "post should have 11 accounts")
	})
}

// TestGolden_TraceTransaction_flatCallTracer tests FlatCallFrame unmarshaling.
func TestGolden_TraceTransaction_flatCallTracer(t *testing.T) {
	if *update {
		saveGolden(t, "debug_traceTransaction_flatCallTracer.json",
			rpcResult(t, "debug_traceTransaction", refEip1559TxHash, map[string]interface{}{"tracer": "flatCallTracer"}))
	}

	var frames []*FlatCallFrame
	require.NoError(t, json.Unmarshal(loadGolden(t, "debug_traceTransaction_flatCallTracer.json"), &frames))

	assert.Len(t, frames, 21)

	// 첫 번째 프레임: root call
	first := frames[0]
	assert.Equal(t, "call", first.Type)
	require.NotNil(t, first.Action.From)
	assert.Equal(t, refEip1559TxFrom, *first.Action.From)
	require.NotNil(t, first.Action.To)
	assert.Equal(t, refEip1559TxTo, *first.Action.To)
	assert.Equal(t, uint64(8), first.Subtraces)
	assert.Empty(t, first.TraceAddress)
	require.NotNil(t, first.Result)
	assert.Equal(t, "0x288ab", *first.Result.GasUsed)
}

// TestGolden_TraceCall tests CallFrame unmarshaling from debug_traceCall with callTracer.
func TestGolden_TraceCall(t *testing.T) {
	if *update {
		saveGolden(t, "debug_traceCall_callTracer.json",
			rpcResult(t, "debug_traceCall",
				map[string]interface{}{
					"from": "0x0000000000000000000000000000000000000000",
					"to":   refUSDTAddress,
					"data": "0x18160ddd",
				},
				refBlockHex,
				map[string]interface{}{"tracer": "callTracer"},
			))
	}

	var frame CallFrame
	require.NoError(t, json.Unmarshal(loadGolden(t, "debug_traceCall_callTracer.json"), &frame))

	assert.Equal(t, "CALL", frame.Type)
	assert.Equal(t, "0x0000000000000000000000000000000000000000", frame.From)
	require.NotNil(t, frame.To)
	assert.Equal(t, refUSDTAddressLower, *frame.To)
	assert.NotEmpty(t, frame.Gas)
	assert.Equal(t, "0x644b", frame.GasUsed)
	// USDT totalSupply 결과: uint256 hex 반환
	require.NotNil(t, frame.Output)
	assert.Len(t, *frame.Output, 66, "output should be 0x + 64 hex chars (uint256)")
}

// TestGolden_TraceTransaction_customTracer tests json.RawMessage unmarshaling
// from debug_traceTransaction with a custom JS tracer that strips input/output.
func TestGolden_TraceTransaction_customTracer(t *testing.T) {
	if *update {
		saveGolden(t, "debug_traceTransaction_customTracer.json",
			rpcResult(t, "debug_traceTransaction", refEip1559TxHash, map[string]interface{}{"tracer": customJSTracer}))
	}

	raw := loadGolden(t, "debug_traceTransaction_customTracer.json")
	assert.NotEmpty(t, raw)

	var frame map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &frame))

	// callTracer 구조와 동일하지만 input/output 필드가 없어야 한다
	assert.Equal(t, "CALL", frame["type"])
	assert.Equal(t, refEip1559TxFrom, frame["from"])
	assert.Equal(t, refEip1559TxTo, frame["to"])
	assert.NotEmpty(t, frame["gas"])
	assert.NotEmpty(t, frame["gasUsed"])
	assert.NotContains(t, frame, "input", "custom tracer should not contain input")
	assert.NotContains(t, frame, "output", "custom tracer should not contain output")
	// 하위 calls 존재
	calls, ok := frame["calls"].([]interface{})
	require.True(t, ok, "calls should be an array")
	assert.NotEmpty(t, calls)
	// 하위 call에도 input/output 없어야 함
	firstCall := calls[0].(map[string]interface{})
	assert.NotContains(t, firstCall, "input")
	assert.NotContains(t, firstCall, "output")
}

// TestGolden_TraceCall_customTracer tests json.RawMessage unmarshaling
// from debug_traceCall with a custom JS tracer that strips input/output.
func TestGolden_TraceCall_customTracer(t *testing.T) {
	if *update {
		saveGolden(t, "debug_traceCall_customTracer.json",
			rpcResult(t, "debug_traceCall",
				map[string]interface{}{
					"from": "0x0000000000000000000000000000000000000000",
					"to":   refUSDTAddress,
					"data": "0x18160ddd",
				},
				refBlockHex,
				map[string]interface{}{"tracer": customJSTracer},
			))
	}

	raw := loadGolden(t, "debug_traceCall_customTracer.json")
	assert.NotEmpty(t, raw)

	var frame map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &frame))

	assert.Equal(t, "CALL", frame["type"])
	assert.Equal(t, "0x0000000000000000000000000000000000000000", frame["from"])
	assert.Equal(t, refUSDTAddressLower, frame["to"])
	assert.NotEmpty(t, frame["gas"])
	assert.NotEmpty(t, frame["gasUsed"])
	assert.NotContains(t, frame, "input", "custom tracer should not contain input")
	assert.NotContains(t, frame, "output", "custom tracer should not contain output")
}

// TestGolden_GetRawHeader tests raw hex string from debug_getRawHeader.
func TestGolden_GetRawHeader(t *testing.T) {
	if *update {
		saveGolden(t, "debug_getRawHeader.json",
			rpcResult(t, "debug_getRawHeader", refBlockHex))
	}

	var raw string
	require.NoError(t, json.Unmarshal(loadGolden(t, "debug_getRawHeader.json"), &raw))

	assert.NotEmpty(t, raw)
	assert.Equal(t, "0x", raw[:2], "raw header should be hex-encoded")
	assert.Greater(t, len(raw), 100, "raw header should have meaningful length")
}

// TestGolden_GetRawTransaction tests raw hex string from debug_getRawTransaction.
func TestGolden_GetRawTransaction(t *testing.T) {
	if *update {
		saveGolden(t, "debug_getRawTransaction.json",
			rpcResult(t, "debug_getRawTransaction", refEip1559TxHash))
	}

	var raw string
	require.NoError(t, json.Unmarshal(loadGolden(t, "debug_getRawTransaction.json"), &raw))

	assert.NotEmpty(t, raw)
	assert.Equal(t, "0x", raw[:2], "raw tx should be hex-encoded")
	// EIP-1559 (type 2) RLP은 0x02로 시작
	assert.Equal(t, "0x02", raw[:4], "EIP-1559 tx RLP should start with 0x02")
}

// TestGolden_AccountProof tests AccountProof unmarshaling from eth_getProof.
func TestGolden_AccountProof(t *testing.T) {
	if *update {
		saveGolden(t, "eth_getProof.json",
			rpcResult(t, "eth_getProof", refUSDTAddress,
				[]string{"0x0000000000000000000000000000000000000000000000000000000000000001"},
				"latest"))
	}

	var proof AccountProof
	require.NoError(t, json.Unmarshal(loadGolden(t, "eth_getProof.json"), &proof))

	// 주소 (소문자)
	assert.Equal(t, refUSDTAddressLower, proof.Address)
	// account proof Merkle path
	assert.NotEmpty(t, proof.AccountProof)
	// balance, codeHash, nonce, storageHash는 hex 문자열
	assert.NotEmpty(t, proof.Balance)
	assert.NotEmpty(t, proof.CodeHash)
	assert.NotEmpty(t, proof.Nonce)
	assert.NotEmpty(t, proof.StorageHash)
	// storage proof: slot 1 조회
	require.Len(t, proof.StorageProof, 1)
	assert.Equal(t,
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		proof.StorageProof[0].Key)
	assert.NotEmpty(t, proof.StorageProof[0].Value)
	assert.NotEmpty(t, proof.StorageProof[0].Proof)
}
