package evmc

import (
	"context"

	"github.com/bbaktaeho/evmc/token"
	"github.com/shopspring/decimal"
)

type erc20Contract struct {
	c contractCaller
}

func (e *erc20Contract) Name(tokenAddress string, numToTag interface{}) (string, error) {
	return e.name(context.Background(), tokenAddress, numToTag)
}

func (e *erc20Contract) NameWithContext(
	ctx context.Context,
	tokenAddress string,
	numToTag interface{},
) (string, error) {
	return e.name(ctx, tokenAddress, numToTag)
}

func (e *erc20Contract) name(
	ctx context.Context,
	tokenAddress string,
	numToTag interface{},
) (string, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numToTag)
	if err != nil {
		return "", err
	}
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		token.TokenNameSig,
		parsedNumOrTag,
	); err != nil {
		return "", err
	}
	name, err := parseSolStringToStr(*result)
	if err != nil {
		return parseSolBytes32ToStr(*result)
	}
	return name, nil
}

func (e *erc20Contract) Symbol(tokenAddress string, numToTag interface{}) (string, error) {
	return e.symbol(context.Background(), tokenAddress, numToTag)
}

func (e *erc20Contract) SymbolWithContext(
	ctx context.Context,
	tokenAddress string,
	numToTag interface{},
) (string, error) {
	return e.symbol(ctx, tokenAddress, numToTag)
}

func (e *erc20Contract) symbol(
	ctx context.Context,
	tokenAddress string,
	numToTag interface{},
) (string, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numToTag)
	if err != nil {
		return "", err
	}
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		token.TokenSymbolSig,
		parsedNumOrTag,
	); err != nil {
		return "", err
	}
	symbol, err := parseSolStringToStr(*result)
	if err != nil {
		return parseSolBytes32ToStr(*result)
	}
	return symbol, nil
}

func (e *erc20Contract) TotalSupply(tokenAddress string, numToTag interface{}) (decimal.Decimal, error) {
	return e.totalSupply(context.Background(), tokenAddress, numToTag)
}

func (e *erc20Contract) TotalSupplyWithContext(
	ctx context.Context,
	tokenAddress string,
	numToTag interface{},
) (decimal.Decimal, error) {
	return e.totalSupply(ctx, tokenAddress, numToTag)
}

func (e *erc20Contract) totalSupply(
	ctx context.Context,
	tokenAddress string,
	numToTag interface{},
) (decimal.Decimal, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numToTag)
	if err != nil {
		return decimal.Zero, err
	}
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		token.ERC20TotalSupplySig,
		parsedNumOrTag,
	); err != nil {
		return decimal.Zero, err
	}
	return parseSolUint256ToDecimal(*result)
}

func (e *erc20Contract) Decimals() {}

func (e *erc20Contract) Approve() {}

func (e *erc20Contract) Transfer() {}

func (e *erc20Contract) TransferFrom() {}

func (e *erc20Contract) BalanceOf() {}

func (e *erc20Contract) Allowance() {}
