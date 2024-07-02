package evmcutils

import (
	"math/big"
	"testing"

	"github.com/bbaktaeho/evmc"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTxInput(t *testing.T) {
	input, err := GenerateTxInput("transfer(address,uint256)", evmc.ParseAddress("0x0a3f6849f78076aefaDf113F5BED87720274dDC0"), evmc.ParseUint256(decimal.NewFromBigInt(big.NewInt(1), 18)))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "0xa9059cbb0000000000000000000000000a3f6849f78076aefadf113f5bed87720274ddc00000000000000000000000000000000000000000000000000de0b6b3a7640000", input)
}
