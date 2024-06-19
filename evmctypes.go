package evmc

import "github.com/shopspring/decimal"

type Balance struct {
	Address string          `json:"address"`
	Value   decimal.Decimal `json:"value"`
}

type TokenOwner struct {
	Token string `json:"token"`
	Owner string `json:"owner"`
}
