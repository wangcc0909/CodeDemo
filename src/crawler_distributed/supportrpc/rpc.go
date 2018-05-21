package supportrpc

import (
	"net/rpc"
	"net"
	"log"
	"net/rpc/jsonrpc"
)

func ServeRpc(host string,server interface{}) error {
	err := rpc.Register(server)
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	log.Printf("listening on %s",host)
	for {
		conn,err := listen.Accept()
		if err != nil {
			log.Printf("accept error %v",err)
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
	return nil
}

func ClientRpc(host string) (*rpc.Client,error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil,err
	}

	return jsonrpc.NewClient(conn),nil
}