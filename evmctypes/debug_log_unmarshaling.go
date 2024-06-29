package evmctypes

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (c *CallLog) UnmarshalJSON(input []byte) error {
	type CallLog struct {
		Address  *string  `json:"address" validate:"required"`
		Topics   []string `json:"topics" validate:"required"`
		Data     *string  `json:"data" validate:"required"`
		Position *string  `json:"position" validate:"-"`
	}
	var dec CallLog
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Address != nil {
		c.Address = *dec.Address
	}
	if dec.Topics != nil {
		c.Topics = dec.Topics
	}
	if dec.Data != nil {
		c.Data = *dec.Data
	}
	if dec.Position != nil {
		position, err := hexutil.DecodeUint64(*dec.Position)
		if err != nil {
			return err
		}
		c.Position = position
	}
	return nil
}
