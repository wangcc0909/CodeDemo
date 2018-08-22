package main

import (
	"flag"
	"runtime"
	"gopush/gateway"
	"os"
	"fmt"
	"time"
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
	//初始化连接管理器
	if err = gateway.InitConnMgr();err != nil {
		fmt.Println(err)
		goto ERR
	}
	//初始化websocket服务器
	if err = gateway.InitWsServer();err != nil {
		fmt.Println(err)
		goto ERR
	}

	//初始化merger合并层
	if err = gateway.InitMerger();err != nil {
		fmt.Println(err)
		goto ERR
	}

	//初始化service接口
	if err = gateway.InitService();err != nil {
		fmt.Println(err)
		goto ERR
	}

	fmt.Println("运行")

	//主协程不能终止
	for {
		time.Sleep(1 * time.Second)
	}

	ERR:
		os.Exit(-1)
}
