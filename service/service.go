package service

import (
	"block_chain/config"
	"block_chain/repository"

	"github.com/inconshreveable/log15"
)

type Service struct {
	config *config.Config
	log    log15.Logger

	repository *repository.Repository
}

func NewService(config *config.Config, repository *repository.Repository) *Service {
	s := &Service{
		config:     config,
		log:        log15.New("module", "service"),
		repository: repository,
	}

	return s
}
