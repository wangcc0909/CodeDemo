package gateway

import (
	"net/http"
	"github.com/gorilla/websocket"
	"net"
	"time"
	"strconv"
	"sync/atomic"
	"fmt"
)

type WSServer struct {
	server *http.Server
	curConnID uint64
}

var (
	G_wsServer *WSServer
	//允许所有CORS跨域请求
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func handlerConnect(w http.ResponseWriter, r *http.Request) {
	var (
		wsSocket *websocket.Conn
		err error
		connId uint64
		wsConn *WSConnection
	)
	//websocket握手
	if wsSocket,err = upgrader.Upgrade(w,r,nil);err != nil {
		return
	}

	//连接的唯一标识
	connId = atomic.AddUint64(&G_wsServer.curConnID,1)
	//初始化websocket读写协程
	wsConn = InitWSConnection(connId,wsSocket)
	//开始处理websocket的消息
	wsConn.WSHandle()
}

func InitWsServer() (err error) {
	var (
		mux *http.ServeMux
		server *http.Server
		listener net.Listener
	)
	//路由
	mux = http.NewServeMux()
	mux.HandleFunc("/connect",handlerConnect)

	//HTTP服务
	server = &http.Server{
		ReadTimeout:time.Duration(G_config.WsReadTimeout) * time.Millisecond,
		WriteTimeout:time.Duration(G_config.WsWriteTimeout) * time.Millisecond,
		Handler:mux,
	}

	//监听端口
	if listener,err = net.Listen("tcp",":" + strconv.Itoa(G_config.WsPort));err != nil {
		fmt.Println(err)
		return
	}

	//赋值全局变量
	G_wsServer = &WSServer{
		server:server,
		curConnID:uint64(time.Now().Unix()),
	}

	//拉起服务
	go server.Serve(listener)
	return
}
