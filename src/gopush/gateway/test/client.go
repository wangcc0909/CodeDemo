package main

import (
	"flag"
	"log"
	"time"
	"github.com/gorilla/websocket"
	"net/url"
	"fmt"
)

var addr = flag.String("addr","localhost:7777","http server address")

func loop() {
	for {
		u := url.URL{Scheme:"ws",Host:*addr,Path:"/connect"}
		c,_,err := websocket.DefaultDialer.Dial(u.String(),nil)
		if err != nil {
			continue
		}
		//循环读消息
		for {
			_,data,err :=c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			log.Println(string(data))
		}
		c.Close()
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	for i := 0; i < 10000; i++ {
		go loop()
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
