package main

import (
	"gopkg.in/olivere/elastic.v3"
	"crawler_distributed/supportrpc"
	"crawler_distributed/persist"
	"log"
	"crawler_distributed/config"
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

	log.Fatal(saveRpc(fmt.Sprintf(":%d",*port),config.RpcIndex))

}

func saveRpc(host,index string) error {
	client,err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://192.168.99.100:9200"),
	)

	if err != nil {
		return err
	}

	return supportrpc.ServeRpc(host,&persist.ItemSaveServer{
		Client:client,
		Index:index,
	})
}
