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
	Time       uint64 `json:"time" bson:"time"`
}

type Block struct {
	Time         int64          `json:"time"`
	Hash         []byte         `json:"hash"`
	PrevHash     []byte         `json:"from"`
	Nonce        int64          `json:"nonce"`
	Height       int64          `json:"height"`
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
