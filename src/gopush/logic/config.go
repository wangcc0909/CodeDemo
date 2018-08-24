package logic

import (
	"io/ioutil"
	"encoding/json"
)

type gatewayConfig struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

type Config struct {
	ServicePort                int             `json:"servicePort"`
	ServiceReadTimeout         int             `json:"serviceReadTimeout"`
	ServiceWriteTimeout        int             `json:"serviceWriteTimeout"`
	GatewayList                []gatewayConfig `json:"gatewayList"`
	GatewayMaxConnection       int             `json:"gatewayMaxConnection"`
	GatewayTimeout             int             `json:"gatewayTimeout"`
	GatewayIdleTimeout         int             `json:"gatewayIdleTimeout"`
	GatewayDispatchWorkerCount int             `json:"gatewayDispatchWorkerCount"`
	GatewayDispatchChannelSize int             `json:"gatewayDispatchChannelSize"`
	GatewayMaxPendingCount     int             `json:"gatewayMaxPendingCount"`
	GatewayPushRetry           int             `json:"gatewayPushRetry"`
}

var (
	G_config *Config
)

func InitConfig(confFile string) (err error) {
	var (
		data []byte
		conf Config
	)

	data,err = ioutil.ReadFile(confFile)
	if err != nil {
		return
	}
	if err = json.Unmarshal(data,&conf);err != nil {
		return
	}
	G_config = &conf
	return
}
