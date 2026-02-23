package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

// UnmarshalJSON implements json.Unmarshaler for PrestateAccount.
// Balance is decoded from a hex string to decimal.Decimal.
// Nonce is an integer (not hex) as returned by geth.
func (pa *PrestateAccount) UnmarshalJSON(input []byte) error {
	type wire struct {
		Balance  *string           `json:"balance,omitempty"`
		Code     *string           `json:"code,omitempty"`
		CodeHash *string           `json:"codeHash,omitempty"`
		Nonce    *uint64           `json:"nonce,omitempty"`
		Storage  map[string]string `json:"storage,omitempty"`
	}
	var dec wire
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Balance != nil {
		balanceBig, err := hexutil.DecodeBig(*dec.Balance)
		if err != nil {
			return err
		}
		balance := decimal.NewFromBigInt(balanceBig, 0)
		pa.Balance = &balance
	}
	if dec.Code != nil {
		pa.Code = *dec.Code
	}
	if dec.CodeHash != nil {
		pa.CodeHash = *dec.CodeHash
	}
	if dec.Nonce != nil {
		pa.Nonce = *dec.Nonce
	}
	if dec.Storage != nil {
		pa.Storage = dec.Storage
	}
	return nil
}
