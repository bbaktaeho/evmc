package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

// UnmarshalJSON unmarshals PrestateAccount from JSON.
func (pa *PrestateAccount) UnmarshalJSON(input []byte) error {
	type PrestateAccount0 struct {
		Balance  *string           `json:"balance,omitempty"`
		Code     *string           `json:"code,omitempty"`
		CodeHash *string           `json:"codeHash,omitempty"`
		Nonce    *uint64           `json:"nonce,omitempty"`
		Storage  map[string]string `json:"storage,omitempty"`
	}
	var dec PrestateAccount0
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}

	// Decode balance from hex string to decimal
	if dec.Balance != nil {
		balanceBig, err := hexutil.DecodeBig(*dec.Balance)
		if err != nil {
			return err
		}
		balance := decimal.NewFromBigInt(balanceBig, 0)
		pa.Balance = &balance
	}

	// Copy code
	if dec.Code != nil {
		pa.Code = *dec.Code
	}

	// Copy codeHash
	if dec.CodeHash != nil {
		pa.CodeHash = *dec.CodeHash
	}

	// Decode nonce from hex string to uint64
	if dec.Nonce != nil {
		pa.Nonce = *dec.Nonce
	}

	// Copy storage map
	if dec.Storage != nil {
		pa.Storage = dec.Storage
	}

	return nil
}
