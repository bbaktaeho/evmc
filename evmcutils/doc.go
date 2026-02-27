// Package evmcutils provides utility functions for working with EVM-compatible
// blockchain data.
//
// It includes helpers for generating transaction input data from Solidity
// function signatures and computing event log topics using Keccak-256 hashing.
//
//	input, err := evmcutils.GenerateTxInput(
//	    "transfer(address,uint256)",
//	    evmcsoltypes.Address("0x..."),
//	    evmcsoltypes.Uint256(amount),
//	)
package evmcutils
