package main

import (
	"net"
	"fmt"
	"io"
)

var msg = make(chan string)
func HandleServerMsg(conn net.Conn) {
	buf := make([]byte,1024 * 2)
	for  {
		n,err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
		}

		fmt.Println(string(buf[:n]))
	}
}

func main() {
	conn,err := net.Dial("tcp","127.0.0.1:8000")
	if err != nil {
		fmt.Println("net.Dial err :",err)
		return
	}
	defer conn.Close()

	go HandleServerMsg(conn)

	//给服务器发送数据
	go SendMsgToServer(conn)
	var in string
	for {
		_,err = fmt.Scan(&in)
		if err != nil {
			fmt.Println("fmt.Scan err :",err)
		}
		msg <- in
	}


}
func SendMsgToServer(conn net.Conn) {
	for {
		message := <-msg

		conn.Write([]byte(message))
	}
}
