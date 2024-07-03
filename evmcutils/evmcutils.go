package evmcutils

import (
	"errors"
	"regexp"

	"github.com/bbaktaeho/evmc/evmcsoltypes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/crypto"
)

// TODO: geth abi.ABI
const (
	defaultMethodCacheCapacity = 100
)

var (
	methodCache   *lru.Cache[string, string]
	methodPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*\((\w+(\[\d*\])?(,\w+(\[\d*\])?)*)?\)$`)
)

func init() {
	methodCache = lru.NewCache[string, string](defaultMethodCacheCapacity)
}

func GenerateTxInput(funcSig string, args ...evmcsoltypes.SolType) (string, error) {
	mid, ok := methodCache.Get(funcSig)
	if !ok {
		if !methodPattern.MatchString(funcSig) {
			return "", errors.New("invalid function signature")
		}
		keccak := crypto.Keccak256([]byte(funcSig))
		id := hexutil.Encode(keccak[:])[0:10]
		mid = id
		methodCache.Add(funcSig, id)
	}
	idBytes := hexutil.MustDecode(mid)
	input := make([]byte, len(idBytes))
	copy(input, idBytes)
	for _, arg := range args {
		input = append(input, arg.([]byte)...)
	}
	return hexutil.Encode(input), nil
}
