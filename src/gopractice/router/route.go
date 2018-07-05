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
	"gopractice/cotroller/collect"
	"gopractice/cotroller/comment"
	"gopractice/cotroller/vote"
)

//这里就是接口
func Route(router *gin.Engine) {
	apiPrefix := config.ServerConfig.APIPrefix

	api := router.Group(apiPrefix, middleware.RefreshTokenCookie) //这个是群组中间件
	{
		api.GET("/siteinfo", common.SiteInfo)
		api.POST("/login", user.Signin)
		api.POST("/signup", user.Signup)
		api.POST("/signout", middleware.SigninRequired, user.Signout) //单个中间件
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

		api.DELETE("/user/career/delete/:id", middleware.SigninRequired, user.DeleteCareer)
		api.DELETE("/user/school/delete/:id", middleware.SigninRequired, user.DeleteSchool)

		api.GET("/messages/unread", middleware.SigninRequired, message.UnRead)
		api.GET("/messages/read/:id", middleware.SigninRequired, message.Read)

		api.GET("/categories", category.List)
		api.GET("/articles", article.List)

		api.GET("/articles/max/bycomment", article.ListMaxComment)
		api.GET("/articles/max/bybrowse", article.ListMaxBrowse)
		api.GET("/articles/top/global", article.Tops)
		api.GET("/articles/info/:id", article.Info)
		api.GET("/articles/user/:userID", article.UserArticleList)

		api.POST("/articles/create", middleware.SigninRequired, article.Create)
		api.POST("/articles/top/:id", middleware.EditorRequired, article.Top)

		api.PUT("/articles/update", middleware.SigninRequired, article.Update)
		api.DELETE("/articles/delete/:id", middleware.SigninRequired, article.Delete)
		api.DELETE("/articles/deletetop/:id", middleware.EditorRequired, article.DeleteTop)

		api.GET("/collects", collect.Collects)
		api.GET("/collects/folders/withsource", middleware.SigninRequired, collect.FoldersWithSource)
		api.GET("/collects/user/:userID/folders", collect.Folders)
		api.POST("/collects/create", middleware.SigninRequired, collect.CreateCollect)
		api.POST("/collects/folder/create", middleware.SigninRequired, collect.CreateFolder)
		api.DELETE("/collects/delete/:id",middleware.SigninRequired,collect.DeleteCollect)

		api.GET("/comments/user/:userID",comment.UserCommentList)
		api.GET("/comments/source/:sourceName/:sourceID",comment.SourceComments)
		api.POST("/comments/create",middleware.SigninRequired,comment.Create)
		api.PUT("/comments/update",middleware.SigninRequired,comment.Update)
		api.DELETE("/comments/delete/:id",middleware.SigninRequired,comment.Delete)

		api.GET("/votes",vote.List)
		api.GET("/votes/info/:id",vote.Info)
		api.GET("/votes/max/bybrowse",vote.ListMaxBrowse)
		api.GET("/votes/max/bycomment",vote.ListMaxComment)
		api.GET("/votes/user/:userID",vote.UserVoteList)
		api.POST("/votes/create",middleware.SigninRequired,vote.Create)
	}
}
