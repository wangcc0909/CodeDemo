package util

import (
	"time"
	"strconv"
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

func GetTodayYMD(sep string) string {
	now := time.Now()  //获取现在的时间
	year := now.Year()
	month := StrToIntMonth(now.Month().String())
	date := now.Day()

	var monthStr string
	var dateStr string

	if month < 9 {
		monthStr = "0" + strconv.Itoa(month + 1)
	}else {
		monthStr = strconv.Itoa(month + 1)
	}

	if date < 10 {
		dateStr = "0" + strconv.Itoa(date)
	}else {
		dateStr = strconv.Itoa(date)
	}

	return strconv.Itoa(year) + sep + monthStr + sep + dateStr
}
