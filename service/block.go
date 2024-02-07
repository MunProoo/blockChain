package service

import (
	"block_chain/types"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hacpy/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBlock(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	var pHash []byte
	var block *types.Block
	latestBlock, err := s.repository.GetLatestBlock()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genesis Block Will be Created")

			genesisMessage := "THis is First Genesis Block"

			// message와 등등 정보를 담은 트랜잭션 생성
			tx := createTransaction(genesisMessage, "0x3710954186c28084f2190ee367a59a66e5b584c9", "", "", 1)

			// 트랜잭션, 해시 담은 블록 생성
			block = createBlockInner([]*types.Transaction{tx}, pHash, height)

			// 연산 생성
			pow := s.NewPow(block)

			// mining
			block.Nonce, block.Hash = pow.RunMining()
		} else {
			// 그냥 에러
			s.log.Crit("Failed to Get Latest Block", "err", err)
			return nil
		}
	} else {
		pHash = latestBlock.Hash

		// create new block
		block = createBlockInner(txs, pHash, height)
		pow := s.NewPow(block)

		// mining
		block.Nonce, block.Hash = pow.RunMining()
	}

	if err = s.repository.SaveBlock(block); err != nil {
		s.log.Crit("Failed to Save Block", "err", err)
		panic(err)
	}

	return block
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

func createTransaction(message, from, to, amount string, block int64) *types.Transaction {
	data := struct {
		Message string `json:"message"`
		From    string `json:"from"`
		To      string `json:"to"`
		Amount  string `json:"amount"`
	}{
		Message: message,
		From:    from,
		To:      to,
		Amount:  amount,
	}

	// hex로 변경
	dataToSign := fmt.Sprintf("%x", data)
	pk := "726d28fcad721dbd6d3badcec9b9cf8e98ee341ad1955e800a49fa55570a25a0"

	if ecdsaPrivateKey, err := crypto.HexToECDSA(pk); err != nil {
		panic(err)
	} else if r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivateKey, []byte(dataToSign)); err != nil {
		// From -> To로 Amount 만큼 이동한다는 해시에 개인키를 통해 서명(데이터 인증, 무결성 보장, 신원 확인)
		panic(err)
	} else {
		signature := append(r.Bytes(), s.Bytes()...)

		//
		return &types.Transaction{
			Block:   block,
			Time:    time.Now().Unix(),
			From:    from,
			To:      to,
			Amount:  amount,
			Message: message,
			Tx:      hex.EncodeToString(signature),
		}
	}
}

// transaction -> byte encoding (해쉬화?)
func HashTransaction(b *types.Block) []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		var encoded bytes.Buffer

		enc := gob.NewEncoder(&encoded)
		if err := enc.Encode(tx); err != nil {
			panic(err)
		} else {
			txHashes = append(txHashes, encoded.Bytes())
		}
	}

	// 하나의 값만 달라져도 최종값이 달라지는 트리 -> 머클 트리
	// 자식 노드 중 하나만 달라져도 루트 값이 달라져버린다
	tree := NewMerkleTree(txHashes)
	return tree.RootNode.Data
}
