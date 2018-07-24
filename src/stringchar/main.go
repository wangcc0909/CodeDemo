package main

import (
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {

	c := cron.New()
	c.AddFunc("8,12,30 * * * * *", func() {
		log.Println("我是定时任务")
	})
	c.Start()

	time.Sleep(time.Second * 60)
	log.Println("结束")
}
