package evmc

import (
	"context"
	"fmt"

	"github.com/bbaktaeho/evmc/token"
	"github.com/shopspring/decimal"
)

type erc20Contract struct {
	c contractCaller
}

func (e *erc20Contract) Name(tokenAddress string, numOrTag interface{}) (string, error) {
	return e.name(context.Background(), tokenAddress, numOrTag)
}

func (e *erc20Contract) NameWithContext(
	ctx context.Context,
	tokenAddress string,
	numOrTag interface{},
) (string, error) {
	return e.name(ctx, tokenAddress, numOrTag)
}

func (e *erc20Contract) name(
	ctx context.Context,
	tokenAddress string,
	numOrTag interface{},
) (string, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numOrTag)
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
	name, err := parseSolStringToString(*result)
	if err != nil {
		return parseSolFixedBytesToString(*result)
	}
	return name, nil
}

func (e *erc20Contract) Symbol(tokenAddress string, numOrTag interface{}) (string, error) {
	return e.symbol(context.Background(), tokenAddress, numOrTag)
}

func (e *erc20Contract) SymbolWithContext(
	ctx context.Context,
	tokenAddress string,
	numOrTag interface{},
) (string, error) {
	return e.symbol(ctx, tokenAddress, numOrTag)
}

func (e *erc20Contract) symbol(
	ctx context.Context,
	tokenAddress string,
	numOrTag interface{},
) (string, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numOrTag)
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
	symbol, err := parseSolStringToString(*result)
	if err != nil {
		return parseSolFixedBytesToString(*result)
	}
	return symbol, nil
}

func (e *erc20Contract) TotalSupply(tokenAddress string, numOrTag interface{}) (decimal.Decimal, error) {
	return e.totalSupply(context.Background(), tokenAddress, numOrTag)
}

func (e *erc20Contract) TotalSupplyWithContext(
	ctx context.Context,
	tokenAddress string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	return e.totalSupply(ctx, tokenAddress, numOrTag)
}

func (e *erc20Contract) totalSupply(
	ctx context.Context,
	tokenAddress string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numOrTag)
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
	return parseSolUintToDecimal(*result)
}

func (e *erc20Contract) Decimals(tokenAddress string, numOrTag interface{}) (decimal.Decimal, error) {
	return e.decimals(context.Background(), tokenAddress, numOrTag)
}

func (e *erc20Contract) DecimalsWithContext(
	ctx context.Context,
	tokenAddress string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	return e.decimals(ctx, tokenAddress, numOrTag)
}

func (e *erc20Contract) decimals(
	ctx context.Context,
	tokenAddress string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numOrTag)
	if err != nil {
		return decimal.Zero, err
	}
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		token.ERC20DecimalsSig,
		parsedNumOrTag,
	); err != nil {
		return decimal.Zero, err
	}
	return parseSolUintToDecimal(*result)
}

func (e *erc20Contract) Approve() {}

func (e *erc20Contract) Transfer() {}

func (e *erc20Contract) TransferFrom() {}

func (e *erc20Contract) BalanceOf(tokenAddress string, owner string, numOrTag interface{}) (decimal.Decimal, error) {
	return e.balanceOf(context.Background(), tokenAddress, owner, numOrTag)
}

func (e *erc20Contract) BalanceOfWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	return e.balanceOf(ctx, tokenAddress, owner, numOrTag)
}

func (e *erc20Contract) balanceOf(
	ctx context.Context,
	tokenAddress string,
	owner string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numOrTag)
	if err != nil {
		return decimal.Zero, err
	}
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		fmt.Sprintf("%s%064s", token.ERC20BalanceOfSig, owner),
		parsedNumOrTag,
	); err != nil {
		return decimal.Zero, err
	}
	return parseSolUintToDecimal(*result)
}

func (e *erc20Contract) Allowance(
	tokenAddress string,
	owner string,
	spender string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	return e.allowance(context.Background(), tokenAddress, owner, spender, numOrTag)
}

func (e *erc20Contract) AllowanceWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	spender string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	return e.allowance(ctx, tokenAddress, owner, spender, numOrTag)
}

func (e *erc20Contract) allowance(
	ctx context.Context,
	tokenAddress string,
	owner string,
	spender string,
	numOrTag interface{},
) (decimal.Decimal, error) {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	if spender[:2] == "0x" {
		spender = spender[2:]
	}
	result := new(string)
	parsedNumOrTag, err := parseNumOrTag(numOrTag)
	if err != nil {
		return decimal.Zero, err
	}
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		fmt.Sprintf("%s%064s%064s", token.ERC20AllowanceSig, owner, spender),
		parsedNumOrTag,
	); err != nil {
		return decimal.Zero, err
	}
	return parseSolUintToDecimal(*result)
}
