package vote

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"strconv"
	"fmt"
	"net/http"
	"time"
	"gopractice/util"
	"github.com/jinzhu/gorm"
)

func List(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var status int
	var hasStatus = false
	var pageNo int
	var pageErr error
	var statusErr error
	var votes []model.Vote

	if pageNo, pageErr = strconv.Atoi(c.Query("pageNo")); pageErr != nil {
		pageNo = 1
	}

	if pageNo < 1 {
		pageNo = 1
	}

	offset := (pageNo - 1) * model.PageSize
	pageSize := model.PageSize

	statusStr := c.Query("status")
	if statusStr == "" {
		hasStatus = false
	} else if status, statusErr = strconv.Atoi(statusStr); statusErr != nil {
		sendErrJson("status不正确", c)
		return
	} else {
		hasStatus = true
	}

	if hasStatus {
		if status != model.VoteUnderway && status != model.VoteOver {
			sendErrJson("status不正确", c)
			return
		}

		if err := model.DB.Where("status = ?", status).Offset(offset).
			Limit(pageSize).Order("created_at DESC").Find(&votes).Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}
	} else {
		if err := model.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&votes).Error; err != nil {
			sendErrJson("error", c)
			return
		}
	}

	for i := 0; i < len(votes); i++ {
		if err := model.DB.Model(&votes[i]).Related(votes[i].User, "users").Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}

		if votes[i].LastUserID != 0 {
			if err := model.DB.Model(&votes[i]).Related(votes[i].LastUser, "users", "last_user_id").Error; err != nil {
				fmt.Println(err.Error())
				sendErrJson("error", c)
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"votes": votes,
		},
	})

}

//info 查询投票
func Info(c *gin.Context) {
	sendErrJson := common.SendErrJson
	id, idErr := strconv.Atoi(c.Param("id"))

	if idErr != nil {
		sendErrJson("无效的ID", c)
		return
	}

	var vote model.Vote

	if err := model.DB.First(&vote, id).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的ID", c)
		return
	}

	vote.BrowseCount++

	if vote.Status != model.VoteOver && vote.EndAt.Unix() < time.Now().Unix() {
		vote.Status = model.VoteOver
	}

	if err := model.DB.Save(&vote).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Model(&vote).Related(&vote.User, "users").Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Model(&vote).Related(&vote.VoteItems, "vote_items").Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Model(&vote).Where("source_name = ?", model.CommentSourceVote).Related(&vote.Comments, "comments").Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	for i := 0; i < len(vote.Comments); i++ {
		if err := model.DB.Model(&vote.Comments[i]).Related(&vote.Comments[i].User, "users").Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}

		vote.Comments[i].HTMLContent = util.MarkdownToHTML(vote.Comments[i].Content)
		parentId := vote.Comments[i].ParentID

		var parents []model.Comment

		//只查看回复的直接父回复

		if parentId != 0 {
			var parent model.Comment
			var parentExist = true

			if err := model.DB.Where("id = ?", parentId).Find(&parent).Error; err != nil {
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
				vote.Comments[i].Parents = parents
			}
		}
	}

	if c.Query("f") != "md" {
		if vote.ContentType == model.ContentTypeMarkdown {
			vote.HTMLContent = util.MarkdownToHTML(vote.Content)
		} else if vote.ContentType == model.ContentTypeHTML {
			vote.HTMLContent = util.AvoidXss(vote.HTMLContent)
		} else {
			vote.HTMLContent = util.MarkdownToHTML(vote.Content)
		}

		vote.Content = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  vote,
	})
}

//访问量最多的投票  返回5条
func ListMaxBrowse(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var votes []model.Vote

	if err := model.DB.Order("browse_count DESC").Limit(5).Find(&votes).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"votes": votes,
		},
	})

}

func ListMaxComment(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var votes []model.Vote

	if err := model.DB.Order("comment_count DESC").Limit(5).Find(&votes).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"votes": votes,
		},
	})

}
