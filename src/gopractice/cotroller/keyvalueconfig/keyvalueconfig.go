package keyvalueconfig

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"github.com/jinzhu/gorm"
	"fmt"
	"net/http"
)

//设置key value
func SetKeyValue(c *gin.Context) {
	sendErrJson := common.SendErrJson
	type ReqData struct {
		KeyName string `json:"key" binding:"required,min=1"`
		Value   string `json:"value" binding:"required,min=1"`
	}

	var reqData ReqData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		sendErrJson("无效的参数", c)
		return
	}

	var keyValueConfig model.KeyValueConfig
	if err := model.DB.Where("key_name = ?", reqData.KeyName).Find(&keyValueConfig).Error; err != nil {
		if err != gorm.ErrRecordNotFound { //这里表示已经存在  但是出错了  返回错误
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}

		var theKeyValueConfig model.KeyValueConfig
		theKeyValueConfig.KeyName = reqData.KeyName
		theKeyValueConfig.Value = reqData.Value

		if err := model.DB.Create(&theKeyValueConfig).Error; err != nil {
			fmt.Println(err.Error())
			sendErrJson("error", c)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.SUCCESS,
			"msg":   "success",
			"data": gin.H{
				"id": theKeyValueConfig.ID,
			},
		})
		return
	}

	//这里表示更新数据
	keyValueConfig.Value = reqData.Value
	if err := model.DB.Save(&keyValueConfig).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"id": keyValueConfig.ID,
		},
	})
}
