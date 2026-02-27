// Package evmctypes defines the core data types returned by EVM-compatible
// blockchain RPC methods.
//
// It includes types for blocks, transactions, receipts, logs, and trace
// results. Each type implements custom JSON unmarshaling to handle the
// hex-encoded values commonly found in Ethereum JSON-RPC responses.
//
// Unmarshaling implementations are separated into dedicated *_unmarshaling.go
// files alongside their type definitions.
package evmctypes
