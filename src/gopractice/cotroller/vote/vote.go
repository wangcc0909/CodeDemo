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
	"errors"
	"strings"
	"unicode/utf8"
)

func Create(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var voteErr error
	var vote model.Vote
	type ReqData struct {
		Vote      model.Vote       `json:"vote"`
		VoteItems []model.VoteItem `json:"voteItems"`
	}

	var reqData ReqData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		fmt.Println(err.Error())
		sendErrJson("参数无效", c)
		return
	}

	if len(reqData.VoteItems) < 2 {
		sendErrJson("至少要添加两个投票项", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)
	tx := model.DB.Begin()

	if vote, voteErr = save(false, reqData.Vote, user, tx); voteErr != nil {
		tx.Rollback()
		fmt.Println(voteErr.Error())
		sendErrJson("error", c)
		return
	}

	for i := 0; i < len(reqData.VoteItems); i++ {
		var voteItem model.VoteItem
		var err error
		reqData.VoteItems[i].Count = 0
		reqData.VoteItems[i].VoteID = vote.ID
		if voteItem, err = saveVoteItem(reqData.VoteItems[i], tx); err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}
		vote.VoteItems = append(vote.VoteItems, voteItem)
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  vote,
	})
}

func save(isEdit bool, vote model.Vote, user model.User, tx *gorm.DB) (model.Vote, error) {
	var queryVote model.Vote
	if isEdit {
		if err := tx.First(&queryVote, vote.ID).Error; err != nil {
			return vote, errors.New("无效的ID")
		}
	} else {
		vote.UserID = user.ID
	}

	if isEdit {
		if vote.Status == model.VoteOver {
			return vote, errors.New("投票已经结束,不能再进行编辑")
		}
		vote.BrowseCount = queryVote.BrowseCount
		vote.CommentCount = queryVote.CommentCount
		vote.Status = queryVote.Status
		vote.CreatedAt = queryVote.CreatedAt
		vote.UpdateAt = time.Now()
		vote.UserID = queryVote.UserID
		vote.ContentType = queryVote.ContentType
	} else {
		vote.BrowseCount = 0
		vote.CommentCount = 0
		vote.Status = model.VoteUnderway
		vote.CreatedAt = time.Now()
		vote.UpdateAt = vote.CreatedAt
		vote.ContentType = model.ContentTypeMarkdown
	}

	vote.Name = strings.TrimSpace(vote.Name)
	vote.Content = strings.TrimSpace(vote.Content)
	vote.Name = util.AvoidXss(vote.Name)

	if vote.Name == "" {
		return vote, errors.New("名称不能为空")
	}

	if vote.Content == "" || utf8.RuneCountInString(vote.Content) <= 0 {
		return vote, errors.New("内容不能为空")
	}

	if utf8.RuneCountInString(vote.Name) > model.MaxNameLen {
		msg := "名称不能超过" + fmt.Sprintf("%d", model.MaxNameLen) + "个字符"
		return vote, errors.New(msg)
	}

	if vote.CreatedAt.Unix() >= vote.EndAt.Unix() {
		return vote, errors.New("结束时间要大于创建时间")
	}

	if isEdit {
		if err := tx.Save(&vote).Error; err != nil {
			fmt.Println(err.Error())
			return vote, errors.New("error")
		}
	} else {
		if err := tx.Create(&vote).Error; err != nil {
			fmt.Println(err.Error())
			return vote, errors.New("error")
		}
	}

	return vote, nil
}

func saveVoteItem(voteItems model.VoteItem, tx *gorm.DB) (model.VoteItem, error) {
	voteItems.Name = util.AvoidXss(voteItems.Name)
	voteItems.Name = strings.TrimSpace(voteItems.Name)

	if voteItems.Name == "" {
		return voteItems, errors.New("名称不能为空")
	}

	if utf8.RuneCountInString(voteItems.Name) > model.MaxNameLen {
		msg := "名称不能超过" + fmt.Sprintf("%d", model.MaxNameLen) + "个字符"
		return voteItems, errors.New(msg)
	}

	var vote model.Vote

	if err := tx.First(&vote, voteItems.VoteID).Error; err != nil {
		fmt.Println(err.Error())
		return voteItems, errors.New("无效的voteID")
	}

	if vote.Status == model.VoteOver {
		return voteItems, errors.New("投票已结束,不能添加投票项")
	}

	if err := tx.Create(&voteItems).Error; err != nil {
		fmt.Println(err.Error())
		return voteItems, errors.New("error")
	}

	return voteItems, nil
}

//创建投票项
func CreateVoteItem(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var voteItem model.VoteItem
	if err := c.ShouldBindJSON(&voteItem); err != nil {
		fmt.Println(err.Error())
		sendErrJson("参数无效", c)
		return
	}

	var itemErr error
	tx := model.DB.Begin()

	if voteItem, itemErr = saveVoteItem(voteItem, tx); itemErr != nil {
		tx.Rollback()
		fmt.Println(itemErr.Error())
		sendErrJson("error", c)
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  voteItem,
	})
}

//获取投票列表
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

//评论最多的投票 返回5条
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

//某用户的投票列表
func UserVoteList(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var orderType int
	var orderTypeErr error
	var orderStr string
	var isDESC int
	var descErr error
	var pageNo int
	var pageNoErr error
	var pageSize int
	var pageSizeErr error
	userID, idErr := strconv.Atoi(c.Param("userID"))

	if idErr != nil {
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

	//1.按日期排序  2.按点赞数排序  3.按评论数排序
	if orderType != 1 && orderType != 2 && orderType != 3 {
		sendErrJson("无效的orderType", c)
		return
	}

	if isDESC, descErr = strconv.Atoi(c.Query("desc")); descErr != nil {
		sendErrJson("无效的desc", c)
		return
	}

	if isDESC != 0 && isDESC != 1 {
		sendErrJson("无效的desc", c)
		return
	}

	if pageSize, pageSizeErr = strconv.Atoi(c.Query("pageSize")); pageSizeErr != nil {
		sendErrJson("无效的pageSize", c)
		return
	}

	if pageSize < 1 || pageSize > model.MaxPageSize {
		sendErrJson("无效的pageSize", c)
		return
	}

	if pageNo, pageNoErr = strconv.Atoi(c.Query("pageNo")); pageNoErr != nil {
		pageNo = 1
	}

	if pageNo < 1 {
		pageNo = 1
	}

	offset := (pageNo - 1) * pageSize

	if orderType == 1 {
		orderStr = "created_at"
	} else if orderType == 2 {
		orderStr = "up_count"
	} else if orderType == 3 {
		orderStr = "comment_count"
	}

	if isDESC == 1 {
		orderStr += " DESC"
	} else {
		orderStr += " ASC"
	}

	var totalCount int
	var userVotes []model.UserVote

	if err := model.DB.Model(&model.UserVote{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Where("user_id = ?", userID).Order(orderStr).Offset(offset).Limit(pageSize).Find(&userVotes).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	var votes []model.Vote
	for i := 0; i < len(userVotes); i++ {
		if err := model.DB.Model(&userVotes[i]).Related(&userVotes[i].Vote, "votes").Error; err != nil {
			sendErrJson("error", c)
			return
		}
		if userVotes[i].Vote.ContentType == model.ContentTypeMarkdown {
			userVotes[i].Vote.HTMLContent = util.MarkdownToHTML(userVotes[i].Vote.Content)
		} else {
			userVotes[i].Vote.HTMLContent = util.AvoidXss(userVotes[i].Vote.HTMLContent)
		}

		userVotes[i].Vote.Content = ""
		votes = append(votes, userVotes[i].Vote)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"votes":      votes,
			"pageSize":   pageSize,
			"pageNo":     pageNo,
			"totalCount": totalCount,
		},
	})

}

//UserVoteVoteItem  用户投了一票
func UserVoteVoteItem(c *gin.Context) {
	sendErrJson := common.SendErrJson
	id, idErr := strconv.Atoi(c.Param("id"))
	if idErr != nil {
		sendErrJson("无效的ID", c)
		return
	}

	var voteItem model.VoteItem

	if err := model.DB.First(&voteItem, id).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的ID", c)
		return
	}

	var vote model.Vote
	if err := model.DB.Model(&voteItem).Related(&vote).Error; err != nil {
		sendErrJson("无效的ID", c)
		return
	}

	if vote.Status == model.VoteOver {
		sendErrJson("投票已经结束", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)

	var existUserVote model.UserVote

	if err := model.DB.Where("user_id = ? AND vote_id = ?", user.ID, vote.ID).Find(&existUserVote).Error; err == nil {
		sendErrJson("已参与过投票", c)
		return
	}

	voteItem.Count++
	if err := model.DB.Save(&voteItem).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	userVote := model.UserVote{
		UserID:     user.ID,
		VoteID:     voteItem.VoteID,
		VoteItemID: voteItem.ID,
	}

	if err := model.DB.Create(&userVote).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Save(&vote).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

//更新投票
func Update(c *gin.Context) {
	sendErrJson := common.SendErrJson

	var vote model.Vote
	if err := c.ShouldBindJSON(&vote); err != nil {
		fmt.Println(err.Error())
		sendErrJson("参数错误", c)
		return
	}

	iUser, _ := c.Get("user")
	user := iUser.(model.User)
	var voteErr error
	tx := model.DB.Begin()

	if vote, voteErr = save(true, vote, user, tx); voteErr != nil {
		tx.Rollback()
		sendErrJson("error", c)
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  vote,
	})
}

//编辑投票选项
func EditVoteItem(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var voteItem model.VoteItem

	if err := c.ShouldBindJSON(&voteItem); err != nil {
		fmt.Println(err.Error())
		sendErrJson("参数无效", c)
		return
	}

	voteItem.Name = util.AvoidXss(voteItem.Name)
	voteItem.Name = strings.TrimSpace(voteItem.Name)

	if voteItem.Name == "" {
		sendErrJson("名称不能为空", c)
		return
	}

	if utf8.RuneCountInString(voteItem.Name) > model.MaxNameLen {
		msg := "名称不能超过" + fmt.Sprintf("%d", model.MaxNameLen) + "个字符"
		sendErrJson(msg, c)
		return
	}

	var queryVoteItem model.VoteItem

	if err := model.DB.First(&queryVoteItem, voteItem.ID).Error; err != nil {
		sendErrJson("无效的ID", c)
		return
	}

	queryVoteItem.Name = voteItem.Name
	if err := model.DB.Save(&queryVoteItem).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"voteItem": voteItem,
		},
	})
}

//删除投票
func Delete(c *gin.Context) {
	sendErrJson := common.SendErrJson
	//只删除投票本身  用户投票记录保留

	voteId, idErr := strconv.Atoi(c.Param("id"))

	if idErr != nil {
		fmt.Println(idErr.Error())
		sendErrJson("无效的ID", c)
		return
	}

	var vote model.Vote

	if err := model.DB.First(&vote, voteId).Error; err != nil {
		sendErrJson("无效的ID", c)
		return
	}

	tx := model.DB.Begin()

	if err := tx.Delete(&vote).Error; err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Where("DELETE FROM vote_items WHERE vote_id = ?", voteId).Error; err != nil {
		tx.Rollback()
		sendErrJson("error", c)
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"voteID": voteId,
		},
	})
}

func DeleteItem(c *gin.Context) {
	sendErrJson := common.SendErrJson

	voteItemId, idErr := strconv.Atoi(c.Param("id"))

	if idErr != nil {
		sendErrJson("无效的ID", c)
		return
	}
	var voteItem model.VoteItem
	if err := model.DB.First(&voteItem, voteItemId).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Delete(&voteItem).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"voteItemId": voteItemId,
		},
	})
}
