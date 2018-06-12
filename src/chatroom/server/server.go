package main

import (
	"net"
	"fmt"
	"io"
)

type Client struct {
	C    chan string
	Name string
	Addr string
}

var onlineMap map[string]Client
var message = make(chan string)

func MakeMsg(cli Client, msg string) (buf string) {
	buf = "[" + cli.Addr + "]" + cli.Name + ":" + msg
	return
}

func HandleMessage() {
	onlineMap = make(map[string]Client)
	for {
		msg := <-message
		for _, cli := range onlineMap {
			cli.C <- msg
		}
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	cliAddr := conn.RemoteAddr().String()


	cli := Client{make(chan string), cliAddr, cliAddr}

	onlineMap[cliAddr] = cli
	//处理客服端发送过来的消息

	go func() {
		result := make([]byte,1024 * 2)

		for {
			n,err := conn.Read(result)
			if err != nil {
				if err == io.EOF {
					continue
				}
			}

			message <- MakeMsg(cli,string(result[:n]))

		}

	}()

	message <- MakeMsg(cli, "login")

	go func() {
		for msg := range cli.C {
			conn.Write([]byte(msg + "\n"))
		}
	}()

	for {

	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	go HandleMessage()


	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		go HandleConn(conn)
	}
}
