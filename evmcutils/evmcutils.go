package evmcutils

import (
	"github.com/bbaktaeho/evmc"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/crypto"
)

type method struct {
	selector *abi.SelectorMarshaling
	id       string
}

const (
	defaultMethodCacheCapacity = 100
)

var (
	methodCache *lru.Cache[string, *method]
)

func init() {
	methodCache = lru.NewCache[string, *method](defaultMethodCacheCapacity)
}

func GenerateTxInput(funcSig string, args ...evmc.SolType) (string, error) {
	m, ok := methodCache.Get(funcSig)
	if !ok {
		selector, err := abi.ParseSelector(funcSig)
		if err != nil {
			return "", err
		}
		keccak := crypto.Keccak256([]byte(funcSig))
		id := hexutil.Encode(keccak[:])[0:10]
		m = &method{selector: &selector, id: id}
		methodCache.Add(funcSig, m)
	}
	idBytes := hexutil.MustDecode(m.id)
	input := make([]byte, len(idBytes))
	copy(input, idBytes)
	for _, arg := range args {
		input = append(input, arg.([]byte)...)
	}
	return hexutil.Encode(input), nil
}
