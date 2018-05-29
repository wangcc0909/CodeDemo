package model

import (
	"gopractice/config"
	"github.com/cactus/go-statsd-client/statsd"
	"time"
	"fmt"
)

var StatterClient *statsd.Statter

func init() {
	if config.StatsDConfig.URL == "" {
		return
	}

	statter,err := statsd.NewBufferedClient(config.StatsDConfig.URL,config.StatsDConfig.Prefix,300 * time.Millisecond,512)
	if err != nil {
		fmt.Println(err.Error())
	}
	StatterClient = &statter
}
