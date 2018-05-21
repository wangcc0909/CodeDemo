package main

import (
	"net"
	"log"
	"net/rpc/jsonrpc"
	"net/rpc"
	"rpc"
)

func main() {

	err := rpc.Register(rpcdemo.JsonServer{})
	if err != nil {
		panic(err)
	}

	listen, err := net.Listen("tcp", ":8123")
	if err != nil {
		panic(err)
	}

	for {
		conn,err := listen.Accept()
		if err != nil {
			log.Printf("accept error %v",err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}
