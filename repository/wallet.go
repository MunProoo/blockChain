package repository

import (
	"block_chain/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateNewWallet(wallet *types.Wallet) error {
	ctx := context.Background()
	wallet.Time = uint64(time.Now().Unix())

	opt := options.Update().SetUpsert(true)
	filter := bson.M{"privateKey": wallet.PrivateKey}
	update := bson.M{"$set": wallet}

	if _, err := r.wallet.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetWallet(pk string) (*types.Wallet, error) {
	ctx := context.Background()
	filter := bson.M{"privateKey": pk}
	var wallet types.Wallet

	if err := r.wallet.FindOne(ctx, filter, options.FindOne()).Decode(&wallet); err != nil {
		return nil, err
	} else {
		return &wallet, nil
	}
}

// 내가 작성한거.. 수정 할 수도.
func (r *Repository) GetWalletByPublicKey(publicKey string) (*types.Wallet, error) {
	ctx := context.Background()
	filter := bson.M{"publicKey": publicKey}
	var wallet types.Wallet

	if err := r.wallet.FindOne(ctx, filter, options.FindOne()).Decode(&wallet); err != nil {
		return nil, err
	} else {
		return &wallet, nil
	}

}
