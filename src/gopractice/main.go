package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopractice/config"
	"gopractice/model"
	"os"
	"io"
	"gopractice/middleware"
)

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

		gin.DefaultWriter = io.MultiWriter(file)
	}

	app := gin.New()

	maxSize := int64(config.ServerConfig.MaxMultipartMemory)

	app.MaxMultipartMemory = maxSize << 20

	app.Use(gin.Logger())

	app.Use(gin.Recovery())

	app.Use(middleware.APIStatsD())


	
}
