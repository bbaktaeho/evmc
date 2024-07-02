package evmc

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

// TODO: dinamic type

var (
	solUint, _       = abi.NewType("uint256", "", nil)
	solString, _     = abi.NewType("string", "", nil)
	solBool, _       = abi.NewType("bool", "", nil)
	solBytes32, _    = abi.NewType("bytes32", "", nil)
	solBytes, _      = abi.NewType("bytes", "", nil)
	solAddress, _    = abi.NewType("address", "", nil)
	solUint64Arr, _  = abi.NewType("uint64[]", "", nil)
	solAddressArr, _ = abi.NewType("address[]", "", nil)

	solUint256Args    = abi.Arguments{{Type: solUint}}
	solStringArgs     = abi.Arguments{{Type: solString}}
	solBoolArgs       = abi.Arguments{{Type: solBool}}
	solBytesArgs      = abi.Arguments{{Type: solBytes}}
	solBytes32Args    = abi.Arguments{{Type: solBytes32}}
	solAddressArgs    = abi.Arguments{{Type: solAddress}}
	solUint64ArrArgs  = abi.Arguments{{Type: solUint64Arr}}
	solAddressArrArgs = abi.Arguments{{Type: solAddressArr}}
)

type SolType interface{}

func ParseUint256(d decimal.Decimal) SolType {
	res, err := solUint256Args.Pack(d.BigInt())
	if err != nil {
		panic(err)
	}
	return res
}

func ParseAddress(address string) SolType {
	res, err := solAddressArgs.Pack(common.HexToAddress(address))
	if err != nil {
		panic(err)
	}
	return res
}

func parseSolStringToString(solReturn string) (string, error) {
	b, err := hexutil.Decode(solReturn)
	if err != nil {
		return "", err
	}
	unpacked, err := solStringArgs.Unpack(b)
	if err != nil {
		return "", err
	}
	return *abi.ConvertType(unpacked[0], new(string)).(*string), nil
}

func parseSolBytesToString(solReturn string) (string, error) {
	b, err := hexutil.Decode(solReturn)
	if err != nil {
		return "", err
	}
	unpacked, err := solBytesArgs.Unpack(b)
	if err != nil {
		return "", err
	}
	preRes := *abi.ConvertType(unpacked[0], new([]byte)).(*[]byte)
	return string(preRes), nil
}

func parseSolFixedBytesToString(solReturn string) (string, error) {
	if len(solReturn) != 66 {
		solReturn += strings.Repeat("0", 66-len(solReturn))
	}
	b, err := hexutil.Decode(solReturn)
	if err != nil {
		return "", err
	}
	unpacked, err := solBytes32Args.Unpack(b)
	if err != nil {
		return "", err
	}
	preRes := *abi.ConvertType(unpacked[0], new([32]byte)).(*[32]byte)
	return string(preRes[:]), nil
}

func parseSolUintToDecimal(solReturn string) (decimal.Decimal, error) {
	b, err := hexutil.Decode(solReturn)
	if err != nil {
		return decimal.Zero, err
	}
	unpacked, err := solUint256Args.Unpack(b)
	if err != nil {
		return decimal.Zero, err
	}
	preRes := *abi.ConvertType(unpacked[0], new(*big.Int)).(**big.Int)
	return decimal.NewFromBigInt(preRes, 0), nil
}
