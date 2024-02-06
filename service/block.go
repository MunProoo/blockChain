package service

import (
	"block_chain/types"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBlock(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	var pHash []byte

	if latestBlock, err := s.repository.GetLatestBlock(); err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genesis Block Will be Created")

			newBlock := createBlockInner(txs, pHash, height)
			pow := s.NewPow(newBlock)

			// mining
			newBlock.Nonce, newBlock.Hash = pow.RunMining()

			return newBlock
		} else {
			// 그냥 에러
			s.log.Crit("Failed to Get Latest Block", "err", err)
		}
	} else {
		pHash = latestBlock.Hash

		// create new block
		newBlock := createBlockInner(txs, pHash, height)
		pow := s.NewPow(newBlock)

		// mining
		newBlock.Nonce, newBlock.Hash = pow.RunMining()

		return newBlock
	}

	return nil
}

func createBlockInner(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	return &types.Block{
		Time:         time.Now().Unix(),
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
	}
}
