package comment

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"strconv"
	"gopractice/model"
	"fmt"
	"gopractice/util"
	"github.com/jinzhu/gorm"
	"net/http"
	"github.com/gomodule/redigo/redis"
	"strings"
	"unicode/utf8"
	"time"
)

//查询用户的评论
func UserCommentList(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var userID int
	var idErr error
	var orderType int
	var orderTypeErr error
	var orderStr string
	var isDesc int
	var descErr error
	var pageSize int
	var pageSizeErr error
	var pageNo int
	var pageNoErr error

	if userID, idErr = strconv.Atoi(c.Param("userID")); idErr != nil {
		sendErrJson("无效的userID", c)
		return
	}

	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		sendErrJson("无效的userID", c)
		return
	}

	if orderType, orderTypeErr = strconv.Atoi(c.Query("orderType")); orderTypeErr != nil {
		sendErrJson("无效的orderType", c)
		return
	}

	//1.按日期排序  2.按点赞数排序
	if orderType != 1 && orderType != 2 {
		sendErrJson("无效的orderType", c)
		return
	}

	if isDesc, descErr = strconv.Atoi(c.Query("desc")); descErr != nil {
		sendErrJson("无效的desc", c)
		return
	}

	if isDesc != 0 && isDesc != 1 {
		sendErrJson("无效的desc", c)
		return
	}

	if pageNo, pageNoErr = strconv.Atoi(c.Query("pageNo")); pageNoErr != nil {
		pageNo = 1
		pageNoErr = nil
	}

	if pageNo < 1 {
		pageNo = 1
	}

	if pageSize, pageSizeErr = strconv.Atoi(c.Query("pageSize")); pageSizeErr != nil {
		sendErrJson("无效的pageSize", c)
		return
	}

	if pageSize < 1 || pageSize > model.MaxPageSize {
		sendErrJson("无效的pageSize", c)
		return
	}

	offset := (pageNo - 1) * pageSize

	if orderType == 1 {
		orderStr = "created_at"
	} else if orderType == 2 {
		orderStr = "up_count"
	}

	if isDesc == 1 {
		orderStr += " DESC"
	} else {
		orderStr += " ASC"
	}

	var comments []model.Comment
	var totalCount int

	if err := model.DB.Model(&model.Comment{}).Where("user_id = ? AND status != ?", userID, model.CommentVertifyFail).
		Count(&totalCount).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Where("user_id = ? AND status != ?", userID, model.CommentVertifyFail).Order(orderStr).
		Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	var results []map[string]interface{}

	for i := 0; i < len(comments); i++ {
		data := make(map[string]interface{})

		var article model.Article
		var vote model.Vote

		data["id"] = comments[i].ID
		if comments[i].ContentType == model.ContentTypeMarkdown {
			data["html"] = util.MarkdownToHTML(comments[i].Content)
		} else {
			data["html"] = util.AvoidXss(comments[i].HTMLContent)
		}

		if comments[i].SourceName == model.CommentSourceArticle {
			if err := model.DB.Model(&comments[i]).Related(&article, "articles", "source_id").Error; err != nil {
				//没有找到话题  已经删除了
				if err != gorm.ErrRecordNotFound {
					fmt.Println(err.Error())
					sendErrJson("error", c)
					return
				}
			}
			data["sourceName"] = model.CommentSourceArticle
			data["articleID"] = article.ID
			data["articleName"] = article.Name
		} else if comments[i].SourceName == model.CommentSourceVote {
			if err := model.DB.Model(&comments[i]).Related(&vote, "votes", "source_id").Error; err != nil {
				//没有找到话题  已经删除了
				if err != gorm.ErrRecordNotFound {
					fmt.Println(err.Error())
					sendErrJson("error", c)
					return
				}
			}
			data["sourceName"] = model.CommentSourceArticle
			data["voteID"] = vote.ID
			data["voteName"] = vote.Name
		}

		if err := model.DB.Model(comments[i]).Related(&comments[i].User, "users").Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"comments":   results,
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
		},
	})
}

//查询话题或投票评论
func SourceComments(c *gin.Context) {
	sendErrJson := common.SendErrJson
	sourceName := c.Param("sourceName")
	if sourceName != model.CommentSourceArticle && sourceName != model.CommentSourceVote {
		sendErrJson("无效的sourceName", c)
		return
	}

	sourceID, idErr := strconv.Atoi(c.Param("sourceID"))

	if idErr != nil {
		sendErrJson("无效的sourceID", c)
		return
	}

	var article model.Article
	var vote model.Vote

	if sourceName == model.CommentSourceArticle {
		if err := model.DB.First(&article, sourceID).Error; err != nil {
			sendErrJson("无效的sourceID", c)
			return
		}
	}

	if sourceName == model.CommentSourceVote {
		if err := model.DB.First(&vote, sourceID).Error; err != nil {
			sendErrJson("无效的sourceID", c)
			return
		}
	}

	var comments []model.Comment

	if err := model.DB.Where("source_id = ? AND source_name = ?", sourceID, sourceName).
		Preload("User").Find(&comments).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	for i := 0; i < len(comments); i++ {
		comments[i].HTMLContent = util.MarkdownToHTML(comments[i].Content)
		//只查看直接父回复
		var parentID = comments[i].ParentID

		var parents []model.Comment
		if parentID != 0 {
			var parent model.Comment
			var parentExist = true
			if err := model.DB.Where("id = ?", parentID).Find(&parent).Error; err != nil {
				parentExist = false
				if err != gorm.ErrRecordNotFound {
					fmt.Println(err.Error())
					sendErrJson("error", c)
					return
				}
			}

			if parentExist {
				if err := model.DB.Model(&parent).Related(&parent.User, "users").Error; err != nil {
					fmt.Println(err.Error())
					sendErrJson("error", c)
					return
				}

				parents = append(parents, parent)
				comments[i].Parents = parents
			}

		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"comments": comments,
		},
	})

}

//创建评论
func Create(c *gin.Context) {
	sendErrJson := common.SendErrJson
	iUser, _ := c.Get("user")

	user := iUser.(model.User)

	RedisConn := model.RedisPool.Get()

	defer RedisConn.Close()

	minuteKey := model.CommentMinuteLimit + fmt.Sprintf("%d", user.ID)

	minuteCount, minuteErr := redis.Int64(RedisConn.Do("GET", minuteKey))
	if minuteErr == nil && minuteCount > model.CommentMinuteLimitCount {
		sendErrJson("您的操作过于频繁,请先休息一会", c)
		return
	}

	minuteRemainingTime, _ := redis.Int64(RedisConn.Do("TTL", minuteKey))

	if minuteRemainingTime < 0 || minuteRemainingTime > 60 {
		minuteRemainingTime = 60
	}

	if _, err := RedisConn.Do("SET", minuteKey, minuteCount+1, "EX", minuteRemainingTime); err != nil {
		fmt.Println("redis set failed err:", err)
		sendErrJson("内部错误", c)
		return
	}

	dayKey := model.CommentDayLimit + fmt.Sprintf("%d", user.ID)

	dayCount, dayErr := redis.Int64(RedisConn.Do("GET", dayKey))

	if dayErr == nil && dayCount > model.CommentMinuteLimitCount {
		sendErrJson("您今天的操作过于频繁,请先休息一会", c)
		return
	}

	dayRemainingTime, _ := redis.Int64(RedisConn.Do("TTL", dayKey))

	secondOfDay := int64(24 * 60 * 60)
	if dayRemainingTime < 0 || dayRemainingTime > secondOfDay {
		dayRemainingTime = secondOfDay
	}

	if _, err := RedisConn.Do("SET", dayKey, dayCount+1, "EX", dayRemainingTime); err != nil {
		fmt.Println("redis set failed err:", err)
		sendErrJson("内部错误", c)
		return
	}

	Save(c, false)
}

//保存评论
/**
先更新用户的信息再更新文章的信息
 */
func Save(c *gin.Context, isEdit bool) {
	sendErrJson := common.SendErrJson
	var comment model.Comment
	var parentComment model.Comment

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	if user.Role == model.UserRoleCrawler {
		sendErrJson("爬虫管理员不能回复", c)
		return
	}

	//编辑评论时只传id 和 content
	if err := c.ShouldBindJSON(&comment); err != nil {
		sendErrJson("参数无效", c)
		return
	}

	var article model.Article
	var vote model.Vote

	//不是重新编辑
	if !isEdit {
		if comment.SourceName != model.CommentSourceArticle && comment.SourceName != model.CommentSourceVote {
			sendErrJson("无效的sourceName", c)
			return
		}

		if comment.SourceName == model.CommentSourceArticle {
			if err := model.DB.First(&article, comment.SourceID).Error; err != nil {
				sendErrJson("无效的sourceName", c)
				return
			}
		}

		if comment.SourceName == model.CommentSourceVote {
			if err := model.DB.First(&vote, comment.SourceID).Error; err != nil {
				sendErrJson("无效的sourceName", c)
				return
			}
		}

		if comment.ParentID != model.NoParent {
			if err := model.DB.First(&parentComment, comment.ParentID).Error; err != nil {
				sendErrJson("无效的parentID", c)
				return
			}

			if parentComment.SourceID != comment.SourceID {
				sendErrJson("无效的parentID", c)
				return
			}
		}
	}

	comment.Content = strings.TrimSpace(comment.Content)

	if comment.Content == "" {
		sendErrJson("评论不能为空", c)
		return
	}

	if utf8.RuneCountInString(comment.Content) > model.MaxCommentLen {
		msg := "评论不能超过" + fmt.Sprintf("%d", model.MaxCommentLen) + "个字符"
		sendErrJson(msg, c)
		return
	}

	comment.Status = model.CommentVertifying //设置为校验评论的合法性
	comment.UserID = user.ID

	var updateComment model.Comment

	if !isEdit {
		comment.ContentType = model.ContentTypeMarkdown

		tx := model.DB.Begin()

		if err := tx.Create(&comment).Error; err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			sendErrJson("error", c)
			return
		}

		updateUserMap := map[string]interface{}{
			"comment_count": user.CommentCount + 1,
			"score":         user.Score + model.CommentScore,
		}

		if err := tx.Model(&user).Updates(updateUserMap).Error; err != nil {
			fmt.Println(err.Error())
			tx.Rollback()
			sendErrJson("error", c)
			return
		}

		if err := model.UserToRedis(user); err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}

		var author model.User //文章作者

		if comment.SourceName == model.CommentSourceArticle {
			author.ID = article.UserID
			articleMap := map[string]interface{}{
				"comment_count":   article.CommentCount + 1,
				"last_user_id":    user.ID,
				"last_comment_at": time.Now(),
			}

			if err := tx.Model(&article).Updates(articleMap).Error; err != nil {
				fmt.Println(err.Error())
				tx.Rollback()
				sendErrJson("error", c)
				return
			}
		} else if comment.SourceName == model.CommentSourceVote {
			author.ID = vote.UserID
			voteMap := map[string]interface{}{
				"comment_count":   vote.CommentCount + 1,
				"last_user_id":    user.ID,
				"last_comment_at": time.Now(),
			}

			if err := tx.Model(&vote).Updates(voteMap).Error; err != nil {
				fmt.Println(err.Error())
				tx.Rollback()
				sendErrJson("error", c)
				return
			}
		}

		//自己评论自己的不增加积分
		if user.ID != author.ID {
			if err := tx.First(&author, author.ID).Error; err != nil {
				fmt.Println(err.Error())
				tx.Rollback()
				sendErrJson("error", c)
				return
			}

			authorScore := author.Score + model.ByCollectScore
			if err := tx.Model(&author).Update("score", authorScore).Error; err != nil {
				fmt.Println(err.Error())
				tx.Rollback()
				sendErrJson("error", c)
				return
			}
		}

		//回复别人的话题时给消息提示
		//对回复进行回复,即使回复属于自己创建的也给父回复发送消息
		if user.ID != author.ID || comment.ParentID != model.NoParent {
			var message model.Message
			message.FromUserId = user.ID
			message.SourceId = comment.SourceID
			message.SourceName = comment.SourceName
			message.CommentId = comment.ID
			message.Readed = false
			if comment.ParentID != model.NoParent {
				message.Type = model.MessageTypeCommentComment
				message.ToUserId = parentComment.UserID
			} else if comment.SourceName == model.CommentSourceArticle {
				message.Type = model.MessageTypeCommentArticle
				message.ToUserId = author.ID
			} else if comment.SourceName == model.CommentSourceVote {
				message.Type = model.MessageTypeCommentVote
				message.ToUserId = author.ID
			}

			if err := model.DB.Create(&message).Error; err != nil {
				fmt.Println(err.Error())
				tx.Rollback()
				sendErrJson("error", c)
				return
			}
		}
		tx.Commit()
	} else {
		if err := model.DB.First(&updateComment, comment.ID).Error; err != nil {
			sendErrJson("无效的ID", c)
			return
		}

		if user.ID != updateComment.UserID {
			sendErrJson("您无权执行此操作", c)
			return
		}

		updateCommentMap := map[string]interface{}{
			"content": comment.Content,
			"status":  model.CommentVertifying,
		}

		if err := model.DB.Model(&updateComment).Updates(updateCommentMap).Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}
	}

	var commentJson model.Comment

	if isEdit {
		commentJson = updateComment
	} else {
		commentJson = comment
	}

	commentJson.User = user

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"comment": commentJson,
		},
	})

}
