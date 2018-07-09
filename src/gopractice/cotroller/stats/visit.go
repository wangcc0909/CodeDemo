package stats

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"gopractice/model"
	"github.com/globalsign/mgo/bson"
	"strconv"
	"time"
	"fmt"
	"net/http"
)

func PV(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var err error
	var userVisit model.UserVisit
	userVisit.ID = bson.NewObjectId()
	userVisit.Platform = c.Query("platform")
	userVisit.ClientID = c.Query("clientId")
	userVisit.OSName = c.Query("osName")
	userVisit.OSVersion = c.Query("osVersion")
	userVisit.Language = c.Query("language")
	userVisit.Country = c.Query("country")
	userVisit.DeviceModel = c.Query("deviceModel")
	userVisit.DeviceWidth, err = strconv.Atoi(c.Query("deviceWidth"))
	if err != nil {
		sendErrJson("无效的deviceWidth", c)
		return
	}

	userVisit.DeviceHeight, err = strconv.Atoi(c.Query("deviceHeight"))
	if err != nil {
		sendErrJson("无效的deviceHeight", c)
		return
	}

	userVisit.IP = c.ClientIP()
	userVisit.Date = time.Now()
	userVisit.Referrer = c.Query("referrer")
	userVisit.URL = c.Query("url")
	userVisit.BrowserName = c.Query("browserName")
	userVisit.BrowserVersion = c.Query("browserVersion")
	if userVisit.ClientID == "" {
		sendErrJson("clientID不能为空", c)
		return
	}

	if err := model.MongoDB.C("userVisit").Insert(&userVisit); err != nil {
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
