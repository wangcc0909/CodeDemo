package main

import (
	"google.golang.org/grpc"
	"net"
	"log"
	pb "testAnswer/code12/order-server/proto"
	"testAnswer/code12/order-server/handler"
	"google.golang.org/grpc/reflection"
)

const (
	address = ":8080"
)

func main() {
	listener,err := net.Listen("tcp",address)
	if err != nil {
		log.Fatalf("net listen error : %v",err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderServer(s,&handler.Service{})
	reflection.Register(s)
	if err = s.Serve(listener);err != nil {
		log.Fatalf("s server error:%v",err)
	}
}
