package logic

import (
	"encoding/json"
	"gopush/common"
)

type PushJob struct {
	pushType int               //推送类型
	roomId   string            //房间Id
	message  []json.RawMessage //推送的消息数组
}

type GateConnMgr struct {
	gateConns    []*GateConn   //到所有的gateway连接数组
	pendingChan  []chan byte   //gateway的并发请求控制
	dispatchChan chan *PushJob //带分发的推送
}

var (
	G_GateConnMgr *GateConnMgr
)

func (connMgr *GateConnMgr) doPush(gatewayIdx int, pushJob *PushJob, itemJsons []byte) {
	if pushJob.pushType == common.TYPE_PUSH_ALL {
		connMgr.gateConns[gatewayIdx].PushAll(itemJsons)
	} else if pushJob.pushType == common.TYPE_PUSH_ROOM {
		connMgr.gateConns[gatewayIdx].PushRoom(pushJob.roomId, itemJsons)
	}
	<-connMgr.pendingChan[gatewayIdx]
}

//消息分发协程
func (connMgr *GateConnMgr) dispatchWorkerMain(workerIdx int) {
	var (
		pushJob    *PushJob
		gatewayIdx int
		itemJsons  []byte
		err        error
	)

	for {
		select {
		case pushJob = <-connMgr.dispatchChan:
			//序列化
			if itemJsons, err = json.Marshal(pushJob.message); err != nil {
				continue
			}
			//分发到所有的gateway
			for gatewayIdx = 0; gatewayIdx < len(connMgr.gateConns); gatewayIdx++ {
				select {
				case connMgr.pendingChan[gatewayIdx] <- 1:
					go connMgr.doPush(gatewayIdx, pushJob, itemJsons)
				default:

				}
			}
		}
	}
}

func InitGateConnMgr() (err error) {
	var (
		gatewayIdx        int
		dispatchWorkerIdx int
		gatewayConfig     gatewayConfig
		gateConnMgr       *GateConnMgr
	)

	gateConnMgr = &GateConnMgr{
		gateConns:    make([]*GateConn, len(G_config.GatewayList)),
		pendingChan:  make([]chan byte, len(G_config.GatewayList)),
		dispatchChan: make(chan *PushJob),
	}

	for gatewayIdx, gatewayConfig = range G_config.GatewayList {
		if gateConnMgr.gateConns[gatewayIdx], err = InitGateConn(&gatewayConfig); err != nil {
			return
		}
		gateConnMgr.pendingChan[gatewayIdx] = make(chan byte, G_config.GatewayMaxPendingCount)
	}

	for dispatchWorkerIdx = 0; dispatchWorkerIdx < G_config.GatewayDispatchWorkerCount; dispatchWorkerIdx++ {
		go gateConnMgr.dispatchWorkerMain(dispatchWorkerIdx)
	}

	G_GateConnMgr = gateConnMgr
	return
}

func (connMgr *GateConnMgr) PushAll(msgArr []json.RawMessage) (err error) {
	var (
		pushJob *PushJob
	)

	pushJob = &PushJob{
		pushType: common.TYPE_PUSH_ALL,
		message:  msgArr,
	}
	select {
	case connMgr.dispatchChan <- pushJob:
		DispatchTotal_INCR(int64(len(msgArr)))
	default:
		DispatchFail_INCR(int64(len(msgArr)))
		err = common.ERR_LOGIC_DISPATCH_CHANNEL_FULL
	}
	return
}

func (connMgr *GateConnMgr) PushRoom(room string, msgArr []json.RawMessage) (err error) {
	var (
		pushJob *PushJob
	)

	pushJob = &PushJob{
		pushType: common.TYPE_PUSH_ROOM,
		roomId:   room,
		message:  msgArr,
	}

	select {
	case connMgr.dispatchChan <- pushJob:
		DispatchTotal_INCR(int64(len(msgArr)))
	default:
		DispatchFail_INCR(int64(len(msgArr)))
		err = common.ERR_LOGIC_DISPATCH_CHANNEL_FULL
	}
	return
}
