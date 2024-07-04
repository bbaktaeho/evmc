package evmcutils

import (
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
	sigCache   *lru.Cache[string, string]
	sigPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*\((\w+(\[\d*\])?(,\w+(\[\d*\])?)*)?\)$`)
)

func init() {
	sigCache = lru.NewCache[string, string](defaultMethodCacheCapacity)
}

func GenerateTxInput(funcSig string, args ...evmcsoltypes.SolType) (string, error) {
	sigHash, err := getSig(funcSig)
	if err != nil {
		return "", err
	}
	sid := sigHash[:10]
	idBytes := hexutil.MustDecode(sid)
	input := make([]byte, len(idBytes))
	copy(input, idBytes)
	for _, arg := range args {
		input = append(input, arg.([]byte)...)
	}
	return hexutil.Encode(input), nil
}

func GenerateLogTopic(eventSig string) (string, error) {
	sigHash, err := getSig(eventSig)
	if err != nil {
		return "", err
	}
	return sigHash, nil
}

func getSig(sig string) (string, error) {
	sigHash, ok := sigCache.Get(sig)
	if !ok {
		if !sigPattern.MatchString(sig) {
			return "", ErrInvalidSig
		}
		keccak := crypto.Keccak256([]byte(sig))
		hash := hexutil.Encode(keccak[:])
		sigCache.Add(sig, hash)
		sigHash = hash
	}
	return sigHash, nil
}
