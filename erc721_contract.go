package evmc

const (
	erc721NameSig              = "0x06fdde03"
	erc721SymbolSig            = "0x95d89b41"
	erc721BalanceOfSig         = "0x70a08231"
	erc721OwnerOfSig           = "0x6352211e"
	erc721TransferSig          = "0xa9059cbb"
	erc721TransferFromSig      = "0x23b872dd"
	erc721ApproveSig           = "0x095ea7b3"
	erc721SetApprovalForAllSig = "0xa22cb465"
	erc721GetApprovedSig       = "0x081812fc"
	erc721IsApprovedForAllSig  = "0xe985e9c5"
)

type erc721Contract struct {
	c contractCaller
}
