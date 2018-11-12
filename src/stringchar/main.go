package main

import (
	"github.com/robfig/cron"
	"log"
)

func main() {
//https://www.cnblogs.com/zuxingyu/p/6023919.html
	c := cron.New()
	c.AddFunc("8,12,30 * * * * *", func() {
		log.Println("我是定时任务")
	})
	c.Start()

	for {

	}
	log.Println("结束")
}
