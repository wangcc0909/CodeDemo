package gateway

type Stats struct {
	//反馈在线长连接的数量
	OnlineConnections int64 `json:"onlineConnections"`
	//反馈客户端推送的压力
	SendMessageTotal int64 `json:"sendMessageTotal"`
	SendMessageFail  int64 `json:"sendMessageFail"`
	//反馈connMgr消息分发模块的压力
	DispatchPending int64 `json:"dispatchPending"`
	PushJobPending  int64 `json:"pushJobPending"`
	DispatchFail    int64 `json:"dispatchFail"`

	//返回在线房间总数 有利于分析内存上涨的原因
	RoomCount int64 `json:"roomCount"`
	//Merger模块处理队列 反馈出消息合拼的压力情况
	MergerPending int64 `json:"mergerPending"`
	//Merger模块合并发送消息总数和失败
	MergerRoomTotal int64 `json:"mergerRoomTotal"`
	MergerAllTotal  int64 `json:"mergerAllTotal"`
	MergerRoomFail  int64 `json:"mergerRoomFail"`
	MergerAllFail   int64 `json:"mergerAllFail"`
}
