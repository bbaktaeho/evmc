package evmc

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

const (
	erc20NameSig         = "0x06fdde03"
	erc20SymbolSig       = "0x95d89b41"
	erc20TotalSupplySig  = "0x18160ddd"
	erc20DecimalsSig     = "0x313ce567"
	erc20BalanceOfSig    = "0x70a08231"
	erc20TransferSig     = "0xa9059cbb"
	erc20TransferFromSig = "0x23b872dd"
	erc20ApproveSig      = "0x095ea7b3"
	erc20AllowanceSig    = "0xdd62ed3e"
)

type erc20Contract struct {
	c contractCaller
}

func (e *erc20Contract) Name(tokenAddress string, blockAndTag BlockAndTag) (string, error) {
	return e.name(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc20Contract) NameWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag BlockAndTag,
) (string, error) {
	return e.name(ctx, tokenAddress, blockAndTag)
}

func (e *erc20Contract) name(
	ctx context.Context,
	tokenAddress string,
	blockAndTag BlockAndTag,
) (string, error) {
	result := new(string)
	parseBT := parseBlockAndTag(blockAndTag)
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		erc20NameSig,
		parseBT,
	); err != nil {
		return "", err
	}
	name, err := parseSolStringToString(*result)
	if err != nil {
		return parseSolFixedBytesToString(*result)
	}
	return name, nil
}

func (e *erc20Contract) Symbol(tokenAddress string, blockAndTag BlockAndTag) (string, error) {
	return e.symbol(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc20Contract) SymbolWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag BlockAndTag,
) (string, error) {
	return e.symbol(ctx, tokenAddress, blockAndTag)
}

func (e *erc20Contract) symbol(
	ctx context.Context,
	tokenAddress string,
	blockAndTag BlockAndTag,
) (string, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		erc20SymbolSig,
		parsedBT,
	); err != nil {
		return "", err
	}
	symbol, err := parseSolStringToString(*result)
	if err != nil {
		return parseSolFixedBytesToString(*result)
	}
	return symbol, nil
}

func (e *erc20Contract) TotalSupply(tokenAddress string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	return e.totalSupply(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc20Contract) TotalSupplyWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	return e.totalSupply(ctx, tokenAddress, blockAndTag)
}

func (e *erc20Contract) totalSupply(
	ctx context.Context,
	tokenAddress string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		erc20TotalSupplySig,
		parsedBT,
	); err != nil {
		return decimal.Zero, err
	}
	return parseSolUintToDecimal(*result)
}

func (e *erc20Contract) Decimals(tokenAddress string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	return e.decimals(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc20Contract) DecimalsWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	return e.decimals(ctx, tokenAddress, blockAndTag)
}

func (e *erc20Contract) decimals(
	ctx context.Context,
	tokenAddress string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		erc20DecimalsSig,
		parsedBT,
	); err != nil {
		return decimal.Zero, err
	}
	return parseSolUintToDecimal(*result)
}

func (e *erc20Contract) Approve() {}

func (e *erc20Contract) Transfer() {}

func (e *erc20Contract) TransferFrom() {}

func (e *erc20Contract) BalanceOf(tokenAddress string, owner string, blockAndTag BlockAndTag) (decimal.Decimal, error) {
	return e.balanceOf(context.Background(), tokenAddress, owner, blockAndTag)
}

func (e *erc20Contract) BalanceOfWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	return e.balanceOf(ctx, tokenAddress, owner, blockAndTag)
}

func (e *erc20Contract) balanceOf(
	ctx context.Context,
	tokenAddress string,
	owner string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		fmt.Sprintf("%s%064s", erc20BalanceOfSig, owner),
		parsedBT,
	); err != nil {
		return decimal.Zero, err
	}
	return parseSolUintToDecimal(*result)
}

func (e *erc20Contract) Allowance(
	tokenAddress string,
	owner string,
	spender string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	return e.allowance(context.Background(), tokenAddress, owner, spender, blockAndTag)
}

func (e *erc20Contract) AllowanceWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	spender string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	return e.allowance(ctx, tokenAddress, owner, spender, blockAndTag)
}

func (e *erc20Contract) allowance(
	ctx context.Context,
	tokenAddress string,
	owner string,
	spender string,
	blockAndTag BlockAndTag,
) (decimal.Decimal, error) {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	if spender[:2] == "0x" {
		spender = spender[2:]
	}
	result := new(string)
	parsedBT := parseBlockAndTag(blockAndTag)
	if err := e.c.contractCall(
		ctx,
		result,
		tokenAddress,
		fmt.Sprintf("%s%064s%064s", erc20AllowanceSig, owner, spender),
		parsedBT,
	); err != nil {
		return decimal.Zero, err
	}
	return parseSolUintToDecimal(*result)
}
