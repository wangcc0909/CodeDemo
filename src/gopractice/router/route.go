package router

import (
	"github.com/gin-gonic/gin"
	"gopractice/config"
	"gopractice/middleware"
	"gopractice/cotroller/user"
	"gopractice/cotroller/common"
	"gopractice/cotroller/crawler"
	"gopractice/cotroller/message"
	"gopractice/cotroller/category"
	"gopractice/cotroller/article"
)

//这里就是接口
func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix, middleware.RefreshTokenCookie)
	{
		api.GET("/siteinfo", common.SiteInfo)
		api.POST("/login", user.Signin)
		api.POST("/signup", user.Signup)
		api.POST("/signout", middleware.SigninRequired, user.Signout)
		api.POST("/upload", middleware.SigninRequired, common.UploadHandler)
		api.POST("crawlnotsavecontent", middleware.EditorRequired, crawler.CrawlNotSaveContent)
		api.POST("/active/sendmail", user.ActiveSendMail)

		api.POST("/active/user/:id/:secret", user.ActiveAccount)

		api.POST("/reset/sendemail", user.ResetPasswordMail)
		api.GET("/reset/verify/:id/:secret", user.VerifyResetPasswordLink)
		api.POST("/reset/password/:id/:secret", user.ResetPassword)

		api.GET("/user/info", middleware.SigninRequired, user.SecretInfo)
		api.GET("/user/score/top10", user.Top10)
		api.GET("/user/score/top100", user.Top100)
		api.GET("/user/info/detail", middleware.SigninRequired, user.InfoDetail)
		api.GET("/user/info/public/:id", user.PublicInfo)

		api.POST("/user/uploadavatar", middleware.SigninRequired, user.UploadAvatar)
		api.POST("/user/career/add", middleware.SigninRequired, user.AddCareer)
		api.POST("/user/school/add", middleware.SigninRequired, user.AddSchool)
		api.PUT("/user/update/:filed", middleware.SigninRequired, user.UpdateInfo)
		api.PUT("/user/password/update", middleware.SigninRequired, user.UpdatePassword)

		api.DELETE("/user/career/delete/:id",middleware.SigninRequired,user.DeleteCareer)
		api.DELETE("/user/school/delete/:id",middleware.SigninRequired,user.DeleteSchool)

		api.GET("/messages/unread",middleware.SigninRequired,message.UnRead)
		api.GET("/messages/read/:id",middleware.SigninRequired,message.Read)

		api.GET("/categories",category.List)
		api.GET("/articles",article.List)


	}
}
