package evmcsoltypes_test

import (
	"testing"

	"github.com/bbaktaeho/evmc/evmcsoltypes"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		solReturn string
		want      bool
		wantErr   bool
	}{
		{
			name:      "true",
			solReturn: "0x0000000000000000000000000000000000000000000000000000000000000001",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "false",
			solReturn: "0x0000000000000000000000000000000000000000000000000000000000000000",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "invalid hex",
			solReturn: "0xZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
			want:      false,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := evmcsoltypes.ParseBool(tt.solReturn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseBool() failed: %v", gotErr)
				}
				return
			}
			if got != tt.want {
				t.Errorf("ParseBool() = %v, want %v", got, tt.want)
			}
			if tt.wantErr {
				t.Fatal("ParseBool() succeeded unexpectedly")
			}
		})
	}
}
