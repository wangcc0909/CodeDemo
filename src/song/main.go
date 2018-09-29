package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"fmt"
	"io"
	"song/config"
	"song/router"
)

func main() {
	fmt.Println("gin version :",gin.Version)

	if config.ServerConfig.Env != "development" {
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
	router.Routers(app)

	/*if config.ServerConfig.StatsEnable {
		cron.New().Start()
	}*/

	app.Run(":" + fmt.Sprintf("%d",config.ServerConfig.Port))
}
