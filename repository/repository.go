package repository

import (
	"block_chain/config"
	"context"

	"github.com/inconshreveable/log15"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
	db     *mongo.Database

	config *config.Config
	log    log15.Logger
}

func NewRepository(config *config.Config) (*Repository, error) {
	r := &Repository{
		config: config,
		log:    log15.New("module", "repository"),
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

	r.db = r.client.Database(mConfig.DB)

	// TODO -> Collection 변경

	r.log.Error("Success to repository", "uri", mConfig.Uri, "db", mConfig.DB)

	return r, nil
}
