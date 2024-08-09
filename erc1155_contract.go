package evmc

const (
	erc1155NameSig                 = "0x06fdde03"
	erc1155SymbolSig               = "0x95d89b41"
	erc155BalanceOfSig             = "0x00fdd58e"
	erc155BalanceOfBatchSig        = "0x4e1273f4"
	erc155SetApprovalForAllSig     = "0xa22cb465"
	erc155IsApprovedForAllSig      = "0xe985e9c5"
	erc155SafeTransferFromSig      = "0xf242432a"
	erc155SafeBatchTransferFromSig = "0xbc197c81"
)

type erc1155Contract struct {
	c caller
}
