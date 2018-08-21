package gateway

import (
	"sync"
	"gopush/common"
)

type Bucket struct {
	rwMutex sync.RWMutex
	index   int                     //我是第几个桶
	id2conn map[int64]*WSConnection //连接列表 key=连接唯一ID
	rooms   map[string]*Room
}

func InitBucket(bucketIdx int) (bucket *Bucket) {
	bucket = &Bucket{
		index:   bucketIdx,
		id2conn: make(map[int64]*WSConnection),
		rooms:   make(map[string]*Room),
	}
	return
}

//推送给bucket内所有的用户
func (bucket *Bucket) PushAll(wsMsg *common.WSMessage) {
	var (
		wsConn *WSConnection
	)
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()
	//全量非阻塞推送
	for _,wsConn = range bucket.id2conn {
		wsConn.SendMessage(wsMsg)
	}
}

func (bucket *Bucket) PushRoom(roomId string,wsMsg *common.WSMessage) {
	var (
		room *Room
		existed bool
	)

	//锁bucket
	bucket.rwMutex.Lock()
	room,existed = bucket.rooms[roomId]
	bucket.rwMutex.Unlock()

	if !existed {
		return
	}
	room.Push(wsMsg)
}
