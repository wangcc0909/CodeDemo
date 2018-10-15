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
		app.GET("/list",song.List)
		app.GET("/crawler",crawler.CrawlerKouGouWithoutUrl)
	}

}


