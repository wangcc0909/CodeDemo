package gateway

import (
	"gopush/common"
	"sync"
)

type Room struct {
	rwMutex sync.RWMutex
	roomId string
	id2Conn map[uint64]*WSConnection
}

func (room *Room) Push(wsMgr *common.WSMessage) {
	var (
		wsConn *WSConnection
	)
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	for _,wsConn = range room.id2Conn {
		wsConn.SendMessage(wsMgr)
	}
}

func InitRoom(roomId string) *Room {
	return &Room{
		roomId:roomId,
		id2Conn:make(map[uint64]*WSConnection),
	}
}

func (room *Room) Join(wsConn *WSConnection) error {
	var (
		err error
		exists bool
	)
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	if _,exists = room.id2Conn[wsConn.curConnId];exists {
		err = common.ERR_JOIN_ROOM_TWICE
	}
	room.id2Conn[wsConn.curConnId] = wsConn
	return err
}

func (room *Room) Leave(wsConn *WSConnection) error {
	var (
		err error
		exists bool
	)
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	if _,exists = room.id2Conn[wsConn.curConnId];!exists {
		err = common.ERR_LEAVE_ROOM_UNEXIST
		return err
	}
	delete(room.id2Conn,wsConn.curConnId)
	return err
}

func (room *Room) Count() int {
	var (
		count int
	)
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	count = len(room.id2Conn)
	return count
}
