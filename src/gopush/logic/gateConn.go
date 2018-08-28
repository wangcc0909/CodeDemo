package logic

import (
	"net/http"
	"strconv"
	"crypto/tls"
	"time"
	"golang.org/x/net/http2"
	"net/url"
)

type GateConn struct {
	schema string
	client *http.Client //内置长连接 + 并发连接数
}

func InitGateConn(config *gatewayConfig) (gateConn *GateConn, err error) {
	var (
		transport *http.Transport
	)
	gateConn = &GateConn{
		schema: "https://" + config.Hostname + ":" + strconv.Itoa(config.Port),
	}

	transport = &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, //不校验服务器端证书
		MaxIdleConns:        G_config.GatewayMaxConnection,
		MaxIdleConnsPerHost: G_config.GatewayMaxConnection,
		IdleConnTimeout:     time.Duration(G_config.GatewayIdleTimeout) * time.Second, //空闲超时时间
	}

	//启动http/2协议
	http2.ConfigureTransport(transport)

	//http/2 客户端
	gateConn.client = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(G_config.GatewayTimeout) * time.Second,
	}
	return
}

//出于性能考虑  消息数组在此之前已经编码成json
func (gateConn *GateConn) PushAll(itemJson []byte) (err error) {
	var (
		apiUrl   string
		from     url.Values
		retry    int
		response *http.Response
	)
	apiUrl = gateConn.schema + "/push/all"
	from = url.Values{}
	from.Set("items", string(itemJson))
	for retry = 0; retry < G_config.GatewayPushRetry; retry++ {
		if response, err = gateConn.client.PostForm(apiUrl, from); err != nil {
			PushFail_INCR()
			continue
		}
		response.Body.Close()
		break
	}
	return
}

//出于性能考虑  消息数组在此之前已经编码成json
func (gateConn *GateConn) PushRoom(roomId string, itemJson []byte) (err error) {
	var (
		apiUrl   string
		from     url.Values
		retry    int
		response *http.Response
	)
	apiUrl = gateConn.schema + "/push/room"
	from = url.Values{}
	from.Set("room", roomId)
	from.Set("items", string(itemJson))
	for retry = 0; retry < G_config.GatewayPushRetry; retry++ {
		if response, err = gateConn.client.PostForm(apiUrl, from); err != nil {
			PushFail_INCR()
			continue
		}
		response.Body.Close()
		break
	}
	return
}
