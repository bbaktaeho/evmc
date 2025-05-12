package evmc

type ChainID uint64

const (
	EthereumMainnet ChainID = 1
	EthereumSepolia ChainID = 11155111
	EthereumHolesky ChainID = 17000
	EthereumHoodi   ChainID = 560048

	PolygonMainnet ChainID = 137
	PolygonAmoy    ChainID = 80002

	ArbitrumMainnet ChainID = 42161
	ArbitrumSepolia ChainID = 421614

	OptimismMainnet ChainID = 10
	OptimismSepolia ChainID = 11155420

	BaseMainnet ChainID = 8453
	BaseSepolia ChainID = 84532

	KaiaMainnet ChainID = 8217
	KaiaKairos  ChainID = 1001

	ChilizMainnet ChainID = 88888
	ChilizSpicy   ChainID = 88882
)

var (
	EthereumIDs = []ChainID{EthereumMainnet, EthereumSepolia, EthereumHolesky, EthereumHoodi}
	PolygonIDs  = []ChainID{PolygonMainnet, PolygonAmoy}
	ArbitrumIDs = []ChainID{ArbitrumMainnet, ArbitrumSepolia}
	// OpStackIDs is built on Optimism’s open-source OP Stack.
	OpStackIDs = []ChainID{OptimismMainnet, OptimismSepolia, BaseMainnet, BaseSepolia}
	KaiaIDs    = []ChainID{KaiaMainnet, KaiaKairos}
	ChilizIDs  = []ChainID{ChilizMainnet, ChilizSpicy}
)

func (id ChainID) Uint64() uint64 {
	return uint64(id)
}

func (id ChainID) Name() string {
	switch id {
	case EthereumMainnet:
		return "ethereum-mainnet"
	case EthereumSepolia:
		return "ethereum-sepolia"
	case EthereumHolesky:
		return "ethereum-holesky"
	case EthereumHoodi:
		return "ethereum-hoodi"
	case PolygonMainnet:
		return "polygon-mainnet"
	case PolygonAmoy:
		return "polygon-amoy"
	case ArbitrumMainnet:
		return "arbitrum-mainnet"
	case ArbitrumSepolia:
		return "arbitrum-sepolia"
	case OptimismMainnet:
		return "optimism-mainnet"
	case OptimismSepolia:
		return "optimism-sepolia"
	case BaseMainnet:
		return "base-mainnet"
	case BaseSepolia:
		return "base-sepolia"
	case KaiaMainnet:
		return "kaia-mainnet"
	case KaiaKairos:
		return "kaia-kairos"
	case ChilizMainnet:
		return "chiliz-mainnet"
	case ChilizSpicy:
		return "chiliz-spicy"
	default:
		return "unknown"
	}
}
