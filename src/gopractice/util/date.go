package util

import (
	"time"
	"strconv"
	"strings"
)

func StrToIntMonth(month string) int {

	var date = map[string]int{
		"January":   0,
		"February":  1,
		"March":     2,
		"April":     3,
		"May":       4,
		"June":      5,
		"July":      6,
		"August":    7,
		"September": 8,
		"October":   9,
		"November":  10,
		"December":  11,
	}
	return date[month]

}

func GetTodayYM(sep string) string {
	now := time.Now()
	year := now.Year()
	month := StrToIntMonth(now.Month().String())

	var monthStr string
	if month < 9 {
		monthStr = "0" + strconv.Itoa(month+1)
	} else {
		monthStr = strconv.Itoa(month + 1)
	}

	return strconv.Itoa(year) + sep + monthStr
}

func GetTodayYMD(sep string) string {
	now := time.Now() //获取现在的时间
	year := now.Year()
	month := StrToIntMonth(now.Month().String())
	date := now.Day()

	var monthStr string
	var dateStr string

	if month < 9 {
		monthStr = "0" + strconv.Itoa(month+1)
	} else {
		monthStr = strconv.Itoa(month + 1)
	}

	if date < 10 {
		dateStr = "0" + strconv.Itoa(date)
	} else {
		dateStr = strconv.Itoa(date)
	}

	return strconv.Itoa(year) + sep + monthStr + sep + dateStr
}

func GetYesterdayYMD(sep string) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	yesterday := today.Unix() - 24*60*60
	yesterdayTime := time.Unix(yesterday,0)
	yesterdayYMD := yesterdayTime.Format("2018-05-28")
	return strings.Replace(yesterdayYMD,"-",sep,-1)
}

//返回今天0点的时间
func GetTodayTime() time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	return today
}

//返回昨天0点的时间  today.Unix()  返回从1970到现在所经过的秒数
func GetYesterdayTime() time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	yesterday := today.Unix() - 24*60*60

	return time.Unix(yesterday,0)
}
