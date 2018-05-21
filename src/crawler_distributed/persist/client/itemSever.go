package client

import (
	"crawler_distributed/supportrpc"
	"crawler/model"
	"log"
	"crawler_distributed/config"
)

func ItemSever(host string) (chan model.Item,error){

	client,err := supportrpc.ClientRpc(host)

	if err != nil {
		return nil,err
	}

	out := make(chan model.Item)
	go func() {
		count := 0
		for {
			item := <-out
			log.Printf("item sava:%d, %v",count,item)
			count++
			result := ""
			err := client.Call(config.RpcServerMethod,item,&result)

			if err != nil {
				log.Printf("client save item err %v, item %v",err,item)
			}
		}
	}()

	return out,nil
}
