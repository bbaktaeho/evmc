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
	erc721NameSig              = "0x06fdde03"
	erc721SymbolSig            = "0x95d89b41"
	erc721BalanceOfSig         = "0x70a08231"
	erc721OwnerOfSig           = "0x6352211e"
	erc721TransferFromSig      = "0x23b872dd"
	erc721SafeTransferFromSig  = "0x42842e0e"
	erc721ApproveSig           = "0x095ea7b3"
	erc721SetApprovalForAllSig = "0xa22cb465"
	erc721GetApprovedSig       = "0x081812fc"
	erc721IsApprovedForAllSig  = "0xe985e9c5"
	erc721TokenURISig          = "0xc87b56dd"
	erc721TotalSupplySig       = "0x18160ddd"
	erc721TokenByIndexSig      = "0x4f6ccce7"
	erc721TokenOfOwnerByIndexSig = "0x2f745c59"

	erc721FuncSigTransferFrom      = "transferFrom(address,address,uint256)"
	erc721FuncSigSafeTransferFrom  = "safeTransferFrom(address,address,uint256)"
	erc721FuncSigApprove           = "approve(address,uint256)"
	erc721FuncSigSetApprovalForAll = "setApprovalForAll(address,bool)"
)

type erc721Contract struct {
	info clientInfo
	c    caller
	ts   transactionSender
}

// --- Generate helpers ---

func GenerateERC721BalanceOf(owner string) string {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	return fmt.Sprintf("%s%064s", erc721BalanceOfSig, owner)
}

func GenerateERC721OwnerOf(tokenID decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		"ownerOf(uint256)",
		evmcsoltypes.Uint256(tokenID),
	)
	return input
}

func GenerateERC721TokenURI(tokenID decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		"tokenURI(uint256)",
		evmcsoltypes.Uint256(tokenID),
	)
	return input
}

func GenerateERC721GetApproved(tokenID decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		"getApproved(uint256)",
		evmcsoltypes.Uint256(tokenID),
	)
	return input
}

func GenerateERC721IsApprovedForAll(owner, operator string) string {
	if owner[:2] == "0x" {
		owner = owner[2:]
	}
	if operator[:2] == "0x" {
		operator = operator[2:]
	}
	return fmt.Sprintf("%s%064s%064s", erc721IsApprovedForAllSig, owner, operator)
}

func GenerateERC721TransferFrom(from, to string, tokenID decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		erc721FuncSigTransferFrom,
		evmcsoltypes.Address(from),
		evmcsoltypes.Address(to),
		evmcsoltypes.Uint256(tokenID),
	)
	return input
}

func GenerateERC721SafeTransferFrom(from, to string, tokenID decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		erc721FuncSigSafeTransferFrom,
		evmcsoltypes.Address(from),
		evmcsoltypes.Address(to),
		evmcsoltypes.Uint256(tokenID),
	)
	return input
}

func GenerateERC721Approve(approved string, tokenID decimal.Decimal) string {
	input, _ := evmcutils.GenerateTxInput(
		erc721FuncSigApprove,
		evmcsoltypes.Address(approved),
		evmcsoltypes.Uint256(tokenID),
	)
	return input
}

func GenerateERC721SetApprovalForAll(operator string, approved bool) string {
	input, _ := evmcutils.GenerateTxInput(
		erc721FuncSigSetApprovalForAll,
		evmcsoltypes.Address(operator),
		evmcsoltypes.Bool(approved),
	)
	return input
}

// --- View methods ---

func (e *erc721Contract) Name(tokenAddress string, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.name(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc721Contract) NameWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.name(ctx, tokenAddress, blockAndTag)
}

func (e *erc721Contract) name(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: erc721NameSig},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return "", fmt.Errorf("Name: %w", err)
	}
	name, err := evmcsoltypes.ParseSolStringToString(*result)
	if err != nil {
		return evmcsoltypes.ParseSolFixedBytesToString(*result)
	}
	return name, nil
}

func (e *erc721Contract) Symbol(tokenAddress string, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.symbol(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc721Contract) SymbolWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.symbol(ctx, tokenAddress, blockAndTag)
}

func (e *erc721Contract) symbol(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: erc721SymbolSig},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return "", fmt.Errorf("Symbol: %w", err)
	}
	symbol, err := evmcsoltypes.ParseSolStringToString(*result)
	if err != nil {
		return evmcsoltypes.ParseSolFixedBytesToString(*result)
	}
	return symbol, nil
}

func (e *erc721Contract) BalanceOf(tokenAddress string, owner string, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	return e.balanceOf(context.Background(), tokenAddress, owner, blockAndTag)
}

func (e *erc721Contract) BalanceOfWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.balanceOf(ctx, tokenAddress, owner, blockAndTag)
}

func (e *erc721Contract) balanceOf(
	ctx context.Context,
	tokenAddress string,
	owner string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC721BalanceOf(owner)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, fmt.Errorf("BalanceOf: %w", err)
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}

func (e *erc721Contract) OwnerOf(tokenAddress string, tokenID decimal.Decimal, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.ownerOf(context.Background(), tokenAddress, tokenID, blockAndTag)
}

func (e *erc721Contract) OwnerOfWithContext(
	ctx context.Context,
	tokenAddress string,
	tokenID decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.ownerOf(ctx, tokenAddress, tokenID, blockAndTag)
}

func (e *erc721Contract) ownerOf(
	ctx context.Context,
	tokenAddress string,
	tokenID decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC721OwnerOf(tokenID)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return "", fmt.Errorf("OwnerOf: %w", err)
	}
	return evmcsoltypes.ParseSolAddress(*result)
}

func (e *erc721Contract) TokenURI(tokenAddress string, tokenID decimal.Decimal, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.tokenURI(context.Background(), tokenAddress, tokenID, blockAndTag)
}

func (e *erc721Contract) TokenURIWithContext(
	ctx context.Context,
	tokenAddress string,
	tokenID decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.tokenURI(ctx, tokenAddress, tokenID, blockAndTag)
}

func (e *erc721Contract) tokenURI(
	ctx context.Context,
	tokenAddress string,
	tokenID decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC721TokenURI(tokenID)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return "", fmt.Errorf("TokenURI: %w", err)
	}
	return evmcsoltypes.ParseSolStringToString(*result)
}

func (e *erc721Contract) GetApproved(tokenAddress string, tokenID decimal.Decimal, blockAndTag evmctypes.BlockAndTag) (string, error) {
	return e.getApproved(context.Background(), tokenAddress, tokenID, blockAndTag)
}

func (e *erc721Contract) GetApprovedWithContext(
	ctx context.Context,
	tokenAddress string,
	tokenID decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	return e.getApproved(ctx, tokenAddress, tokenID, blockAndTag)
}

func (e *erc721Contract) getApproved(
	ctx context.Context,
	tokenAddress string,
	tokenID decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (string, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC721GetApproved(tokenID)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return "", fmt.Errorf("GetApproved: %w", err)
	}
	return evmcsoltypes.ParseSolAddress(*result)
}

func (e *erc721Contract) IsApprovedForAll(
	tokenAddress string,
	owner string,
	operator string,
	blockAndTag evmctypes.BlockAndTag,
) (bool, error) {
	return e.isApprovedForAll(context.Background(), tokenAddress, owner, operator, blockAndTag)
}

func (e *erc721Contract) IsApprovedForAllWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	operator string,
	blockAndTag evmctypes.BlockAndTag,
) (bool, error) {
	return e.isApprovedForAll(ctx, tokenAddress, owner, operator, blockAndTag)
}

func (e *erc721Contract) isApprovedForAll(
	ctx context.Context,
	tokenAddress string,
	owner string,
	operator string,
	blockAndTag evmctypes.BlockAndTag,
) (bool, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: GenerateERC721IsApprovedForAll(owner, operator)},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return false, fmt.Errorf("IsApprovedForAll: %w", err)
	}
	return evmcsoltypes.ParseBool(*result)
}

// --- ERC721Enumerable extension ---

func (e *erc721Contract) TotalSupply(tokenAddress string, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	return e.totalSupply(context.Background(), tokenAddress, blockAndTag)
}

func (e *erc721Contract) TotalSupplyWithContext(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.totalSupply(ctx, tokenAddress, blockAndTag)
}

func (e *erc721Contract) totalSupply(
	ctx context.Context,
	tokenAddress string,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: erc721TotalSupplySig},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, fmt.Errorf("TotalSupply: %w", err)
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}

func (e *erc721Contract) TokenByIndex(tokenAddress string, index decimal.Decimal, blockAndTag evmctypes.BlockAndTag) (decimal.Decimal, error) {
	return e.tokenByIndex(context.Background(), tokenAddress, index, blockAndTag)
}

func (e *erc721Contract) TokenByIndexWithContext(
	ctx context.Context,
	tokenAddress string,
	index decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.tokenByIndex(ctx, tokenAddress, index, blockAndTag)
}

func (e *erc721Contract) tokenByIndex(
	ctx context.Context,
	tokenAddress string,
	index decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	input, _ := evmcutils.GenerateTxInput("tokenByIndex(uint256)", evmcsoltypes.Uint256(index))
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: input},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, fmt.Errorf("TokenByIndex: %w", err)
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}

func (e *erc721Contract) TokenOfOwnerByIndex(
	tokenAddress string,
	owner string,
	index decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.tokenOfOwnerByIndex(context.Background(), tokenAddress, owner, index, blockAndTag)
}

func (e *erc721Contract) TokenOfOwnerByIndexWithContext(
	ctx context.Context,
	tokenAddress string,
	owner string,
	index decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	return e.tokenOfOwnerByIndex(ctx, tokenAddress, owner, index, blockAndTag)
}

func (e *erc721Contract) tokenOfOwnerByIndex(
	ctx context.Context,
	tokenAddress string,
	owner string,
	index decimal.Decimal,
	blockAndTag evmctypes.BlockAndTag,
) (decimal.Decimal, error) {
	input, _ := evmcutils.GenerateTxInput(
		"tokenOfOwnerByIndex(address,uint256)",
		evmcsoltypes.Address(owner),
		evmcsoltypes.Uint256(index),
	)
	var (
		result = new(string)
		params = []any{
			&evmctypes.QueryParams{To: tokenAddress, Data: input},
			evmctypes.ParseBlockAndTag(blockAndTag),
		}
	)
	if err := e.c.call(ctx, result, EthCall, params...); err != nil {
		return decimal.Zero, fmt.Errorf("TokenOfOwnerByIndex: %w", err)
	}
	return evmcsoltypes.ParseSolUintToDecimal(*result)
}

// --- Write methods ---

func (e *erc721Contract) TransferFrom(
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	tokenID decimal.Decimal,
) (string, error) {
	return e.transferFrom(context.Background(), tx, wallet, from, to, tokenID)
}

func (e *erc721Contract) TransferFromWithContext(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	tokenID decimal.Decimal,
) (string, error) {
	return e.transferFrom(ctx, tx, wallet, from, to, tokenID)
}

func (e *erc721Contract) transferFrom(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	tokenID decimal.Decimal,
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
	tx.Data = GenerateERC721TransferFrom(from, to, tokenID)
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

func (e *erc721Contract) SafeTransferFrom(
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	tokenID decimal.Decimal,
) (string, error) {
	return e.safeTransferFrom(context.Background(), tx, wallet, from, to, tokenID)
}

func (e *erc721Contract) SafeTransferFromWithContext(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	tokenID decimal.Decimal,
) (string, error) {
	return e.safeTransferFrom(ctx, tx, wallet, from, to, tokenID)
}

func (e *erc721Contract) safeTransferFrom(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	from string,
	to string,
	tokenID decimal.Decimal,
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
	tx.Data = GenerateERC721SafeTransferFrom(from, to, tokenID)
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

func (e *erc721Contract) Approve(
	tx *Tx,
	wallet *Wallet,
	approved string,
	tokenID decimal.Decimal,
) (string, error) {
	return e.approve(context.Background(), tx, wallet, approved, tokenID)
}

func (e *erc721Contract) ApproveWithContext(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	approved string,
	tokenID decimal.Decimal,
) (string, error) {
	return e.approve(ctx, tx, wallet, approved, tokenID)
}

func (e *erc721Contract) approve(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	approved string,
	tokenID decimal.Decimal,
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
	tx.Data = GenerateERC721Approve(approved, tokenID)
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

func (e *erc721Contract) SetApprovalForAll(
	tx *Tx,
	wallet *Wallet,
	operator string,
	approved bool,
) (string, error) {
	return e.setApprovalForAll(context.Background(), tx, wallet, operator, approved)
}

func (e *erc721Contract) SetApprovalForAllWithContext(
	ctx context.Context,
	tx *Tx,
	wallet *Wallet,
	operator string,
	approved bool,
) (string, error) {
	return e.setApprovalForAll(ctx, tx, wallet, operator, approved)
}

func (e *erc721Contract) setApprovalForAll(
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
	tx.Data = GenerateERC721SetApprovalForAll(operator, approved)
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
