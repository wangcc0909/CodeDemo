package gateway

import (
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	WsPort               int    `json:"wsPort"`
	WsReadTimeout        int    `json:"wsReadTimeout"`
	WsWriteTimeout       int    `json:"wsWriteTimeout"`
	WsInChannelSize      int    `json:"wsInChannelSize"`
	WsOutChannelSize     int    `json:"wsOutChannelSize"`
	WsHeartbeatInternal  int    `json:"wsHeartbeatInternal"`
	MaxMergerDelay       int    `json:"maxMergerDelay"`
	MaxMergerBatchSize   int    `json:"maxMergerBatchSize"`
	MergerWorkerCount    int    `json:"mergerWorkerCount"`
	MaxChannelSize       int    `json:"maxChannelSize"`
	ServerPort           int    `json:"serverPort"`
	ServerReadTimeout    int    `json:"serverReadTimeout"`
	ServerWriteTimeout   int    `json:"serverWriteTimeout"`
	ServerPem            string `json:"serverPem"`
	ServerKey            string `json:"serverKey"`
	BucketCount          int    `json:"bucketCount"`
	BucketWorkerCount    int    `json:"bucketWorkerCount"`
	MaxJoinRoom          int    `json:"maxJoinRoom"`
	DispatchChannelSize  int    `json:"dispatchChannelSize"`
	DispatchWorkerCount  int    `json:"dispatchWorkerCount"`
	BucketJobChannelSize int    `json:"bucketJobChannelSize"`
	BucketJobWorkerCount int    `json:"bucketJobWorkerCount"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_config = &conf
	return
}
