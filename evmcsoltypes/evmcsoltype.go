package evmcsoltypes

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
)

var (
	solInt, _        = abi.NewType("int256", "", nil)
	solUint, _       = abi.NewType("uint256", "", nil)
	solString, _     = abi.NewType("string", "", nil)
	solBool, _       = abi.NewType("bool", "", nil)
	solFixedBytes, _ = abi.NewType("bytes32", "", nil)
	solBytes, _      = abi.NewType("bytes", "", nil)
	solAddress, _    = abi.NewType("address", "", nil)
	solUint64Arr, _  = abi.NewType("uint64[]", "", nil)
	solAddressArr, _ = abi.NewType("address[]", "", nil)

	solUint256Args    = abi.Arguments{{Type: solUint}}
	solStringArgs     = abi.Arguments{{Type: solString}}
	solBoolArgs       = abi.Arguments{{Type: solBool}}
	solBytesArgs      = abi.Arguments{{Type: solBytes}}
	solFixedBytesArgs = abi.Arguments{{Type: solFixedBytes}}
	solAddressArgs    = abi.Arguments{{Type: solAddress}}
	solUint64ArrArgs  = abi.Arguments{{Type: solUint64Arr}}
	solAddressArrArgs = abi.Arguments{{Type: solAddressArr}}
)

type SolType interface{}

func Int256(d decimal.Decimal) SolType {
	res, err := solUint256Args.Pack(d.BigInt())
	if err != nil {
		panic(err)
	}
	return res
}

func Uint256(d decimal.Decimal) SolType {
	res, err := solUint256Args.Pack(d.BigInt())
	if err != nil {
		panic(err)
	}
	return res
}

func Address(address string) SolType {
	res, err := solAddressArgs.Pack(common.HexToAddress(address))
	if err != nil {
		panic(err)
	}
	return res
}

func FixedBytes(b []byte) SolType {
	// slice to 32 bytes
	if len(b) > 32 {
		panic("bytes length should be less than 32")
	}
	var fixedBytes [32]byte
	copy(fixedBytes[:32], b)
	res, err := solFixedBytesArgs.Pack(fixedBytes)
	if err != nil {
		panic(err)
	}
	return res
}

func AddressArr(addresses []string) SolType {
	var addrs []common.Address
	for _, address := range addresses {
		addrs = append(addrs, common.HexToAddress(address))
	}
	res, err := solAddressArrArgs.Pack(addrs)
	if err != nil {
		panic(err)
	}
	return res
}

func ParseSolStringToString(solReturn string) (string, error) {
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

func ParseSolBytesToString(solReturn string) (string, error) {
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

func ParseSolFixedBytesToString(solReturn string) (string, error) {
	if len(solReturn) != 66 {
		solReturn += strings.Repeat("0", 66-len(solReturn))
	}
	b, err := hexutil.Decode(solReturn)
	if err != nil {
		return "", err
	}
	unpacked, err := solFixedBytesArgs.Unpack(b)
	if err != nil {
		return "", err
	}
	preRes := *abi.ConvertType(unpacked[0], new([32]byte)).(*[32]byte)
	return string(preRes[:]), nil
}

func ParseSolUintToDecimal(solReturn string) (decimal.Decimal, error) {
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
