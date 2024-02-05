package app

import (
	"block_chain/config"
	"block_chain/repository"
	"block_chain/service"
	. "block_chain/types"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/inconshreveable/log15"
)

type App struct {
	config *config.Config

	service    *service.Service
	repository *repository.Repository

	log log15.Logger
}

func NewApp(config *config.Config) {
	a := &App{
		config: config,
		log:    log15.New("module", "app"),
	}

	var err error
	if a.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	}

	a.log.Info("Module Started", "time", time.Now().Unix())
	a.service = service.NewService(config, a.repository)

	sc := bufio.NewScanner(os.Stdin)

	useCase()
	for {
		sc.Scan()
		fmt.Println(sc.Text())

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
	fmt.Println("2. ", TransferCoin, " <To> <Amount>")
	fmt.Println("3. ", MintCoin, " <To> <Amount>")

	fmt.Println()
}

func (a *App) inputValueAssessment(input []string) (msg error) {
	if len(input) == 0 {
		msg = errors.New("check Use Case")
		return
	}

	switch input[0] {
	case "1":
		fmt.Println("Create Wallet is switch")
		a.service.MakeWallet()

		fmt.Println("Success to Create Wallet")
	case "2":
		fmt.Println("TransferCoin is switch")
	case "3":
		fmt.Println("MintCoin is switch")
	default:
		return msg
	}

	fmt.Println()

	return nil
}
