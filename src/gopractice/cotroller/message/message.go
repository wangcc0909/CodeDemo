package message

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"strconv"
	"gopractice/util"
	"net/http"
)

//未读消息
func UnRead(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var messages []model.Message

	pageNo, pageErr := strconv.Atoi(c.Query("pageNo"))

	if pageErr != nil || pageNo < 1 {
		pageNo = 1
	}

	pageSize, sizeErr := strconv.Atoi(c.Query("pageSize"))
	if sizeErr != nil || pageSize < 1 {
		pageSize = 10
	}

	if pageSize > model.MaxPageSize {
		pageSize = model.MaxPageSize
	}

	iUser, _ := c.Get("user")

	user := iUser.(model.User)

	if model.DB.Where("readed = ? AND to_user_id = ?", 0, user.ID).Order("create_at DESC").
		Offset((pageNo - 1) * pageSize).Limit(pageSize).Find(&messages).Error != nil {
		sendErrJson("error", c)
		return
	}

	var count int
	if err := model.DB.Model(&model.Message{}).Where("readed = ? AND to_user_id = ?", 0, user.ID).Count(&count).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	for i := len(messages) - 1; i >= 0; i-- {
		if err := model.DB.Model(&messages[i]).Related(&messages[i].FromUser, "user", "from_user_id").Error; err != nil {
			sendErrJson("error", c)
			return
		}

		if messages[i].FromUser.Status == model.UserStatusFrozen {
			// 如果用户被冻结了，那么他的操作所产生的对其他用户的消息提示就没了
			messages = append(messages[:i], messages[i+1:]...)
			continue
		}

		if messages[i].Type == model.MessageTypeCommentArticle {
			var article model.Article
			var isSkip = false
			if err := model.DB.Select("name,status").Where("id = ?", messages[i].SourceId).First(&article).Error; err != nil {
				isSkip = true
			}

			if article.Status == model.ArticleVerifyFail {
				isSkip = true
			}

			//出错了,或话题被删除了,或话题未审核通过,就不返回这个消息提示
			if isSkip {
				messages = append(messages[:i], messages[i+1:]...)
				continue
			}
			messages[i].Data.Title = article.Name
		} else if messages[i].Type == model.MessageTypeCommentVote {
			var vote model.Vote
			if err := model.DB.Select("name").Where("id = ?", messages[i].SourceId).First(&vote).Error; err != nil {
				//出错了，或投票被删除了，就不返回这个消息提示
				messages = append(messages[:i], messages[i+1:]...)
				continue
			}
			messages[i].Data.Title = vote.Name
		} else if messages[i].Type == model.MessageTypeCommentComment {
			var comment model.Comment

			if err := model.DB.Where("id = ?", messages[i].CommentId).First(&comment).Error; err != nil {
				//出错了，或评论被删除了，就不返回这个消息提示
				messages = append(messages[:i], messages[i+1:]...)
				continue
			}

			if comment.SourceName == model.CommentSourceArticle {
				var article model.Article
				var isSkip = false
				if err := model.DB.Select("id,status").Where("id = ?", comment.SourceID).First(&article).Error; err != nil {
					isSkip = true
				}
				if article.Status == model.ArticleVerifyFail {
					isSkip = true
				}

				//出错了,或话题被删除了,或话题未审核通过,就不返回这个消息提示
				if isSkip {
					messages = append(messages[:i], messages[i+1:]...)
					continue
				}
			} else if comment.SourceName == model.CommentSourceVote {
				var vote model.Vote
				if err := model.DB.Select("id").Where("id = ?", comment.SourceID).First(&vote).Error; err != nil {
					//出错了，或投票被删除了，就不返回这个消息提示
					messages = append(messages[:i], messages[i+1:]...)
					continue
				}
			}

			var commentContent string

			if comment.ContentType == model.ContentTypeMarkdown {
				commentContent = util.MarkdownToHTML(comment.Content)
			} else {
				commentContent = util.AvoidXss(comment.Content)
			}
			messages[i].Data.CommentContent = commentContent
		}

		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.SUCCESS,
			"msg":   "success",
			"data": gin.H{
				"message":  messages,
				"pageNo":   pageNo,
				"pageSize": pageSize,
				"count":    count,
			},
		})
	}
}

func Read(c *gin.Context) {
	sendErrJson := common.SendErrJson

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Model(&model.Message{}).Where("id = ?", id).Update("readed", true).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}
