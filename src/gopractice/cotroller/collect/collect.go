package collect

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"strconv"
	"fmt"
	"github.com/jinzhu/gorm"
	"gopractice/util"
	"net/http"
)

// Collects 根据收藏夹查询用户已收藏的话题或投票
func Collects(c *gin.Context) {
	sendErrJson := common.SendErrJson

	var collects []model.Collect

	userID, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		sendErrJson("无效的UserID", c)
		return
	}

	var user model.User
	if err := model.DB.First(&user, userID).Error; err != nil {
		sendErrJson("无效用户ID", c)
		return
	}

	pageNo, pageErr := strconv.Atoi(c.Query("pageNo"))

	if pageErr != nil || pageNo < 1 {
		pageNo = 1
	}

	folderID, folderErr := strconv.Atoi(c.Query("folderID"))

	if folderErr != nil {
		sendErrJson("无效的folderID", c)
		return
	}

	var folder model.Folder
	if err := model.DB.First(&folder, folderID).Error; err != nil {
		sendErrJson("无效的folderID", c)
		return
	}

	var pageSize int
	var pageSizeErr error

	if pageSize, pageSizeErr = strconv.Atoi(c.Query("pageSize")); pageSizeErr != nil {
		sendErrJson("无效的pageSize", c)
		return
	}

	if pageSize < 1 || pageSize > model.MaxPageSize {
		sendErrJson("无效的pageSize", c)
		return
	}

	offset := (pageNo - 1) * pageSize

	if err := model.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).Offset(offset).
		Limit(pageSize).Order("created_at DESC").Find(&collects).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	var totalCount int

	if err := model.DB.Model(&model.Collect{}).Where("folder_id = ? AND user_id = ?", folderID, userID).
		Count(&totalCount).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	var results []map[string]interface{}

	for i := 0; i < len(collects); i++ {
		data := make(map[string]interface{})
		var article model.Article
		var vote model.Vote
		data["id"] = collects[i].ID

		if collects[i].SourceName == model.CollectSourceArticle {
			if err := model.DB.Model(&collects[i]).Related(&article, "articles", "source_id").Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					fmt.Println(err.Error())
					sendErrJson("error", c)
					return
				}
			}
			data["sourceName"] = model.CollectSourceArticle
			data["articleID"] = article.ID
			data["articleName"] = article.Name
			if article.ContentType == model.ContentTypeMarkdown {
				data["htmlContent"] = util.MarkdownToHTML(article.Content)
			} else {
				data["htmlContent"] = util.AvoidXss(article.HTMLContent)
			}
		} else if collects[i].SourceName == model.CollectSourceVote {
			if err := model.DB.Model(&collects[i]).Related(&vote, "votes", "source_id").Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					fmt.Println(err.Error())
					sendErrJson("error", c)
					return
				}
			}
			data["sourceName"] = model.CollectSourceVote
			data["voteID"] = vote.ID
			data["voteName"] = vote.Name
			if vote.ContentType == model.ContentTypeMarkdown {
				data["htmlContent"] = util.MarkdownToHTML(vote.Content)
			} else {
				data["htmlContent"] = util.AvoidXss(vote.HTMLContent)
			}

		}
		results = append(results, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"folderID":   folderID,
			"folderName": folder.Name,
			"collects":   results,
			"pageNo":     pageNo,
			"pageSize":   pageSize,
			"totalCount": totalCount,
		},
	})

}
