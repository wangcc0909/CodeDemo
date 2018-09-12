package main

import (
	"flag"
	"runtime"
	"gopush/logic"
	"os"
	"fmt"
	"time"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "config", "src/gopush/logic/cli/logic.json", "where is logic.json.")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	initArgs()
	initEnv()

	if err = logic.InitConfig(confFile); err != nil {
		goto ERR
	}

	if err = logic.InitStats(); err != nil {
		goto ERR
	}

	if err = logic.InitGateConnMgr();err != nil {
		goto ERR
	}

	if err = logic.InitService();err != nil {
		goto ERR
	}

	fmt.Println("运行")
	for {
		time.Sleep(1 * time.Second)
	}
	return
ERR:
	fmt.Println(err)
	os.Exit(-1)
	return
}
