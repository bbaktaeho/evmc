package evmc

import (
	"context"
	"fmt"

	"github.com/bbaktaeho/evmc/evmcsoltypes"
	"github.com/bbaktaeho/evmc/evmctypes"
	"github.com/bbaktaeho/evmc/evmcutils"
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

	erc20FuncSigTransfer     = "transfer(address,uint256)"
	erc20FuncSigTransferFrom = "transferFrom(address,address,uint256)"
	erc20FuncSigApprove      = "approve(address,uint256)"
)

type erc20Contract struct {
	info clientInfo
	c    caller
	ts   transactionSender
}

func GenerateERC20BalanceOf(owner string) string {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	return fmt.Sprintf("%s%064s", erc20BalanceOfSig, owner)
}

func GenerateERC20Allowance(owner, spender string) string {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	if spender[:2] == "0x" {
		spender = spender[2:]
	}
	return fmt.Sprintf("%s%064s%064s", erc20AllowanceSig, owner, spender)
}

func GenerateERC20Transfer(recipient string, amount decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		erc20FuncSigTransfer,
		evmcsoltypes.Address(recipient),
		evmcsoltypes.Uint256(amount),
	)
	return input
}

func (e *erc20Contract) Name(tokenAddress string, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.name(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc20Contract) NameWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.name(ctx, tokenAddress, blockAndTag)
}

func (e *erc20Contract) name(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	var (
		result = new(string)
		params = []interface{}{
			&evmctypes.QueryParams{To: tokenAddress, Data: erc20NameSig},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return "", err
	}
	name, err := evmcsoltypes.ParseSolStringToString(*result)
	if err != nil {
		return evmcsoltypes.ParseSolFixedBytesToString(*result)
	}
	return name, nil
}

func (e *erc20Contract) Symbol(tokenAddress string, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.symbol(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc20Contract) SymbolWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.symbol(ctx, tokenAddress, blockAndTag)
}

func (e *erc20Contract) symbol(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	var (
		result = new(string)
		params = []interface{}{
			evmctypes.QueryParams{To: tokenAddress, Data: erc20SymbolSig},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return "", err
	}
	symbol, err := evmcsoltypes.ParseSolStringToString(*result)
	if err != nil {
		return evmcsoltypes.ParseSolFixedBytesToString(*result)
	}
	return symbol, nil
}

func (e *erc20Contract) TotalSupply(tokenAddress string, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	return e.totalSupply(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc20Contract) TotalSupplyWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.totalSupply(ctx, tokenAddress, blockAndTag)
}

func (e *erc20Contract) totalSupply(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	var (
		result = new(string)
		params = []interface{}{
			evmctypes.QueryParams{To: tokenAddress, Data: erc20TotalSupplySig},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, err
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}

func (e *erc20Contract) Decimals(tokenAddress string, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	return e.decimals(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc20Contract) DecimalsWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.decimals(ctx, tokenAddress, blockAndTag)
}

func (e *erc20Contract) decimals(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	var (
		result = new(string)
		params = []interface{}{
			evmctypes.QueryParams{To: tokenAddress, Data: erc20DecimalsSig},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, err
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}

// TODO: implement
func (e *erc20Contract) Approve(tx *Tx, wallet *Wallet, spender string, amount decimal.Decimal) {}

// TODO: check gas price and base fee
func (e *erc20Contract) Transfer(tx *Tx, wallet *Wallet, recipient string, amount decimal.Decimal) (string, error) {
	return e.transfer(context.Background(), tx, wallet, recipient, amount)
}

func (e *erc20Contract) TransferWithContext(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	recipient string,
	amount decimal.Decimal,
) (string, error) {
	return e.transfer(ctx, tx, wallet, recipient, amount)
}

func (e *erc20Contract) transfer(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	recipient string,
	amount decimal.Decimal,
) (string, error) {
	if tx == nil {
		return "", ErrTxRequired
	}
	if wallet == nil {
		return "", ErrWalletRequired
	}
	if err := tx.valid(); err != nil {
		return "", err
	}
	tx.Data = GenerateERC20Transfer(recipient, amount)
	sendingTx, err := NewSendingTx(tx)
	if err != nil {
		return "", err
	}
	chainID, err := e.info.ChainID()
	if err != nil {
		return "", err
	}
	_, rawTx, err := wallet.SignTx(sendingTx, chainID)
	if err != nil {
		return "", err
	}
	return e.ts.sendRawTransaction(ctx, rawTx)
}

func (e *erc20Contract) TransferFrom() {}

func (e *erc20Contract) BalanceOf(tokenAddress string, owner string, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	return e.balanceOf(context.Background(), tokenAddress, owner, blockAndTag)
}

func (e *erc20Contract) BalanceOfWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.balanceOf(ctx, tokenAddress, owner, blockAndTag)
}

func (e *erc20Contract) balanceOf(
	ctx context.Context,
	tokenAddress string,
	owner string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	var (
		result = new(string)
		params = []interface{}{
			evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC20BalanceOf(owner)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, err
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}

func (e *erc20Contract) Allowance(
	tokenAddress string,
	owner string,
	spender string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.allowance(context.Background(), tokenAddress, owner, spender, blockAndTag)
}

func (e *erc20Contract) AllowanceWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	spender string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.allowance(ctx, tokenAddress, owner, spender, blockAndTag)
}

func (e *erc20Contract) allowance(
	ctx context.Context,
	tokenAddress string,
	owner string,
	spender string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	var (
		result = new(string)
		params = []interface{}{
			evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC20Allowance(owner, spender)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, err
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}
