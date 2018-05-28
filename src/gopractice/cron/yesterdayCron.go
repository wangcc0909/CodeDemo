package cron

import (
	"gopractice/util"
	"gopractice/model"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

//这里是数据统计功能
func yesterdayCron() {
	var yesterdaySignupUserCount uint //昨天新建用户数
	var yesterdayTopicCount uint      //昨天新建话题数
	var yesterdayCommentCount uint    //昨天回复数
	var yesterdayBoolCount uint       //昨天新建图书数
	var yesterdayPV int               //昨天的PV
	var yesterdayUV uint               //昨天的UV

	todayTime := util.GetTodayTime()
	yesterdayTime := util.GetYesterdayTime()

	if err := model.DB.Model(&model.User{}).Where("activate_at >= ? AND activate_at < ?", yesterdayTime, todayTime).
		Count(&yesterdaySignupUserCount).Error; err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := model.DB.Model(&model.Article{}).Where("activate_at >= ? AND activate_at < ?", yesterdayTime, todayTime).
		Count(&yesterdayTopicCount).Error; err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := model.DB.Model(&model.Comment{}).Where("activate_at >= ? AND activate_at < ?", yesterdayTime, todayTime).
		Count(&yesterdayCommentCount).Error; err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := model.DB.Model(&model.Book{}).Where("activate_at >= ? AND activate_at < ?", yesterdayTime, todayTime).
		Count(&yesterdayBoolCount).Error; err != nil {
		fmt.Println(err.Error())
		return
	}

	var uvCount map[string]uint
	uvErr := model.MongoDB.C("userVisit").Pipe(
		[]bson.M{
			{
				"$match": bson.M{
					"date": bson.M{
						"$gte": yesterdayTime,
						"$lt":  todayTime,
					},
				},
			},
			{
				"$group":bson.M{
					"_id":"$clientID",
				},
			},
			{"$count":"uv"},
		},
	).AllowDiskUse().One(&uvCount)

	if uvErr != nil {
		fmt.Println(uvErr.Error())
	}else {
		yesterdayUV = uvCount["uv"]
	}

	//todo
}
