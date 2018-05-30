package common

import (
	"github.com/gin-gonic/gin"
	"gopractice/model"
	"fmt"
	"encoding/json"
	"net/http"
)

//返回网站信息
func SiteInfo(c *gin.Context) {
	var userCount int
	var topicCount int
	var replyCount int

	if err := model.DB.Model(&model.User{}).Count(&userCount).Error; err != nil {
		SendErrJson(err.Error(), c)
		return
	}

	if err := model.DB.Model(&model.Article{}).Count(&topicCount).Error; err != nil {
		SendErrJson(err.Error(), c)
		return
	}

	if err := model.DB.Model(&model.Comment{}).Count(&replyCount).Error; err != nil {
		SendErrJson(err.Error(), c)
		return
	}

	var keyValueConfig model.KeyValueConfig

	siteConfig := make(map[string]interface{})
	siteConfig["name"] = "还没有想好"
	siteConfig["icp"] = ""
	siteConfig["title"] = ""
	siteConfig["description"] = ""
	siteConfig["keywords"] = ""
	siteConfig["loginURL"] = "/images/logo.png"
	siteConfig["bdStatsID"] = ""
	siteConfig["luosimaoSiteKey"] = ""

	if err := model.DB.Where("key_name = \"site_config\"").Find(&keyValueConfig).Error; err != nil {
		fmt.Println(err.Error())
	}

	err := json.Unmarshal([]byte(keyValueConfig.Value), &keyValueConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	var baiduAdKeyValue model.KeyValueConfig
	baiduAdKeyConfig := make(map[string]interface{})
	baiduAdKeyConfig["banner760x90"] = ""
	baiduAdKeyConfig["banner2_760x90"] = ""
	baiduAdKeyConfig["banner3_760x90"] = ""
	baiduAdKeyConfig["ad250x250"] = ""
	baiduAdKeyConfig["ad120x90"] = ""
	baiduAdKeyConfig["ad20_3"] = ""
	baiduAdKeyConfig["ad20_3A"] = ""
	baiduAdKeyConfig["allowBaiduAd"] = false

	if err := model.DB.Where("key_name = \"baidu_ad_config\"").Find(&baiduAdKeyValue).Error; err != nil {
		fmt.Println(err.Error())
	}

	err = json.Unmarshal([]byte(baiduAdKeyValue.Value), &baiduAdKeyConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"siteConfig":    siteConfig,
			"baiduAdConfig": baiduAdKeyConfig,
			"userCount":     userCount,
			"topicCount":    topicCount,
			"replyCount":    replyCount,
		},
	})
}
