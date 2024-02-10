package service

import (
	"block_chain/types"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	"github.com/hacpy/go-ethereum/common/hexutil"
)

// 채굴 작업

type PowWork struct {
	Block      *types.Block `json:"block"`
	Target     *big.Int     `json:"target"`
	Difficulty int64        `json:"difficulty"`
}

// PowWork를 생성
func (s *Service) NewPow(b *types.Block) *PowWork {
	t := new(big.Int).SetInt64(1)

	// 비트 마스크 연산이라 보면 됨.
	// t.Lsh(t,1) 이면 1<<1 : 2 ,  t.Lsh(t,2)면 1<<2 : 2^2
	t.Lsh(t, uint(256-s.difficulty))

	return &PowWork{Block: b, Target: t, Difficulty: s.difficulty}

}

// 채굴작업
func (p *PowWork) RunMining() (int64, string) {
	var iHash big.Int
	var hash [32]byte

	nonce := 0

	for nonce < math.MaxInt64 {
		// 이전해시, nonce, difficulty를 통해 해싱될 데이터 생성
		d := p.makeHash(nonce)

		// sha256의 체크섬 반환 (해싱)
		hash = sha256.Sum256(d)

		fmt.Printf("\r%x", hash)

		// 체크섬을 장착? 한다고 볼까
		iHash.SetBytes(hash[:])

		// 타겟보다 작은 수를 찾아가는 것.
		if iHash.Cmp(p.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println()
	return int64(nonce), hexutil.Encode(hash[:])
}

func (p *PowWork) makeHash(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(p.Block.PrevHash),
			// ToDO make Transaction To Byte
			HashTransaction(p.Block),
			intToHex(p.Difficulty),
			intToHex(int64(nonce)),
		},
		[]byte{}, // 구분자
	)
}

func intToHex(number int64) []byte {
	b := new(bytes.Buffer)

	if err := binary.Write(b, binary.BigEndian, number); err != nil {
		panic(err)
	} else {
		return b.Bytes()
	}
}
