package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	ws2 "gowebsocket/ws"
	"time"
)
var (
	ws = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter,r *http.Request)()  {
	var (
		wsConn *websocket.Conn
		data []byte
		conn *ws2.Connection
	)
	wsConn,err := ws.Upgrade(w,r,nil)
	if err != nil {
		goto ERR
	}

	conn = ws2.InitConnection(wsConn)

	go func() {
		var err error
		for {
			err = conn.WriteMessage([]byte("heartbeat"))
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		if data,err = conn.ReadMessage();err != nil {
			goto ERR
		}

		if err = conn.WriteMessage(data);err != nil {
			goto ERR
		}
	}

	ERR:
		conn.Close()
}

func main() {
	http.HandleFunc("/ws",wsHandler)
	http.ListenAndServe(":7777",nil)
}
