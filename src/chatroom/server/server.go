package main

import (
	"net"
	"fmt"
	"io"
	"strings"
	"time"
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
	isQuite := make(chan bool)
	hasData := make(chan bool)

	go func() {
		result := make([]byte, 1024*2)

		for {
			n, err := conn.Read(result)

			if n == 0 {
				fmt.Println("conn.Read err:",err)
				isQuite <- true
				return
			}

			if err != nil {
				if err == io.EOF {
					continue
				}
			}

			msg := string(result[:n])
			//fmt.Println("len(msg) ",len(msg))

			if len(msg) == 3 && msg == "who" {
				conn.Write([]byte("user list : \n"))
				for _, temp := range onlineMap {
					msg = temp.Addr + ":" + temp.Name
					conn.Write([]byte(msg + "\n"))
				}
			} else if len(msg) >= 8 && msg[:6] == "rename" {
				cli.Name = strings.Split(msg, "|")[1]
				onlineMap[cliAddr] = cli
				conn.Write([]byte("rename ok \n"))

			} else {
				message <- MakeMsg(cli, msg)
			}
			hasData <- true
		}
	}()

	message <- MakeMsg(cli, "login")

	go func() {
		for msg := range cli.C {
			conn.Write([]byte(msg + "\n"))
		}
	}()

	for {
		select {
		case <-isQuite:
			delete(onlineMap, cliAddr)
			message <- MakeMsg(cli, "logout")
			return
		case <-hasData:
		case <-time.After(30 * time.Second):
			delete(onlineMap, cliAddr)
			message <- MakeMsg(cli, "time out leave out")
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
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
