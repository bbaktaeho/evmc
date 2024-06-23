package token

// TODO: Add more function signatures for extentions and safe ERC20, ERC721, and ERC1155

const (
	TokenNameSig   = "0x06fdde03"
	TokenSymbolSig = "0x95d89b41"
)

const (
	ERC20TotalSupplySig  = "0x18160ddd"
	ERC20DecimalsSig     = "0x313ce567"
	ERC20BalanceOfSig    = "0x70a08231"
	ERC20TransferSig     = "0xa9059cbb"
	ERC20TransferFromSig = "0x23b872dd"
	ERC20ApproveSig      = "0x095ea7b3"
	ERC20AllowanceSig    = "0xdd62ed3e"
)

const (
	ERC721BalanceOfSig         = "0x70a08231"
	ERC721OwnerOfSig           = "0x6352211e"
	ERC721TransferSig          = "0xa9059cbb"
	ERC721TransferFromSig      = "0x23b872dd"
	ERC721ApproveSig           = "0x095ea7b3"
	ERC721SetApprovalForAllSig = "0xa22cb465"
	ERC721GetApprovedSig       = "0x081812fc"
	ERC721IsApprovedForAllSig  = "0xe985e9c5"
)

const (
	ERC1155BalanceOfSig             = "0x00fdd58e"
	ERC1155BalanceOfBatchSig        = "0x4e1273f4"
	ERC1155SetApprovalForAllSig     = "0xa22cb465"
	ERC1155IsApprovedForAllSig      = "0xe985e9c5"
	ERC1155SafeTransferFromSig      = "0xf242432a"
	ERC1155SafeBatchTransferFromSig = "0xbc197c81"
)
