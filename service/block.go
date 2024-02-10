package service

import (
	"block_chain/types"
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/hacpy/go-ethereum/common"
	"github.com/hacpy/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBlock(from, to, value string) {
	var block *types.Block
	fromBalance := "0"
	toBalance := "0"

	if latestBlock, err := s.repository.GetLatestBlock(); err != nil {
		// 첫 블록 생성 작업 (따로 빼도 되지만 우선은 이대로 작업하자)
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genesis Block Will be Created")
			genesisMessage := "THis is First Genesis Block"

			if pk, _, err := newWallet(); err != nil {
				panic(err)
			} else {
				// message와 등등 정보를 담은 트랜잭션 생성
				tx := createTransaction(genesisMessage, common.Address{}.String(), to, value, pk, 1)
				// 트랜잭션, 해시 담은 블록 생성
				block = createBlockInner([]*types.Transaction{tx}, "", 1)
			}
		}
	} else {
		// 기존 블록이 있는 경우 : Mint, Transfer 처리
		var tx *types.Transaction

		if common.HexToAddress(from) == (common.Address{}) {
			// Mint
			if pk, _, err := newWallet(); err != nil {
				panic(err)
			} else {
				// message와 등등 정보를 담은 트랜잭션 생성
				tx = createTransaction("Mint Coin", common.Address{}.String(), to, value, pk, 1)

				// to에 대한 wallet가져와서 balance 업데이트 해주기
				wallet, err := s.repository.GetWalletByPublicKey(to)
				if err != nil {
					panic(err)
				}

				toDecimalBalance, _ := decimal.NewFromString(wallet.Balance)
				valueDecimal, _ := decimal.NewFromString(value)
				toDecimalBalance = toDecimalBalance.Add(valueDecimal)

				toBalance = toDecimalBalance.String()
			}
		} else {
			// Transfer
			if wallet, err := s.repository.GetWalletByPublicKey(from); err != nil {
				panic(err)
			} else if toWallet, err := s.repository.GetWalletByPublicKey(to); err != nil {
				if err == mongo.ErrNoDocuments {
					s.log.Debug("Failed to Find wallet. PublicKey is Nil", "to", to)
				} else {
					panic(err)
				}
				return
			} else {
				// 0. 본인이 본인에게 인지 체크
				if strings.EqualFold(wallet.PublicKey, to) {
					s.log.Debug("Same Address", "from's public", wallet.PublicKey, "to", to)
					return
				}
				// 1. From의 밸런스(잔고) 체크
				// 2. From에서 차감, To Balacne에 증가
				fromDecimalBalance, _ := decimal.NewFromString(wallet.Balance)
				toBalanceDecimal, _ := decimal.NewFromString(toWallet.Balance)
				valueDecimal, _ := decimal.NewFromString(value)

				if fromDecimalBalance.Cmp(valueDecimal) == -1 { // From 잔고 충분치 못함
					s.log.Debug("Failed to transfer Coin by from Balance", "From", from, "balance", wallet.Balance, "amount", value)
					return
				} else {
					fromDecimalBalance = fromDecimalBalance.Sub(valueDecimal)
					fromBalance = fromDecimalBalance.String()
				}

				toBalanceDecimal = toBalanceDecimal.Add(valueDecimal)
				toBalance = toBalanceDecimal.String()

				tx = createTransaction("Transfer Coin", from, to, value, wallet.PrivateKey, 1)

			}

		}
		// create new block
		block = createBlockInner([]*types.Transaction{tx}, latestBlock.Hash, latestBlock.Height+1)
	}

	// 채굴
	// block 특정, 연산법 특정
	pow := s.NewPow(block)
	// mining
	block.Nonce, block.Hash = pow.RunMining()

	// 채굴 종료, 잔고 업데이트
	if err := s.repository.UpsertWalletsWhenTransfer(from, to, fromBalance, toBalance); err != nil {
		panic(err)
	}

	if err := s.repository.SaveBlock(block); err != nil {
		s.log.Crit("Failed to Save Block", "err", err)
		panic(err)
	}
	s.log.Info("Mining is Successed")
}

func createBlockInner(txs []*types.Transaction, prevHash string, height int64) *types.Block {
	return &types.Block{
		Time:         time.Now().Unix(),
		Hash:         "",
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
	}
}

func createTransaction(message, from, to, amount, pk string, block int64) *types.Transaction {
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
	dataToSign := fmt.Sprintf("%x\n", data)

	pk = strings.TrimPrefix(pk, "0x")

	if ecdsaPrivateKey, err := crypto.HexToECDSA(pk); err != nil {
		// error(*errors.errorString) *{s: "invalid hex character 'x' in private key"}
		// 0x로 시작하도록 저장중이므로 pk를 그대로 꺼내서 쓰면 위의 에러가 발생함.
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
