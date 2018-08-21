package gateway

import (
	"time"
	"gopush/common"
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
)

//每隔一秒检测是否健康连接
func (wsConn *WSConnection) heartbeatChecker() {
	var (
		timer *time.Timer
	)
	timer = time.NewTimer(time.Duration(G_config.WsHeartbeatInternal) * time.Second)
	for {
		select {
		case <-timer.C:
			if !wsConn.isAlive() {
				wsConn.Closed()
			}
			timer.Reset(time.Duration(G_config.WsHeartbeatInternal) * time.Second)
		case <-wsConn.closeChan:
			timer.Stop()
		}
	}
}

func (wsConn *WSConnection) LeaveAll() {
	var (
		roomId string
	)

	//退出所有房间
	for roomId = range wsConn.rooms {
		G_connMgr.LeaveRoom(roomId,wsConn)
		delete(wsConn.rooms, roomId)
	}
}

func (wsConn *WSConnection) handleLeave(bizReq *common.BizMessage) (bizResp *common.BizMessage, err error) {
	var (
		bizLeaveData common.BizLeaveData
		exists bool
	)
	//解析请求数据
	if err = json.Unmarshal(bizReq.Data,&bizLeaveData);err != nil {
		fmt.Println(err)
		return
	}
	//判断房间是否存在
	if len(bizLeaveData.Room) == 0 {
		err = common.ERR_ROOM_ID_INVALID
		return
	}
	//判断是否在房间中
	if exists = wsConn.rooms[bizLeaveData.Room];!exists {
		err = common.ERR_LEAVE_ROOM_UNEXIST
		return
	}
	//从连接池中移除
	if err = G_connMgr.LeaveRoom(bizLeaveData.Room,wsConn);err != nil {
		return
	}
	//删除房间和连接池之间的关系
	delete(wsConn.rooms,bizLeaveData.Room)
	return
}

//处理Join请求
func (wsConn *WSConnection) handleJoin(bizReq *common.BizMessage) (bizResp *common.BizMessage, err error) {
	var (
		bizJoinData common.BizJoinData
		exists bool
	)
	//解析请求数据
	if err = json.Unmarshal(bizReq.Data,&bizJoinData);err != nil {
		fmt.Println(err)
		return
	}
	//判断房间是否处在
	if len(bizJoinData.Room) == 0 {
		err = common.ERR_ROOM_ID_INVALID
		return
	}
	//判断是否超过房间数量的上限
	if len(wsConn.rooms) > G_config.MaxJoinRoom {
		err = common.ERR_MAX_ROOM
		return
	}
	//判断是否已在房间中
	if _,exists = wsConn.rooms[bizJoinData.Room];exists {
		err = common.ERR_JOINED_ROOM
		return
	}
	//建立和房间的联系
	if err = G_connMgr.JoinRoom(bizJoinData.Room,wsConn);err != nil {
		return
	}
	//建立和房间的联系
	wsConn.rooms[bizJoinData.Room] = true
	return
}

//处理Ping请求
func (wsConn *WSConnection) handlePing(bizReq *common.BizMessage) (bizResp *common.BizMessage, err error) {
	var (
		buf []byte
	)
	wsConn.keepAlive()
	if buf, err = json.Marshal(common.BizPongData{}); err != nil {
		return
	}
	bizResp = &common.BizMessage{
		Type: "PONG",
		Data: json.RawMessage(buf),
	}
	return
}

func (wsConn *WSConnection) WSHandle() {
	var (
		wsMsg   *common.WSMessage
		bizReq  *common.BizMessage
		bizResp *common.BizMessage
		err     error
		buf     []byte
	)

	//连接加入管理器 可以推送端查到
	G_connMgr.AddConn(wsConn)
	//心跳检测
	go wsConn.heartbeatChecker()

	//请求协程处理
	for {
		if wsMsg, err = wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		//只处理文本消息
		if wsMsg.MsgType != websocket.TextMessage {
			continue
		}

		//解析消息体
		if bizReq, err = common.DecodeBizMessage(wsMsg.MsgData); err != nil {
			goto ERR
		}

		bizResp = nil
		switch bizReq.Type {
		case "PING":
			if bizResp, err = wsConn.handlePing(bizReq); err != nil {
				goto ERR
			}
		case "JOIN":
			if bizResp, err = wsConn.handleJoin(bizReq); err != nil {
				goto ERR
			}
		case "LEAVE":
			if bizResp, err = wsConn.handleLeave(bizReq); err != nil {
				goto ERR
			}
		}
		if bizResp != nil {
			if buf, err = json.Marshal(*bizResp); err != nil {
				goto ERR
			}
			//socket缓冲区写满不是致命错误
			if err = wsConn.SendMessage(&common.WSMessage{MsgType: websocket.TextMessage, MsgData: buf}); err != nil {
				if err == common.ERR_SEND_MESSAGE_FULL {
					goto ERR
				} else {
					err = nil
				}
			}
		}
	}

ERR:
//确保连接关闭
	wsConn.Closed()
	wsConn.LeaveAll()
	//从连接池中移除
	G_connMgr.DelConn(wsConn)
}
