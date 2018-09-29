package router

import (
	"github.com/gin-gonic/gin"
	"song/controller/crawler"
	"song/controller/song"
)

func Routers(route *gin.Engine) {
	apiPrefix := "/api"
	app := route.Group(apiPrefix)
	{
		app.POST("/crawler",crawler.CrawlerKouGou)
		app.GET("/list",song.List)
	}

}


