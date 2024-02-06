package main

import (
	"block_chain/app"
	"block_chain/config"
	"flag"
	"fmt"
)

var (
	// 프로그램 실행 시 넣어주는 플래그 값  ex ) go run . -environment ./config/environment.toml
	configFlag  = flag.String("environment", "./environment.toml", "environment toml file not found")
	difficuilty = flag.Int("difficulty", 12, "difficulty err")
)

func main() {
	flag.Parse()

	c := config.NewConfig(*configFlag)

	app.NewApp(c, int64(*difficuilty))
	fmt.Println(c.Info.Version)

	fmt.Println("test")
}
