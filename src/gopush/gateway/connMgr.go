package gateway

import (
	"gopush/common"
)

//推送任务
type PushJob struct {
	PushType int                `json:"pushType"`
	RoomID   string             `json:"roomId"` //房间ID
	BizMgr   *common.BizMessage `json:"bizMgr"` //未序列化的消息
	WsMgr    *common.WSMessage  `json:"wsMgr"`  //已序列化的消息
}

//连接管理器
type ConnMgr struct {
	Buckets      []*Bucket
	JobChan      []chan *PushJob //每个bucket对应一个job Queue
	DispatchChan chan *PushJob
}

var (
	G_connMgr *ConnMgr
)

//消息分发到bucket
func (connMgr *ConnMgr) jobDispatchMain(dispatchIdx int) {
	var (
		bucketIdx int
		pushJob   *PushJob
		err       error
	)

	for {
		select {
		case pushJob = <-connMgr.DispatchChan:
			DispatchPending_DESC()
			//序列化
			if pushJob.WsMgr, err = common.EncodeWsMessage(pushJob.BizMgr); err != nil {
				continue
			}
			//分发给所有bucket 若bucket阻塞则等待
			for bucketIdx = range connMgr.Buckets {
				PushJobPending_INCR()
				connMgr.JobChan[bucketIdx] <- pushJob
			}
		}
	}
}

//job负责消息广播给客户端
func (connMgr *ConnMgr) jobWorkerMain(jobWorkerInx int, bucketIdx int) {
	var (
		bucket  = connMgr.Buckets[bucketIdx]
		pushJob *PushJob
	)

	for {
		select {
		case pushJob = <-connMgr.JobChan[bucketIdx]: //从bucket的job Queue取出一个任务
			PushJobPending_DESC()
			if pushJob.PushType == common.TYPE_PUSH_ROOM {
				bucket.PushRoom(pushJob.RoomID, pushJob.WsMgr)
			} else if pushJob.PushType == common.TYPE_PUSH_ALL {
				bucket.PushAll(pushJob.WsMgr)
			}
		}
	}
}

func InitConnMgr() (err error) {
	var (
		bucketIdx    int
		jobWorkerIdx int
		dispatchIdx  int
		connMgr      *ConnMgr
	)
	connMgr = &ConnMgr{
		Buckets:      make([]*Bucket, G_config.BucketCount),
		JobChan:      make([]chan *PushJob, G_config.BucketCount),
		DispatchChan: make(chan *PushJob, G_config.DispatchChannelSize),
	}

	for bucketIdx = range connMgr.Buckets {
		connMgr.Buckets[bucketIdx] = InitBucket(bucketIdx)                              //初始化bucket
		connMgr.JobChan[bucketIdx] = make(chan *PushJob, G_config.BucketJobChannelSize) //bucket队列长度
		//bucket的jobWorker
		for jobWorkerIdx = 0; jobWorkerIdx < G_config.BucketJobWorkerCount; jobWorkerIdx++ {
			go connMgr.jobWorkerMain(jobWorkerIdx, bucketIdx)
		}
	}

	//初始化分发协程,用于将消息扇出各个bucket
	for dispatchIdx = 0; dispatchIdx < G_config.DispatchChannelSize; dispatchIdx++ {
		go connMgr.jobDispatchMain(dispatchIdx)
	}
	G_connMgr = connMgr
	return
}
