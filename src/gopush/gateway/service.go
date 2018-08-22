package gateway

import (
	"net/http"
	"net"
	"crypto/tls"
	"gopush/common"
	"strconv"
	"fmt"
	"time"
	"encoding/json"
	"unicode/utf8"
)

type Service struct {
	server *http.Server
}

var (
	G_service *Service
)

//全量推送
func handlerPushAll(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		items  string
		msgArr []json.RawMessage
		msgIds int
	)

	if err = r.ParseForm(); err != nil {
		return
	}

	items = r.PostForm.Get("items")
	if err = json.Unmarshal([]byte(items), &msgArr); err != nil {
		fmt.Println(err)
		return
	}
	for msgIds = range msgArr {
		G_merger.PushAll(&msgArr[msgIds])
	}

}

func handlerPushRoom(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		room   string
		items  string
		msgArr []json.RawMessage
		msgIdx int
	)
	if err = r.ParseForm(); err != nil {
		return
	}

	room = r.PostForm.Get("room")
	items = r.PostForm.Get("items")

	if utf8.RuneCountInString(room) <= 0 || room == "" {
		return
	}

	if err = json.Unmarshal([]byte(items), &msgArr); err != nil {
		fmt.Println(err)
		return
	}

	for msgIdx = range msgArr {
		G_merger.PushRoom(room, &msgArr[msgIdx])
	}
}

//统计
func handlerStats(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		data []byte
	)
	if data, err = G_stats.Dump(); err != nil {
		return
	}
	w.Write(data)
}

func InitService() error {
	var (
		err      error
		mux      *http.ServeMux
		server   *http.Server
		listener net.Listener
	)
	//路由
	mux = http.NewServeMux()
	mux.HandleFunc("/push/all", handlerPushAll)
	mux.HandleFunc("/push/room", handlerPushRoom)
	mux.HandleFunc("/stats", handlerStats)
	//TLS证书解析验证
	if _, err = tls.LoadX509KeyPair(G_config.ServerPem, G_config.ServerKey); err != nil {
		err = common.ERR_CERT_NOT_INVALID
		return err
	}
	//http2 服务
	server = &http.Server{
		ReadTimeout:  time.Duration(G_config.ServerReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ServerWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	//监听端口
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ServerPort)); err != nil {
		fmt.Println(err)
		err = common.ERR_SERVER_RUN_FAIL
		return err
	}

	//赋值全局变量
	G_service = &Service{
		server: server,
	}
	//拉起服务
	go server.ServeTLS(listener, G_config.ServerPem, G_config.ServerKey)

	return err
}
