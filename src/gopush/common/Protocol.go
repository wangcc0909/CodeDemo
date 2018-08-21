package common

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

//推送类型
const (
	TYPE_PUSH_ROOM = 1
	TYPE_PUSH_ALL  = 2
)

//webSocket消息对象
type WSMessage struct {
	MsgType int
	MsgData []byte
}

type BizPongData struct {

}

type BizJoinData struct {
	Room string `json:"room"`
}

type BizLeaveData struct {
	Room string `json:"room"`
}

//业务消息的固定格式 (type+data)
type BizMessage struct {
	Type string          `json:"type"` //type 消息类型 Ping Pong Join Leave Push
	Data json.RawMessage `json:"data"`
}

//序列化 返回json格式字符串
func EncodeWsMessage(message *BizMessage) (wsMsg *WSMessage, err error) {
	var (
		buf []byte
	)

	if buf, err = json.Marshal(message); err != nil {
		return
	}

	wsMsg = &WSMessage{
		MsgType: websocket.TextMessage,
		MsgData: buf,
	}
	return
}

func DecodeBizMessage(data []byte) (bizMsg *BizMessage, err error) {
	var (
		bizMsgObj BizMessage
	)

	if err = json.Unmarshal(data, &bizMsgObj); err != nil {
		return
	}
	bizMsg = &bizMsgObj
	return
}

func BuildMessage(msgType int, data []byte) (wsMsg *WSMessage) {
	wsMsg = &WSMessage{
		MsgType: msgType,
		MsgData: data,
	}
	return
}
