package main

import (
	"flag"
	"runtime"
	"gopush/gateway"
	"os"
	"fmt"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "config", "src/gopush/gateway/cli/gateway.json", "where gateway.json is.")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	//初始化环境
	initArgs()
	initEnv()

	//加载配置
	if err = gateway.InitConfig(confFile); err != nil {
		fmt.Println(err)
		goto ERR
	}
	//统计
	if err = gateway.InitStats();err != nil {
		fmt.Println(err)
		goto ERR
	}

	ERR:
		os.Exit(-1)
}
