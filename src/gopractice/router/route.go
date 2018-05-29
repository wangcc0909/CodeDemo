package router

import (
	"github.com/gin-gonic/gin"
	"gopractice/config"
	"gopractice/middleware"
)

func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix, middleware.RefreshTokenCookie)
	{
		api.GET("/siteinfo",)
		api.POST("/login",user.Sin)
	}

}
