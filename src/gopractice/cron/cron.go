package cron

import (
	"github.com/robfig/cron"
	"gopractice/config"
	"gopractice/model"
)

var cronMap = map[string]func(){}
func init() {
	if config.ServerConfig.Env != model.DevelopmentMode {
		cronMap["0 0 3 * * *"] = yesterdayCron
	}
}

func New() *cron.Cron {
	//定时器
	c := cron.New()
	for spec,cmd := range cronMap{
		c.AddFunc(spec,cmd) //每天凌晨3点执行一次
	}
	return c
}
