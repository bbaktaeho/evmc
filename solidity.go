package evmc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

var (
	solUint256, _    = abi.NewType("uint256", "", nil)
	solUint32, _     = abi.NewType("uint32", "", nil)
	solUint16, _     = abi.NewType("uint16", "", nil)
	solString, _     = abi.NewType("string", "", nil)
	solBool, _       = abi.NewType("bool", "", nil)
	solBytes, _      = abi.NewType("bytes", "", nil)
	solBytes32, _    = abi.NewType("bytes32", "", nil)
	solAddress, _    = abi.NewType("address", "", nil)
	solUint64Arr, _  = abi.NewType("uint64[]", "", nil)
	solAddressArr, _ = abi.NewType("address[]", "", nil)
	solInt8, _       = abi.NewType("int8", "", nil)

	solUint256Args    = abi.Arguments{{Type: solUint256}}
	solUint32Args     = abi.Arguments{{Type: solUint32}}
	solUint16Args     = abi.Arguments{{Type: solUint16}}
	solStringArgs     = abi.Arguments{{Type: solString}}
	solBoolArgs       = abi.Arguments{{Type: solBool}}
	solBytesArgs      = abi.Arguments{{Type: solBytes}}
	solBytes32Args    = abi.Arguments{{Type: solBytes32}}
	solAddressArgs    = abi.Arguments{{Type: solAddress}}
	solUint64ArrArgs  = abi.Arguments{{Type: solUint64Arr}}
	solAddressArrArgs = abi.Arguments{{Type: solAddressArr}}
	solInt8Args       = abi.Arguments{{Type: solInt8}}
)

func parseSolStringToStr(solReturn string) (string, error) {
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

func parseSolBytesToStr(solReturn string) (string, error) {
	b, err := hexutil.Decode(solReturn)
	if err != nil {
		return "", err
	}
	unpacked, err := solBytesArgs.Unpack(b)
	if err != nil {
		return "", err
	}
	return *abi.ConvertType(unpacked[0], new(string)).(*string), nil
}

func parseSolBytes32ToStr(solReturn string) (string, error) {
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

func parseSolUint256ToDecimal(solReturn string) (decimal.Decimal, error) {
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
