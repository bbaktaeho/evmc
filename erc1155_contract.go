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
	erc1155BalanceOfSig             = "0x00fdd58e"
	erc1155BalanceOfBatchSig        = "0x4e1273f4"
	erc1155SetApprovalForAllSig     = "0xa22cb465"
	erc1155IsApprovedForAllSig      = "0xe985e9c5"
	erc1155SafeTransferFromSig      = "0xf242432a"
	erc1155SafeBatchTransferFromSig = "0x2eb2c2d6"
	erc1155URISig                   = "0x0e89341c"

	erc1155FuncSigSafeTransferFrom  = "safeTransferFrom(address,address,uint256,uint256,bytes)"
	erc1155FuncSigSetApprovalForAll = "setApprovalForAll(address,bool)"
)

type erc1155Contract struct {
	info clientInfo
	c    caller
	ts   transactionSender
}

// --- Generate helpers ---

func GenerateERC1155BalanceOf(owner string, id decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		"balanceOf(address,uint256)",
		evmcsoltypes.Address(owner),
		evmcsoltypes.Uint256(id),
	)
	return input
}

func GenerateERC1155IsApprovedForAll(owner, operator string) string {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	if operator[:2] == "0x" {
		operator = operator[2:]
	}
	return fmt.Sprintf("%s%064s%064s", erc1155IsApprovedForAllSig, owner, operator)
}

func GenerateERC1155URI(id decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		"uri(uint256)",
		evmcsoltypes.Uint256(id),
	)
	return input
}

func GenerateERC1155SetApprovalForAll(operator string, approved bool) string {
	input, _ := evmcutils.GenerateTxInput(
		erc1155FuncSigSetApprovalForAll,
		evmcsoltypes.Address(operator),
		evmcsoltypes.Bool(approved),
	)
	return input
}

func GenerateERC1155SafeTransferFrom(from, to string, id, amount decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		erc1155FuncSigSafeTransferFrom,
		evmcsoltypes.Address(from),
		evmcsoltypes.Address(to),
		evmcsoltypes.Uint256(id),
		evmcsoltypes.Uint256(amount),
	)
	return input
}

// --- View methods ---

func (e *erc1155Contract) BalanceOf(
	tokenAddress string,
	owner string,
	id decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.balanceOf(context.Background(), tokenAddress, owner, id, blockAndTag)
}

func (e *erc1155Contract) BalanceOfWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	id decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.balanceOf(ctx, tokenAddress, owner, id, blockAndTag)
}

func (e *erc1155Contract) balanceOf(
	ctx context.Context,
	tokenAddress string,
	owner string,
	id decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC1155BalanceOf(owner, id)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, fmt.Errorf("BalanceOf: %w", err)
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}

func (e *erc1155Contract) BalanceOfBatch(
	tokenAddress string,
	owners []string,
	ids []decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) ([]decimal.Decimal, error) {
	return e.balanceOfBatch(context.Background(), tokenAddress, owners, ids, blockAndTag)
}

func (e *erc1155Contract) BalanceOfBatchWithContext(
	ctx context.Context,
	tokenAddress string,
	owners []string,
	ids []decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) ([]decimal.Decimal, error) {
	return e.balanceOfBatch(ctx, tokenAddress, owners, ids, blockAndTag)
}

func (e *erc1155Contract) balanceOfBatch(
	ctx context.Context,
	tokenAddress string,
	owners []string,
	ids []decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) ([]decimal.Decimal, error) {
	packed, err := evmcsoltypes.PackBalanceOfBatch(owners, ids)
	if err != nil {
		return nil, fmt.Errorf("BalanceOfBatch: %w", err)
	}
	data := fmt.Sprintf("%s%x", erc1155BalanceOfBatchSig, packed)
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: data},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return nil, fmt.Errorf("BalanceOfBatch: %w", err)
	}
	return evmcsoltypes.ParseSolUintArray(*result)
}

func (e *erc1155Contract) IsApprovedForAll(
	tokenAddress string,
	owner string,
	operator string,
	blockAndTag evmctypes.BlockAndTag,
) (bool, error) {
	return e.isApprovedForAll(context.Background(), tokenAddress, owner, operator, blockAndTag)
}

func (e *erc1155Contract) IsApprovedForAllWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	operator string,
	blockAndTag evmctypes.BlockAndTag,
) (bool, error) {
	return e.isApprovedForAll(ctx, tokenAddress, owner, operator, blockAndTag)
}

func (e *erc1155Contract) isApprovedForAll(
	ctx context.Context,
	tokenAddress string,
	owner string,
	operator string,
	blockAndTag evmctypes.BlockAndTag,
) (bool, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC1155IsApprovedForAll(owner, operator)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return false, fmt.Errorf("IsApprovedForAll: %w", err)
	}
	return evmcsoltypes.ParseBool(*result)
}

func (e *erc1155Contract) URI(
	tokenAddress string,
	id decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.uri(context.Background(), tokenAddress, id, blockAndTag)
}

func (e *erc1155Contract) URIWithContext(
	ctx context.Context,
	tokenAddress string,
	id decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.uri(ctx, tokenAddress, id, blockAndTag)
}

func (e *erc1155Contract) uri(
	ctx context.Context,
	tokenAddress string,
	id decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC1155URI(id)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return "", fmt.Errorf("URI: %w", err)
	}
	return evmcsoltypes.ParseSolStringToString(*result)
}

// --- Write methods ---

func (e *erc1155Contract) SafeTransferFrom(
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	id decimal.Decimal,
	amount decimal.Decimal,
) (string, error) {
	return e.safeTransferFrom(context.Background(), tx, wallet, from, to, id, amount)
}

func (e *erc1155Contract) SafeTransferFromWithContext(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	id decimal.Decimal,
	amount decimal.Decimal,
) (string, error) {
	return e.safeTransferFrom(ctx, tx, wallet, from, to, id, amount)
}

func (e *erc1155Contract) safeTransferFrom(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	id decimal.Decimal,
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
	tx.Data = GenerateERC1155SafeTransferFrom(from, to, id, amount)
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

func (e *erc1155Contract) SetApprovalForAll(
	tx *Tx,
	wallet *Wallet,
	operator string,
	approved bool,
) (string, error) {
	return e.setApprovalForAll(context.Background(), tx, wallet, operator, approved)
}

func (e *erc1155Contract) SetApprovalForAllWithContext(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	operator string,
	approved bool,
) (string, error) {
	return e.setApprovalForAll(ctx, tx, wallet, operator, approved)
}

func (e *erc1155Contract) setApprovalForAll(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	operator string,
	approved bool,
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
	tx.Data = GenerateERC1155SetApprovalForAll(operator, approved)
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
