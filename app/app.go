package app

import (
	"block_chain/config"
	"block_chain/global"
	"block_chain/repository"
	"block_chain/service"
	. "block_chain/types"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hacpy/go-ethereum/common"
	"github.com/inconshreveable/log15"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	config *config.Config

	service    *service.Service
	repository *repository.Repository

	log log15.Logger
}

func NewApp(config *config.Config, difficulty int64) {
	a := &App{
		config: config,
		log:    log15.New("module", "app"),
	}

	var err error
	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	}

	a.log.Info("Module Started", "time", time.Now().Unix())
	a.service = service.NewService(a.repository, difficulty)

	sc := bufio.NewScanner(os.Stdin)

	for {
		useCase()
		from := global.FROM()

		if from != "" {
			a.log.Info("Current Connected Wallet", "from", from)
			fmt.Println()
		}

		sc.Scan()
		// fmt.Println(sc.Text())

		// UseCase 입력 파싱
		input := strings.Split(sc.Text(), " ")
		if err = a.inputValueAssessment(input); err != nil {
			a.log.Error("Failed to call CLI", "err", err, "input", input)
			fmt.Println()
		}

	}
}

func useCase() {
	fmt.Println()

	fmt.Println("This is MunProoo's Module for BlockChain Core With Mongo")
	fmt.Println()
	fmt.Println("Use Case")

	fmt.Println("1. ", CreateWallet)
	fmt.Println("2. ", ConnectWallet, " <PK>")
	fmt.Println("3. ", ChangeWallet, " <PK>")
	fmt.Println("4. ", TransferCoin, " <To> <Amount>")
	fmt.Println("5. ", MintCoin, " <To> <Amount>")

	fmt.Println()
}

func (a *App) inputValueAssessment(input []string) (msg error) {
	if len(input) == 0 {
		msg = errors.New("check Use Case")
		return
	} else {
		from := global.FROM()

		switch input[0] {
		case CreateWallet:
			fmt.Println("Create Wallet Command is inputed")
			if wallet := a.service.MakeWallet(); wallet == nil {
				panic("Failed to create wallet")
			} else {
				a.log.Info("Success to Create Wallet", "pk", wallet.PrivateKey, "pu", wallet.PublicKey)
			}
		case ConnectWallet:
			if from != "" {
				a.log.Debug("Already Connected Wallet", "from", from)
				fmt.Println()
				return
			}

			pk := input[1]
			if wallet, err := a.service.GetWallet(pk); err != nil {
				if err == mongo.ErrNoDocuments {
					a.log.Debug("Failed to Find wallet. PK is Nil", "pk", pk)
				} else {
					a.log.Crit("Failed to Find Wallet", "pk", pk, "err", err)
				}
			} else {
				global.SetFROM(wallet.PublicKey)
				fmt.Println()
				a.log.Info("Success To Connect Wallet", "from", wallet.PublicKey)
			}

		case ChangeWallet:
			if from == "" {
				a.log.Debug("Connect Wallet First")
				fmt.Println()
				return
			}
			pk := input[1]
			if strings.EqualFold(pk, from) {
				a.log.Info("Same Address", "pk", pk)
				fmt.Println()
				return
			}
			if wallet, err := a.service.GetWallet(pk); err != nil {
				if err == mongo.ErrNoDocuments {
					a.log.Debug("Failed to Find wallet. PK is Nil", "pk", pk)
				} else {
					a.log.Crit("Failed to Find Wallet", "pk", pk, "err", err)
				}
			} else {
				global.SetFROM(wallet.PublicKey)
				fmt.Println()
				a.log.Info("Success To Change Wallet", "from", wallet.PublicKey)
			}
		case TransferCoin: // 코인을 건내준다. 이동.
			if from == "" {
				a.log.Debug("Connect Wallet First")
				fmt.Println()
				return
			} else if len(input) < 3 {
				a.log.Debug("Insufficient Input")
				fmt.Println()
				return
			}

			to, value := input[1], input[2]
			if to == "" || value == "" {
				a.log.Debug("Request value, to is uncorrect")
				fmt.Println()
				return
			}
			a.service.CreateBlock(from, to, value)

		case MintCoin: // 관리자 -> 사용자 계좌에 입금 (자산 변환)
			if len(input) < 3 {
				a.log.Debug("Insufficient Input")
				fmt.Println()
				return
			}
			to, value := input[1], input[2]
			if to == "" || value == "" {
				a.log.Debug("Request value, to is uncorrect")
				fmt.Println()
				return
			}
			a.service.CreateBlock((common.Address{}).String(), to, value)

		default:
			return msg
		}

		fmt.Println()
	}

	return nil
}
