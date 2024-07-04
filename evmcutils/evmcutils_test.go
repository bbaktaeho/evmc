package evmcutils

import (
	"math/big"
	"testing"

	"github.com/bbaktaeho/evmc/evmcsoltypes"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTxInput(t *testing.T) {
	type args struct {
		funcSig string
		args    []evmcsoltypes.SolType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "erc20 transfer",
			args: args{
				funcSig: "transfer(address,uint256)",
				args: []evmcsoltypes.SolType{
					evmcsoltypes.Address("0x0a3f6849f78076aefaDf113F5BED87720274dDC0"),
					evmcsoltypes.Uint256(decimal.NewFromBigInt(big.NewInt(1), 18)),
				},
			},
			want: "0xa9059cbb0000000000000000000000000a3f6849f78076aefadf113f5bed87720274ddc00000000000000000000000000000000000000000000000000de0b6b3a7640000",
		},
		{
			name: "erc20 transferFrom",
			args: args{
				funcSig: "transferFrom(address,address,uint256)",
				args: []evmcsoltypes.SolType{
					evmcsoltypes.Address("0x0a3f6849f78076aefaDf113F5BED87720274dDC0"),
					evmcsoltypes.Address("0x0a3f6849f78076aefaDf113F5BED87720274dDC0"),
					evmcsoltypes.Uint256(decimal.NewFromBigInt(big.NewInt(1), 18)),
				},
			},
			want: "0x23b872dd0000000000000000000000000a3f6849f78076aefadf113f5bed87720274ddc00000000000000000000000000a3f6849f78076aefadf113f5bed87720274ddc00000000000000000000000000000000000000000000000000de0b6b3a7640000",
		},
		{
			name: "erc20 approve",
			args: args{
				funcSig: "approve(address,uint256)",
				args: []evmcsoltypes.SolType{
					evmcsoltypes.Address("0x0a3f6849f78076aefaDf113F5BED87720274dDC0"),
					evmcsoltypes.Uint256(decimal.NewFromBigInt(big.NewInt(1), 18)),
				},
			},
			want: "0x095ea7b30000000000000000000000000a3f6849f78076aefadf113f5bed87720274ddc00000000000000000000000000000000000000000000000000de0b6b3a7640000",
		},
		{
			name: "bulk transfer",
			args: args{
				funcSig: "transfers(address[])",
				args: []evmcsoltypes.SolType{
					evmcsoltypes.AddressArr([]string{"0x0a3f6849f78076aefaDf113F5BED87720274dDC0", "0x0a3f6849f78076aefaDf113F5BED87720274dDC0"}),
				},
			},
			want: "0x91ca8b56000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000a3f6849f78076aefadf113f5bed87720274ddc00000000000000000000000000a3f6849f78076aefadf113f5bed87720274ddc0",
		},
		{
			name: "error",
			args: args{
				funcSig: "invalid",
				args:    []evmcsoltypes.SolType{},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateTxInput(tt.args.funcSig, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTxInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateLogTopic(t *testing.T) {
	type args struct {
		eventSig string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "erc20 transfer",
			args: args{
				eventSig: "Transfer(address,address,uint256)",
			},
			want: "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		},
		{
			name: "erc20 approval",
			args: args{
				eventSig: "Approval(address,address,uint256)",
			},
			want: "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateLogTopic(tt.args.eventSig)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateLogTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateLogTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
