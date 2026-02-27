// Package evmcsoltypes provides ABI encoding and decoding helpers for
// Solidity types used in EVM smart contract interactions.
//
// Encoding functions pack Go values into ABI-encoded byte slices:
//
//	evmcsoltypes.Address("0x...")
//	evmcsoltypes.Uint256(decimal.NewFromInt(1000))
//	evmcsoltypes.Bool(true)
//
// Parsing functions decode hex-encoded return values from eth_call responses
// back into Go types:
//
//	name, err := evmcsoltypes.ParseSolStringToString(hexResult)
//	balance, err := evmcsoltypes.ParseSolUintToDecimal(hexResult)
package evmcsoltypes
