package main

import (
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	pb "testAnswer/code11/search"
	"io"
	"strconv"
)

const (
	port = ":50051"
)

type service struct {
}

func (s *service) Search(stream pb.SearchService_SearchServer) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("接受客户端发送context的终止信号")
			return ctx.Err()
		default:
			result,err := stream.Recv()
			if err == io.EOF {
				log.Println("客户端发送的数据流结束")
				return nil
			}
			if err != nil {
				log.Println("客户端发送的数据流出错")
				return err
			}
			//正常接受
			switch result.Input {
			case "结束对话\n":
				log.Println("收到结束对话")
				if err = stream.Send(&pb.Response{Outpout:"收到结束指令\n"});err !=nil {
					return err
				}
				return nil
			case "返回数据流\n":
				log.Println("收到返回数据流")
				for i := 0; i < 10; i++ {
					if err= stream.Send(&pb.Response{Outpout:"数据流#"+strconv.Itoa(i)});err != nil {
						return err
					}
				}
			default:
				//返回数据
				log.Println("收到消息:",result.Input)
				if err= stream.Send(&pb.Response{Outpout:"服务端返回:"+result.Input});err != nil {
					return err
				}
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listener:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSearchServiceServer(s, &service{})
	reflection.Register(s)
	if err = s.Serve(listener); err != nil {
		log.Fatal("server error ", err)
	}
}
