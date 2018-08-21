package gateway

import (
	"gopush/common"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type WSConnection struct {
	mutex             sync.Mutex
	curConnId         uint64
	wsSocket          *websocket.Conn
	inChan            chan *common.WSMessage
	outChan           chan *common.WSMessage
	closeChan         chan byte
	isClosed          bool
	lastHeartbeatTime time.Time       //最近一次的心跳时间
	rooms             map[string]bool //加入了那些房间
}

func InitWSConnection(connId uint64, conn *websocket.Conn) (wsConn *WSConnection) {
	wsConn = &WSConnection{
		curConnId:         connId,
		wsSocket:          conn,
		inChan:            make(chan *common.WSMessage, G_config.WsInChannelSize),
		outChan:           make(chan *common.WSMessage, G_config.WsOutChannelSize),
		closeChan:         make(chan byte, 1),
		lastHeartbeatTime: time.Now(),
		rooms:             make(map[string]bool),
	}

	go wsConn.readLoop()
	go wsConn.writeLoop()
	return
}

func (wsConn *WSConnection) SendMessage(wsMsg *common.WSMessage) (err error) {
	select {
	case wsConn.outChan <- wsMsg:
		SendMessageTotal_INCR()
	case <-wsConn.closeChan:
		err = common.ERR_CONNECTION_LOSS
	default:
		//写操作不会阻塞 ,因为channel已经预留给websocket一定缓冲空间
		err = common.ERR_SEND_MESSAGE_FULL
		SendMessageFail_INCR()
	}
	return
}

func (wsConn *WSConnection) ReadMessage() (wsMsg *common.WSMessage, err error) {
	select {
	case wsMsg = <-wsConn.inChan:
	case <-wsConn.closeChan:
		err = common.ERR_CONNECTION_LOSS
	}
	return
}

func (wsConn *WSConnection) readLoop() {
	var (
		msgType int
		data    []byte
		wsMsg   *common.WSMessage
		err     error
	)
	for {
		if msgType, data, err = wsConn.wsSocket.ReadMessage(); err != nil {
			goto ERR
		}
		wsMsg = common.BuildMessage(msgType, data)
		select {
		case wsConn.inChan <- wsMsg:
		case <-wsConn.closeChan:
			goto CLOSED
		}
	}
ERR:
	wsConn.Closed()
CLOSED:
}

func (wsConn *WSConnection) writeLoop() {
	var (
		wsMsg *common.WSMessage
		err   error
	)
	for {
		select {
		case wsMsg = <-wsConn.outChan:
		case <-wsConn.closeChan:
			goto CLOSED
		}
		if err = wsConn.wsSocket.WriteMessage(wsMsg.MsgType, wsMsg.MsgData); err != nil {
			goto ERR
		}
	}
ERR:
	wsConn.Closed()
CLOSED:
}

func (wsConn *WSConnection) Closed() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}
