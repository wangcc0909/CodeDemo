package main

import (
	"net"
	"net/rpc/jsonrpc"
	"rpc"
	"fmt"
	"log"
)

func main() {
	conn, err := net.Dial("tcp", ":8123")
	if err != nil {
		panic(err)
	}

	client := jsonrpc.NewClient(conn)

	var reply float64
	err = client.Call("JsonServer.Div", rpcdemo.Args{X: 1, Y: 1}, &reply)
	if err != nil {
		log.Printf("err : %v", err)
	} else {
		fmt.Printf("result = %v", reply)
	}
}
