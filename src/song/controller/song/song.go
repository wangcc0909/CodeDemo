package song

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"song/db"
	"song/model"
	"net/http"
)

func List(c *gin.Context) {
	pageNo, err := strconv.Atoi(c.Query("pageNo"))
	if err != nil {
		fmt.Println(err)
		pageNo = 1
	}

	if pageNo < 1 {
		pageNo = 1
	}

	pageSize := 20

	offset := (pageNo - 1) * pageSize
	var songs []model.Song
	if err = db.DB.Model(model.Song{}).Offset(offset).Limit(pageSize).Find(&songs).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "",
			"data": gin.H{},
		})
		return
	}

	var totalCount int
	if err = db.DB.Model(model.Song{}).Count(&totalCount).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "",
			"data": gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"totalCount": totalCount,
			"songList":   songs,
		},
	})
}
