package gateway

import (
	"sync"
	"gopush/common"
	"log"
)

type Bucket struct {
	rwMutex sync.RWMutex
	index   int                     //我是第几个桶
	id2conn map[uint64]*WSConnection //连接列表 key=连接唯一ID
	rooms   map[string]*Room
}

func InitBucket(bucketIdx int) (bucket *Bucket) {
	bucket = &Bucket{
		index:   bucketIdx,
		id2conn: make(map[uint64]*WSConnection),
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
	log.Println("是否存在房间 = ",existed)
	if !existed {
		return
	}
	room.Push(wsMsg)
}

func (bucket *Bucket) AddConn(wsConn *WSConnection) {
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()
	bucket.id2conn[wsConn.curConnId] = wsConn
}

func (bucket *Bucket) DelConn(wsConn *WSConnection) {
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()
	delete(bucket.id2conn,wsConn.curConnId)
}

func (bucket *Bucket) JoinRoom(roomId string,wsConn *WSConnection) (err error) {
	var (
		exists bool
		room *Room
	)

	//判断房间是否存在
	if room,exists = bucket.rooms[roomId];!exists {
		//不存在则创建
		room = InitRoom(roomId)
		bucket.rooms[roomId] = room
		RoomCount_INCR()
	}
	//加入房间
	err = room.Join(wsConn)
	return
}

func (bucket *Bucket) LeaveRoom(roomId string,wsConn *WSConnection) (err error) {
	var (
		exists bool
		room *Room
	)
	if room,exists = bucket.rooms[roomId];!exists {
		err = common.ERR_LEAVE_ROOM_UNEXIST
	}
	err = room.Leave(wsConn)
	//如果房间的数量为空 则删除
	if room.Count() == 0 {
		delete(bucket.rooms,roomId)
		RoomCount_DESC()
	}
	return
}
