package types

const (
	CreateWallet = "CreateWallet"
	TransferCoin = "TransferCoin"
	MintCoin     = "MintCoin"
)

type Wallet struct {
	PrivateKey string `json:"privateKey" bson:"privateKey"`
	PublicKey  string `json:"publicKey" bson:"publicKey"`
	Time       uint64 `json:"time" bson:"time"`
}
