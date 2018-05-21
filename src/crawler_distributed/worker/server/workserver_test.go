package main

import (
	"testing"
	"crawler_distributed/supportrpc"
	"crawler_distributed/worker"
	"time"
	"log"
)

func TestWorkServer(t *testing.T) {
	go supportrpc.ServeRpc(":9000",worker.WorkServer{})
	time.Sleep(time.Second)

	client,err := supportrpc.ClientRpc(":9000")
	if err != nil {
		panic(err)
	}

	request := worker.Request{
		Url:"http://album.zhenai.com/u/108906739",
		SerializeParser:worker.SerializeParser{
			Name:"ProfileParser",
			Args:"安静的雪",
		},
	}

	result := worker.ParserResult{}

	err = client.Call("WorkServer.Worker",request,&result)
	if err == nil {
		log.Printf("test result %v",result)
	}else{
		t.Errorf("test error %v,%v",result,err)
	}


}
