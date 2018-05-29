package router

import (
	"github.com/gin-gonic/gin"
	"gopractice/config"
	"gopractice/middleware"
	"gopractice/cotroller/user"
)

func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix, middleware.RefreshTokenCookie)
	{
		api.GET("/siteinfo",)
		api.POST("/login",user.Signin)
	}

}
