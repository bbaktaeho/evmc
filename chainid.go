package evmc

type ChainID uint64

const (
	EthereumMainnet ChainID = 1
	EthereumSepolia ChainID = 11155111
	EthereumHolesky ChainID = 17000

	PolygonMainnet ChainID = 137
	PolygonAmoy    ChainID = 80002

	ArbitrumMainnet ChainID = 42161
	ArbitrumSepolia ChainID = 421614

	OptimismMainnet ChainID = 10
	OptimismSepolia ChainID = 11155420

	BaseMainnet ChainID = 8453
	BaseSepolia ChainID = 84532
)

var (
	EthereumIDs = []ChainID{EthereumMainnet, EthereumSepolia, EthereumHolesky}
	PolygonIDs  = []ChainID{PolygonMainnet, PolygonAmoy}
	ArbitrumIDs = []ChainID{ArbitrumMainnet, ArbitrumSepolia}
	// OpStackIDs is built on Optimism’s open-source OP Stack.
	OpStackIDs = []ChainID{OptimismMainnet, OptimismSepolia, BaseMainnet, BaseSepolia}
)

func (id ChainID) Name() string {
	switch id {
	case EthereumMainnet:
		return "ethereum-mainnet"
	case EthereumSepolia:
		return "ethereum-sepolia"
	case EthereumHolesky:
		return "ethereum-holesky"
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
	default:
		return "unknown"
	}
}
