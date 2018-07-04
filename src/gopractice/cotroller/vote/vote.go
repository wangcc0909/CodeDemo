package vote

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"strconv"
	"fmt"
	"net/http"
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
