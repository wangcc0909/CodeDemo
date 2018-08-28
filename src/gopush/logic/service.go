package logic

import (
	"net/http"
	"net"
	"encoding/json"
	"fmt"
	"unicode/utf8"
	"strings"
	"time"
	"strconv"
	"log"
)

type service struct {
	server *http.Server
}

var (
	G_service *service
)

//全量推送 msg = {}
func handlePushAll(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		items  string
		msgArr []json.RawMessage
	)

	if err = r.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}

	items = r.PostForm.Get("items")
	log.Println(items)
	if err = json.Unmarshal([]byte(items), &msgArr); err != nil {
		fmt.Println("err = ",err.Error())
		return
	}
	G_GateConnMgr.PushAll(msgArr)
}

func handlePushRoom(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		items  string
		room   string
		msgArr []json.RawMessage
	)

	if err = r.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}
	room = r.PostForm.Get("room")
	items = r.PostForm.Get("items")
	room = strings.TrimSpace(room)
	if utf8.RuneCountInString(room) == 0 || room == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("room不能为空"))
		return
	}

	if err = json.Unmarshal([]byte(items), &msgArr); err != nil {
		fmt.Println(err)
		return
	}
	G_GateConnMgr.PushRoom(room, msgArr)
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	var (
		data []byte
		err  error
	)

	if data, err = G_stats.Dump(); err != nil {
		return
	}
	w.Write(data)
}

func InitService() (err error) {
	var (
		mux      *http.ServeMux
		server   *http.Server
		listener net.Listener
	)

	//路由
	mux = http.NewServeMux()
	mux.HandleFunc("/push/all", handlePushAll)
	mux.HandleFunc("/push/room", handlePushRoom)
	mux.HandleFunc("/stats", handleStats)

	//http/1服务
	server = &http.Server{
		ReadTimeout:  time.Duration(G_config.ServiceReadTimeout) * time.Second,
		WriteTimeout: time.Duration(G_config.ServiceWriteTimeout) * time.Second,
		Handler:      mux,
	}

	//端口
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ServicePort)); err != nil {
		fmt.Println(err)
		return
	}
	//赋值全局变量
	G_service = &service{
		server: server,
	}
	//拉起服务
	go server.Serve(listener)
	return
}
