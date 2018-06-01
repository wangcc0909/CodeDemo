package router

import (
	"github.com/gin-gonic/gin"
	"gopractice/config"
	"gopractice/middleware"
	"gopractice/cotroller/user"
	"gopractice/cotroller/common"
)

//这里就是接口
func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix, middleware.RefreshTokenCookie)
	{
		api.GET("/siteinfo",common.SiteInfo)
		api.POST("/login",user.Signin)
		api.POST("/signup",user.Signup)
		api.POST("/signout",middleware.SigninRequired,user.Signout)
		api.POST("/upload",middleware.SigninRequired,common.UploadHandler)
		api.POST("crawlnotsavecontent",middleware.EditorRequired,)


	}

}
