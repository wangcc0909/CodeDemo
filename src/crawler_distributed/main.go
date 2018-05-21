package main

import (
	"crawler/engine"
	"crawler/simple"
	"crawler/zhenai/parser"
	"crawler_distributed/persist/client"
	worker "crawler_distributed/worker/client"
	"net/rpc"
	"crawler_distributed/supportrpc"
	"log"
	"flag"
	"strings"
)

var itemServerPort = flag.String("item_port","","the port for me to listen on")
var workHosts = flag.String("work_hosts","","work hosts (comma separated)")

func main() {
	flag.Parse()
	itemChan,err := client.ItemSever(*itemServerPort)
	if err != nil {
		panic(err)
	}

	p := workerPool(strings.Split(*workHosts,","))
	processor := worker.WorkProcessor(p)

	e := engine.CurrentEngine{
		Scheduler:&simple.QueuedScheduler{},
		WorkerCount:100,
		ItemServer:itemChan,
		Processor:processor,
	}

	e.Run(engine.Request{
		Url:"http://www.zhenai.com/zhenghun",
		Parser:engine.NewFuncParser(parser.ParserCityList,"ParserCityList"),
	})
}

func workerPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _,h := range hosts{
		c,err := supportrpc.ClientRpc(h)
		if err != nil {
			log.Printf("host create client error %v: %v",h,err)
			continue
		}
		clients = append(clients,c)
	}

	out := make(chan *rpc.Client)
	go func() {
		for {
			for _,c := range clients{
				out <- c
			}
		}
	}()
	return out
}
