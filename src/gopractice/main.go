package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopractice/config"
	"gopractice/model"
	"os"
	"io"
	"gopractice/middleware"
	"gopractice/router"
	"gopractice/cron"
)

/**
gin 中参数 使用param是将参数放在url中 例如:/user/:id   使用Query是将参数用?隔开 用键值对的形式传递 例如:  ?username=bob&password=123456
 */
func main() {
	fmt.Println("gin version :",gin.Version)

	if config.ServerConfig.Env != model.DevelopmentMode {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()

		file,err :=os.OpenFile(config.ServerConfig.LogFile,os.O_WRONLY|os.O_APPEND|os.O_CREATE,0666)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

		//将打印信息写入到文件中   默认是输出到屏幕
		gin.DefaultWriter = io.MultiWriter(file)
	}

	app := gin.New()

	maxSize := int64(config.ServerConfig.MaxMultipartMemory)

	app.MaxMultipartMemory = maxSize << 20

	//使用中间件
	app.Use(gin.Logger())

	app.Use(gin.Recovery())

	app.Use(middleware.APIStatsD())

	router.Route(app)

	if config.ServerConfig.StatsEnable {
		cron.New().Start()
	}

	app.Run(":" + fmt.Sprintf("%d",config.ServerConfig.Port))
}
