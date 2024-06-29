package evmc

import "github.com/shopspring/decimal"

type Balance struct {
	Address string          `json:"address"`
	Value   decimal.Decimal `json:"value"`
}

type ContractCallParams struct {
	To       string `json:"to"`
	Data     string `json:"data"`
	NumOrTag interface{}
}

type ContractCallResp struct {
	To     string
	Data   string
	Result string
	Error  error
}
