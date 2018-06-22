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

//收藏文章或收藏投票
func CreateCollect(c *gin.Context) {
	sendErrJson := common.SendErrJson

	var article model.Article
	var collect model.Collect
	var vote model.Vote

	if err := c.ShouldBindJSON(&collect); err != nil {
		fmt.Println(err.Error())
		sendErrJson("参数无效", c)
		return
	}

	if collect.SourceName != model.CollectSourceArticle && collect.SourceName != model.CollectSourceVote {
		sendErrJson("sourceName无效", c)
		return
	}

	if collect.SourceName == model.CollectSourceArticle {
		if err := model.DB.First(&article, collect.SourceID).Error; err != nil {
			sendErrJson("sourceID无效", c)
			return
		}
	}

	if collect.SourceName == model.CollectSourceVote {
		if err := model.DB.First(&vote, collect.SourceID).Error; err != nil {
			sendErrJson("sourceID无效", c)
			return
		}
	}

	if err := model.DB.First(&collect.Folder, collect.FolderID).Error; err != nil {
		sendErrJson("FolderID无效", c)
		return
	}

	var theCollect model.Collect
	if err := model.DB.Where("source_id = ? AND source_name = ?", collect.SourceID, collect.SourceName).
		First(&theCollect).Error; err != nil {
		sendErrJson("之前已经收藏过", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)
	collect.UserID = user.ID

	//开始一个事务
	tx := model.DB.Begin()

	if err := tx.Save(&collect).Error; err != nil {
		fmt.Println(err.Error())
		//回滚事务
		tx.Rollback()
		sendErrJson("error", c)
		return
	}

	if err := tx.Model(&user).Update("collect_count", user.CollectCount+1).Error; err != nil {
		fmt.Println(err.Error())
		//回滚事务
		tx.Rollback()
		sendErrJson("error", c)
		return
	}

	if model.UserToRedis(user) != nil {
		sendErrJson("error", c)
		return
	}

	if collect.SourceName == model.CollectSourceArticle {
		if err := tx.Model(&article).Update("collect_count", article.CollectCount+1).Error; err != nil {
			fmt.Println(err.Error())
			//回滚事务
			tx.Rollback()
			sendErrJson("error", c)
			return
		}

		//获取相关的user
		if err := tx.Model(&article).Related(&article.User).Error; err != nil {
			tx.Rollback()
			sendErrJson("error", c)
			return
		}

		//自己收藏自己的话题不增加积分
		if article.User.ID != user.ID {
			if err := tx.Model(&article.User).Update("score", article.User.Score+model.ByCollectScore).Error; err != nil {
				fmt.Println(err.Error())
				sendErrJson("err", c)
				tx.Rollback()
				return
			}
		}
	}

	if collect.SourceName == model.CollectSourceVote {
		if err := tx.Model(&vote).Update("collect_count", vote.CollectCount+1).Error; err != nil {
			fmt.Println(err.Error())
			//回滚事务
			tx.Rollback()
			sendErrJson("error", c)
			return
		}

		//获取相关的user
		if err := tx.Model(&vote).Related(&vote.User).Error; err != nil {
			tx.Rollback()
			sendErrJson("error", c)
			return
		}

		//自己收藏自己的话题不增加积分
		if vote.User.ID != user.ID {
			if err := tx.Model(&vote.User).Update("score", vote.User.Score+model.ByCollectScore).Error; err != nil {
				fmt.Println(err.Error())
				sendErrJson("err", c)
				tx.Rollback()
				return
			}
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  collect,
	})
}

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

//查询用户的收藏夹列表
func Folders(c *gin.Context) {
	senErrJson := common.SendErrJson
	userID, userIDErr := strconv.Atoi(c.Param("userID"))
	if userIDErr != nil {
		senErrJson("无效的userID", c)
		return
	}

	var folders []model.Folder
	var foldersErr error
	if folders, foldersErr = queryFolders(userID); foldersErr != nil {
		senErrJson("无效的userID", c)
		return
	}

	var results []map[string]interface{}
	for i := 0; i < len(folders); i++ {
		var data = map[string]interface{}{
			"id":        folders[i].ID,
			"createdAt": folders[i].CreatedAt,
			"updateAt":  folders[i].UpdateAt,
			"deleteAt":  folders[i].DeleteAt,
			"name":      folders[i].Name,
			"userID":    folders[i].UserID,
			"parentID":  folders[i].ParentID,
		}

		var collectCount uint
		if err := model.DB.Model(&model.Collect{}).Where("folder_id = ?", folders[i].ID).Count(&collectCount).Error; err != nil {
			senErrJson("error", c)
			return
		}

		data["collectCount"] = collectCount
		results = append(results, data)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"folders": results,
		},
	})

}

func queryFolders(userId int) ([]model.Folder, error) {
	var user model.User
	if err := model.DB.First(&user, userId).Error; err != nil {
		return nil, err
	}

	var folders []model.Folder
	if err := model.DB.Where("user_id = ?", userId).Order(&folders).Error; err != nil {
		return nil, err
	}

	return folders, nil
}

// FoldersWithSource 查询用户的收藏夹列表，并且返回每个收藏夹中收藏了哪些话题或投票
func FoldersWithSource(c *gin.Context) {
	sendErrJson := common.SendErrJson
	iUser, isExist := c.Get("user")

	if !isExist {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.SUCCESS,
			"msg":   "success",
			"data": gin.H{
				"folders": make([]interface{}, 0),
			},
		})
		return
	}

	user := iUser.(model.User)
	var folders []model.Folder
	var queryFloderErr error

	if folders, queryFloderErr = queryFolders(int(user.ID)); queryFloderErr != nil {
		fmt.Println(queryFloderErr.Error())
		sendErrJson("error", c)
		return
	}

	var results []interface{}

	for i := 0; i < len(folders); i++ {
		var collects []model.Collect
		if err := model.DB.Where("folder_id = ?", folders[i].ID).Find(&collects).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				sendErrJson("error", c)
				return
			}
		}
		results = append(results, gin.H{
			"id":       folders[i].ID,
			"name":     folders[i].Name,
			"collects": collects,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"folders": results,
		},
	})
}
