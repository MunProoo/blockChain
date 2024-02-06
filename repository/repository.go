package repository

import (
	"block_chain/config"
	"block_chain/types"
	"context"

	"github.com/inconshreveable/log15"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
	wallet *mongo.Collection
	tx     *mongo.Collection
	block  *mongo.Collection

	// config *config.Config
	log log15.Logger
}

func NewRepository(config *config.Config) (*Repository, error) {
	r := &Repository{
		// config: config,
		log: log15.New("module", "repository"),
	}

	var err error
	ctx := context.Background()

	mConfig := config.Mongo
	if r.client, err = mongo.Connect(ctx, options.Client().ApplyURI(mConfig.Uri)); err != nil {
		r.log.Error("Failed to connect to mongo", "uri", mConfig.Uri)
		return nil, err
	} else if err = r.client.Ping(ctx, nil); err != nil {
		r.log.Error("failed to ping to mongo", "uri", mConfig.Uri)
		return nil, err
	}

	db := r.client.Database(mConfig.DB)

	r.wallet = db.Collection("wallet")
	r.tx = db.Collection("tx")
	r.block = db.Collection("block")

	r.log.Info("Success to repository", "uri", mConfig.Uri, "db", mConfig.DB)

	return r, nil
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
