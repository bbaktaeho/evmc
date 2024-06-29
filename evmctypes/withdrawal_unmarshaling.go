package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (w *Withdrawal) UnmarshalJSON(input []byte) error {
	type Withdrawal struct {
		Index          *string `json:"index" validate:"-"`
		ValidatorIndex *string `json:"validatorIndex" validate:"-"`
		Address        *string `json:"address" validate:"required"`
		Amount         *string `json:"amount" validate:"required"`
	}
	var dec Withdrawal
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Index != nil {
		index, err := hexutil.DecodeUint64(*dec.Index)
		if err != nil {
			return err
		}
		w.Index = index
	}
	if dec.ValidatorIndex != nil {
		validatorIndex, err := hexutil.DecodeUint64(*dec.ValidatorIndex)
		if err != nil {
			return err
		}
		w.ValidatorIndex = validatorIndex
	}
	if dec.Address != nil {
		w.Address = *dec.Address
	}
	if dec.Amount != nil {
		amount, err := hexutil.DecodeUint64(*dec.Amount)
		if err != nil {
			return err
		}
		w.Amount = amount
	}
	return nil
}
