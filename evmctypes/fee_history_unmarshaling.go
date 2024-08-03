package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type feeHistory struct {
	OldestBlock       uint64     `json:"oldestBlock" validate:"required"`
	BaseFeePerGas     []string   `json:"baseFeePerGas,omitempty" validate:"-"`
	BaseFeePerBlobGas []string   `json:"baseFeePerBlobGas,omitempty" validate:"-"`
	GasUsedRatio      []float64  `json:"gasUsedRatio,omitempty" validate:"-"`
	BlobGasUsedRatio  []float64  `json:"blobGasUsedRatio,omitempty" validate:"-"`
	Reward            [][]string `json:"reward,omitempty" validate:"-"`
}

type _feeHistory struct {
	OldestBlock       string     `json:"oldestBlock"`
	BaseFeePerGas     []string   `json:"baseFeePerGas,omitempty"`
	BaseFeePerBlobGas []string   `json:"baseFeePerBlobGas,omitempty"`
	GasUsedRatio      []float64  `json:"gasUsedRatio,omitempty"`
	BlobGasUsedRatio  []float64  `json:"blobGasUsedRatio,omitempty"`
	Reward            [][]string `json:"reward,omitempty"`
}

func (f *feeHistory) UnmarshalJSON(input []byte) error {
	var _f _feeHistory
	if err := json.Unmarshal(input, &_f); err != nil {
		return err
	}
	f.OldestBlock = hexutil.MustDecodeUint64(_f.OldestBlock)
	f.BaseFeePerGas = _f.BaseFeePerGas
	f.BaseFeePerBlobGas = _f.BaseFeePerBlobGas
	f.GasUsedRatio = _f.GasUsedRatio
	f.BlobGasUsedRatio = _f.BlobGasUsedRatio
	f.Reward = _f.Reward
	return nil
}
