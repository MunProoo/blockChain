package types

const (
	CreateWallet  = "CreateWallet"
	TransferCoin  = "TransferCoin"
	MintCoin      = "MintCoin"
	ConnectWallet = "ConnectWallet"
	ChangeWallet  = "ChangeWallet"
)

type Wallet struct {
	PrivateKey string `json:"privateKey" bson:"privateKey"`
	PublicKey  string `json:"publicKey" bson:"publicKey"`
	Balance    string `json:"balance" bson:"balance"` // mongo에서 Decimal을 지원하지 않음, uint64를 사용하기엔 decimal이 너무 크기 때문에 string 사용
	Time       uint64 `json:"time" bson:"time"`
}

type Block struct {
	Time         int64          `json:"time"`
	Hash         string         `json:"hash"`
	PrevHash     string         `json:"from"`
	Nonce        int64          `json:"nonce"`
	Height       int64          `json:"height"` // 블록의 크기 or 트랜잭션이 얼마나 담겼나를 의미. 블록의 순서 X
	Transactions []*Transaction `json:"transaction"`
}

type Transaction struct {
	Block   int64  `json:"block"`
	Time    int64  `json:"time"`
	From    string `json:"from"`
	To      string `json:"to"`
	Amount  string `json:"amount"`
	Message string `json:"message"`
	Tx      string `json:"tx"`
}
