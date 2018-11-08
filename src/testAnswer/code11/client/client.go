package main

import (
	"google.golang.org/grpc"
	"log"
	"golang.org/x/net/context"
	pb "testAnswer/code11/search"
	"bufio"
	"os"
	"io"
)

func main() {
	conn,err := grpc.Dial("localhost:50051",grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSearchServiceClient(conn)

	ctx := context.Background()
	stream,err :=c.Search(ctx)
	if err != nil {
		log.Printf("创建数据流失败 : %v\n",err)
	}

	//启动一个goroutine 结束输入数据
	go func() {
		log.Println("请输入消息")
		r := bufio.NewReader(os.Stdin)
		for {
			//以回车作为结束
			result,_ := r.ReadString('\n')
			if err = stream.Send(&pb.Request{Input:result});err != nil {
				log.Fatalf("client send error :%v",err)
				return
			}
		}
	}()

	//接受服务器返回的数据
	for {
		resp,err := stream.Recv()
		if err == io.EOF {
			log.Println("接受到服务器的结束信号")
			break
		}
		if err != nil {
			log.Printf("接受消息出错 %v",err)
		}

		log.Printf("【客户端收到】:%s",resp.Outpout)
	}

}
