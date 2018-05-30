package router

import (
	"github.com/gin-gonic/gin"
	"gopractice/config"
	"gopractice/middleware"
	"gopractice/cotroller/user"
	"gopractice/cotroller/common"
)

func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix, middleware.RefreshTokenCookie)
	{
		api.GET("/siteinfo",common.SiteInfo)
		api.POST("/login",user.Signin)
		api.POST("/signup",user.Signup)

	}

}
