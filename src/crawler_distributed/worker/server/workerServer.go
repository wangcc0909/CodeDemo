package main

import (
	"crawler_distributed/supportrpc"
	"crawler_distributed/worker"
	"log"
	"flag"
	"fmt"
)

var port = flag.Int("port",0,"the port for me listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		log.Printf("no port listener")
		return
	}
	log.Fatal(supportrpc.ServeRpc(fmt.Sprintf(":%d",*port),worker.WorkServer{}))
}
