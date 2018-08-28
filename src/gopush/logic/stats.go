package logic

import (
	"sync/atomic"
	"encoding/json"
)

type Stats struct {
	DispatchTotal int64 `json:"dispatchTotal"`
	DispatchFail int64 `json:"dispatchFail"`
	PushFail int64 `json:"pushFail"`
}

var (
	G_stats *Stats
)

func InitStats() (err error) {
	G_stats = &Stats{}
	return
}

func DispatchTotal_INCR(batchSize int64) {
	atomic.AddInt64(&G_stats.DispatchTotal,batchSize)
}

func DispatchFail_INCR(batchSize int64) {
	atomic.AddInt64(&G_stats.DispatchFail,batchSize)
}

func PushFail_INCR() {
	atomic.AddInt64(&G_stats.PushFail,1)
}

func (stats *Stats) Dump() (data []byte,err error) {
	return json.Marshal(G_stats)
}
