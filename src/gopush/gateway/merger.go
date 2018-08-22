package gateway

import (
	"encoding/json"
	"time"
	"gopush/common"
)

type PushContext struct {
	message *json.RawMessage
	room    string
}

type PushBatch struct {
	items      []*json.RawMessage
	commitTime *time.Timer
	room       string
}

type MergerWorker struct {
	mergerType  int //合并类型 广播,room,uid...
	contextChan chan *PushContext
	timeoutChan chan *PushBatch
	room2batch  map[string]*PushBatch //房间合并
	allBatch    *PushBatch            //广播合并
}

//广播消息,房间消息的合并
type Merger struct {
	roomWorker      []*MergerWorker //房间合并
	broadcastWorker *MergerWorker   //广播合并
}

var (
	G_merger *Merger
)

func (worker *MergerWorker) autoCommit(batch *PushBatch) func() {
	return func() {
		worker.timeoutChan <- batch
	}
}

func (worker *MergerWorker) commitBatch(batch *PushBatch) (err error) {
	var (
		bizPushData *common.BizPushData
		bizMessage  *common.BizMessage
		buf         []byte
	)
	bizPushData = &common.BizPushData{
		Items: batch.items,
	}
	if buf, err = json.Marshal(*bizPushData); err != nil {
		return
	}

	bizMessage = &common.BizMessage{
		Type: "PUSH",
		Data: json.RawMessage(buf),
	}
	//打包发送
	if worker.mergerType == common.TYPE_PUSH_ROOM {
		delete(worker.room2batch, batch.room)
		err = G_connMgr.PushRoom(batch.room, bizMessage)
	} else if worker.mergerType == common.TYPE_PUSH_ALL {
		worker.allBatch = nil
		err = G_connMgr.PushAll(bizMessage)
	}

	return
}

func (worker *MergerWorker) mergerWorkerMain() {
	var (
		context      *PushContext
		batch        *PushBatch
		timeoutBatch *PushBatch
		existed      bool
		isCreated    bool
		err          error
	)

	for {
		select {
		case context = <-worker.contextChan:
			MergerPending_DESC()
			isCreated = false
			//按房间合并
			if worker.mergerType == common.TYPE_PUSH_ROOM {
				if batch, existed = worker.room2batch[context.room]; !existed {
					batch = &PushBatch{room: context.room}
					worker.room2batch[context.room] = batch
					isCreated = true
				}
			} else if worker.mergerType == common.TYPE_PUSH_ALL { //广播合并
				batch = worker.allBatch
				if batch == nil {
					batch = &PushBatch{}
					worker.allBatch = batch
					isCreated = true
				}
			}
			//合并消息
			batch.items = append(batch.items, context.message)
			//新建批次 启动超时自动提交
			if isCreated {
				batch.commitTime = time.AfterFunc(time.Duration(G_config.MaxMergerDelay)*time.Millisecond, worker.autoCommit(batch))
			}

			//批次未满 继续等待下次提交
			if len(batch.items) < G_config.MaxMergerBatchSize {
				continue
			}
			//批次已满  取消自动提交
			batch.commitTime.Stop()
		case timeoutBatch = <-worker.timeoutChan:
			if worker.mergerType == common.TYPE_PUSH_ROOM {
				//定时器触发时  批次已经提交
				if batch, existed = worker.room2batch[context.room]; !existed {
					continue
				}
				//定时器触发时,前一个批次已经提交 新批次已经建立
				if batch != timeoutBatch {
					continue
				}
			} else if worker.mergerType == common.TYPE_PUSH_ALL {
				batch = worker.allBatch
				//定时器触发时  批次已经提交
				if batch != timeoutBatch {
					continue
				}
			}
		}
		err = worker.commitBatch(batch)
		//打点统计
		if worker.mergerType == common.TYPE_PUSH_ROOM {
			MergerRoomTotal_INCR(int64(len(batch.items)))
			if err != nil {
				MergerRoomFail_INCR(int64(len(batch.items)))
			}
		} else if worker.mergerType == common.TYPE_PUSH_ALL {
			MergerAllTotal_INCR(int64(len(batch.items)))
			if err != nil {
				MergerAllFail_INCR(int64(len(batch.items)))
			}
		}
	}
}

func (worker *MergerWorker) pushAll(msg *json.RawMessage) (err error) {
	var (
		context *PushContext
	)
	context = &PushContext{
		message: msg,
	}
	select {
	case worker.contextChan <- context:
		MergerPending_INCR()
	default:
		err = common.ERR_MERGER_CHANNEL_FULL
	}
	return err
}

func (worker *MergerWorker) pushRoom(room string, msg *json.RawMessage) (err error) {
	var (
		context *PushContext
	)
	context = &PushContext{
		message: msg,
		room:    room,
	}
	select {
	case worker.contextChan <- context:
		MergerPending_INCR()
	default:
		err = common.ERR_MERGER_CHANNEL_FULL
	}
	return err
}

func initMergerWorker(mergerType int) (mergerWorker *MergerWorker) {
	mergerWorker = &MergerWorker{
		mergerType:  mergerType,
		contextChan: make(chan *PushContext, G_config.MergerChannelSize),
		timeoutChan: make(chan *PushBatch, G_config.MergerChannelSize),
		room2batch:  make(map[string]*PushBatch),
	}

	go mergerWorker.mergerWorkerMain()
	return
}

func InitMerger() (err error) {
	var (
		workerIdx int
		merger    *Merger
	)
	merger = &Merger{
		roomWorker: make([]*MergerWorker, G_config.MergerWorkerCount),
	}
	for workerIdx = 0; workerIdx < G_config.MergerWorkerCount; workerIdx++ {
		merger.roomWorker[workerIdx] = initMergerWorker(common.TYPE_PUSH_ROOM)
	}
	merger.broadcastWorker = initMergerWorker(common.TYPE_PUSH_ALL)
	G_merger = merger
	return
}

//广播合并推送
func (merger *Merger) PushAll(msg *json.RawMessage) (err error) {
	return merger.broadcastWorker.pushAll(msg)
}

//房间合并推送
func (merger *Merger) PushRoom(room string, msg *json.RawMessage) (err error) {
	//计算room hash到某个worker
	var (
		workerIdx uint32 = 0
		ch        byte
	)
	for _, ch = range []byte(room) {
		workerIdx = (workerIdx + uint32(ch)*33) % uint32(G_config.MergerWorkerCount)
	}

	return merger.roomWorker[workerIdx].pushRoom(room, msg)
}
