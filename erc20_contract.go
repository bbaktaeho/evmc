package evmc

import "context"

type erc20Contract struct {
	c contractCaller
}

func (e *erc20Contract) Name(
	ctx context.Context,
	tokenAddress string,
	numToTag interface{},
) (string, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numToTag)
	if err != nil {
		return "", err
	}
	if err := e.c.contractCall(ctx, result, tokenAddress, tokenNameSig, parsedNumOrTag); err != nil {
		return "", err
	}
	return *result, nil
}

func (e *erc20Contract) Symbol() {}

func (e *erc20Contract) TotalSupply() {}

func (e *erc20Contract) Decimals() {}

func (e *erc20Contract) Approve() {}

func (e *erc20Contract) Transfer() {}

func (e *erc20Contract) TransferFrom() {}

func (e *erc20Contract) BalanceOf() {}

func (e *erc20Contract) Allowance() {}
