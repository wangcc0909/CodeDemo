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
