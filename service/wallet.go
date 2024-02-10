package service

import (
	"block_chain/types"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"

	"github.com/hacpy/go-ethereum/common/hexutil"
	"github.com/hacpy/go-ethereum/crypto"
)

func newWallet() (string, string, error) {
	p256 := elliptic.P256()

	// 암호화 알고리즘(타원 곡선 등)을 통해 메세지 서명 or 키를 암호화
	if private, err := ecdsa.GenerateKey(p256, rand.Reader); err != nil {
		return "", "", err
	} else if private == nil {
		return "", "", errors.New("PrivateKey is nil")
	} else {
		privateKeyBytes := crypto.FromECDSA(private)
		// [199 53 4 254 73 254 144 31 246 125 125 43 129 46 1 123 233 248 108 128 139 190 223 230 38 209 176 21 213 70 205 67]

		/*
			hexutil의 유용함.
			일반적으로 블록체인은 대소문자를 구별하는데, 구별 안하는 곳에서 사용이 된다면 오류 여지 있다.
			hexutil 라이브러리를 사용하여 오류 발생하지 않도록 할 수 있음
			DB작업에서도 편리함. (find query 등)
		*/
		privateKey := hexutil.Encode(privateKeyBytes)
		// 0xc73504fe49fe901ff67d7d2b812e017be9f86c808bbedfe626d1b015d546cd43
		// -> 실제 사용 가능한 private 지갑 주소 생성 완료

		// hex string으로 변환 한 것을 다시 []byte로 변환 해야지 publicKey 생성 가능.
		// []byte 상태에서 바로 privateKey 만드는건 panic : invalid on Operation..
		againPrivateKey, err := crypto.HexToECDSA(privateKey[2:])
		if err != nil {
			return "", "", err
		}
		cPublicKey := againPrivateKey.Public()
		publicKeyECDSA, ok := cPublicKey.(*ecdsa.PublicKey)

		if !ok {
			return "", "", errors.New("error casting Public Key type")
		}

		publicKeyCommonAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
		publicKey := hexutil.Encode(publicKeyCommonAddress[:])

		return privateKey, publicKey, nil
	}

	// return "", "", nil
}

func (s *Service) MakeWallet() *types.Wallet {
	// var wallet types.Wallet
	// wallet.Balance = ""

	// 초기화할 값이 있으니까 위에 처럼 "타입 추론" 안 시켜줘도 될듯
	wallet := types.Wallet{
		Balance: "0",
	}
	var err error

	if wallet.PrivateKey, wallet.PublicKey, err = newWallet(); err != nil {
		return nil
	} else if err = s.repository.CreateNewWallet(&wallet); err != nil {
		return nil
	}

	return &wallet
}

func (s *Service) GetWallet(pk string) (*types.Wallet, error) {
	if wallet, err := s.repository.GetWallet(pk); err != nil {
		return nil, err
	} else {
		return wallet, err
	}
}
