package kaiatypes

// 블록 헤더 타입
// https://github.com/kaiachain/kaia-sdk/blob/dev/web3rpc/rpc-specs/components/schemas/common/Common.yaml 참고
// https://github.com/kaiachain/kaia-sdk/blob/dev/web3rpc/rpc-specs/components/schemas/common/KaiaTransactionTypes.yaml 참고
type Header struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Miner            string   `json:"miner"`
	Difficulty       string   `json:"difficulty"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             string   `json:"size"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Timestamp        string   `json:"timestamp"`
	BaseFeePerGas    string   `json:"baseFeePerGas,omitempty"`
	Proposer         string   `json:"proposer,omitempty"`
	Committee        []string `json:"committee,omitempty"`
}

// 블록 타입 (트랜잭션 해시만 포함)
type Block struct {
	Header
	Transactions []string `json:"transactions"`
	Uncles       []string `json:"uncles"`
}

// 블록 타입 (트랜잭션 전체 객체 포함)
type BlockIncTxs struct {
	Header
	Transactions []Transaction `json:"transactions"`
	Uncles       []string      `json:"uncles"`
}

// 트랜잭션 타입
type Transaction struct {
	Hash                 string        `json:"hash"`
	Nonce                string        `json:"nonce"`
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	TransactionIndex     string        `json:"transactionIndex"`
	From                 string        `json:"from"`
	To                   string        `json:"to"`
	Value                string        `json:"value"`
	GasPrice             string        `json:"gasPrice"`
	Gas                  string        `json:"gas"`
	Input                string        `json:"input"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
	Type                 string        `json:"type,omitempty"`
	TxType               string        `json:"txType,omitempty"`
	FeePayer             string        `json:"feePayer,omitempty"`
	FeePayerSignatures   []Signature   `json:"feePayerSignatures,omitempty"`
	Signatures           []Signature   `json:"signatures,omitempty"`
	ChainId              string        `json:"chainId,omitempty"`
	AccessList           []AccessTuple `json:"accessList,omitempty"`
	MaxFeePerGas         string        `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas,omitempty"`
}

// 트랜잭션 서명 타입
type Signature struct {
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`
}

// 트랜잭션 영수증(Receipt) 타입
type Receipt struct {
	TransactionHash    string      `json:"transactionHash"`
	TransactionIndex   string      `json:"transactionIndex"`
	BlockHash          string      `json:"blockHash"`
	BlockNumber        string      `json:"blockNumber"`
	From               string      `json:"from"`
	To                 string      `json:"to"`
	CumulativeGasUsed  string      `json:"cumulativeGasUsed"`
	GasUsed            string      `json:"gasUsed"`
	ContractAddress    string      `json:"contractAddress,omitempty"`
	Logs               []Log       `json:"logs"`
	LogsBloom          string      `json:"logsBloom"`
	Status             string      `json:"status"`
	Type               string      `json:"type,omitempty"`
	EffectiveGasPrice  string      `json:"effectiveGasPrice,omitempty"`
	FeePayer           string      `json:"feePayer,omitempty"`
	FeePayerSignatures []Signature `json:"feePayerSignatures,omitempty"`
	TxType             string      `json:"txType,omitempty"`
	ChainId            string      `json:"chainId,omitempty"`
	Root               string      `json:"root,omitempty"`
	RevertReason       string      `json:"revertReason,omitempty"`
	RevertReasonHex    string      `json:"revertReasonHex,omitempty"`
}

// 로그(Log) 타입
type Log struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	BlockNumber string   `json:"blockNumber"`
	TxHash      string   `json:"transactionHash"`
	TxIndex     string   `json:"transactionIndex"`
	BlockHash   string   `json:"blockHash"`
	LogIndex    string   `json:"logIndex"`
	Removed     bool     `json:"removed"`
}

// 필요에 따라 아래 타입은 유지/삭제
// type Reward struct { ... }
// type Committee struct { ... }
// type Council struct { ... }
// type SendTransactionRequest struct { ... }
// type SendTransactionResponse struct { ... }
// type ResendRequest struct { ... }
// type ResendResponse struct { ... }
// type PendingTransaction struct { ... }
// type DecodedAnchoringTransaction struct { ... }

type AccessTuple struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

type Rewards struct {
	BurntFee int64              `json:"burntFee"`
	KGF      int64              `json:"kgf"`
	KIR      int64              `json:"kir"`
	Minted   float64            `json:"minted"`
	Proposer float64            `json:"proposer"`
	Rewards  map[string]float64 `json:"rewards"`
	Stakers  int64              `json:"stakers"`
	TotalFee int64              `json:"totalFee"`
}
