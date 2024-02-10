package repository

import (
	"block_chain/types"
	"context"
	"time"

	"github.com/hacpy/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateNewWallet(wallet *types.Wallet) error {
	ctx := context.Background()
	wallet.Time = uint64(time.Now().Unix())

	opt := options.Update().SetUpsert(true)
	filter := bson.M{"privateKey": wallet.PrivateKey}
	update := bson.M{"$set": bson.M{
		"privateKey": wallet.PrivateKey,
		"publicKey":  wallet.PublicKey,
		"balance":    wallet.Balance,
		"time":       wallet.Time,
	}}

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

func (r *Repository) UpsertWalletsWhenTransfer(from, to, fromBalance, toBalance string) error {
	ctx := context.Background()
	opt := options.Update().SetUpsert(true)

	// TODO : 한 곳에서 2가지 쿼리를 사용하는게 좋은 구조는 아니니까 메서드를 2번 호출하도록 변경등의 리팩토링

	if from != (common.Address{}.String()) {
		// MintCoin이 아닌 경우
		// from 지갑 update
		filter := bson.M{"publicKey": from}
		update := bson.M{"$set": bson.M{"balance": fromBalance}}
		if _, err := r.wallet.UpdateOne(ctx, filter, update, opt); err != nil {
			return err
		}
	}

	// to 지갑 update
	filter := bson.M{"publicKey": to}
	update := bson.M{"$set": bson.M{
		"balance": toBalance,
	}}

	if _, err := r.wallet.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	}

	return nil
}
