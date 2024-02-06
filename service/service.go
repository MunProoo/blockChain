package service

import (
	"block_chain/repository"

	"github.com/inconshreveable/log15"
)

type Service struct {
	// config     *config.Config
	log        log15.Logger
	difficulty int64

	repository *repository.Repository
}

func NewService(repository *repository.Repository, difficulty int64) *Service {
	s := &Service{
		log:        log15.New("module", "service"),
		repository: repository,
		difficulty: difficulty,
	}

	return s
}
