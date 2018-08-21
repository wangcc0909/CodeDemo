package gateway

import (
	"gopush/common"
	"sync"
)

type Room struct {
	rwMutex sync.RWMutex
	roomId string
	id2Conn map[string]*WSConnection
}

func (room *Room) Push(wsMgr *common.WSMessage) {

}
