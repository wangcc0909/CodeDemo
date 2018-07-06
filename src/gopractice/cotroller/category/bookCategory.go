package category

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"fmt"
	"net/http"
)

//图书分类列表
func BookCategoryList(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var categories []model.Category

	if err := model.DB.Order("sequence asc").Find(&categories).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"categories": categories,
		},
	})

}
