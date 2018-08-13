package main

import (
	"net/http"
	"github.com/gorilla/websocket"
)
var (
	ws = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter,r *http.Request)()  {
	conn,err := ws.Upgrade(w,r,nil)
	if err != nil {
		goto ERR
	}
	for {
		if _,data,err := conn.ReadMessage();err != nil {
			goto ERR
		}else {
			conn.WriteMessage(websocket.TextMessage,data)
		}
	}

	ERR:
		conn.Close()
}

func main() {
	http.HandleFunc("/ws",wsHandler)
	http.ListenAndServe(":7777",nil)
}
