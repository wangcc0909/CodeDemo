package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
	"gopractice/config"
	"strings"
	"regexp"
	"gopractice/model"
	"fmt"
	"net/http"
)

//append 每次添加一个元素  多个用...展开
func getReqPath(c *gin.Context) string {
	pathArr := strings.Split(c.Request.URL.Path, "/")

	for i := len(pathArr) - 1; i >= 0; i-- {
		if pathArr[i] == "" {
			pathArr = append(pathArr[:i], pathArr[i:]...)
		}
	}

	for i, path := range pathArr {
		if match, err := regexp.MatchString("^[0-9]+$", path); match && err == nil {
			pathArr[i] = "id"
		}
	}

	pathArr = append([]string{strings.ToLower(c.Request.Method)}, pathArr...)
	return strings.Join(pathArr, "_")
}

func APIStatsD() gin.HandlerFunc { //gin中间件

	return func(c *gin.Context) {
		t := time.Now()
		c.Next()

		if config.StatsDConfig.URL == "" {
			return
		}

		duration := time.Since(t)
		durationMS := int64(duration / 1e6) //转成毫秒

		reqPath := getReqPath(c)

		if err := (*model.StatterClient).Timing(reqPath, durationMS, 1); err != nil {
			fmt.Println(err.Error())
		}

		stats := c.Writer.Status()
		if stats != http.StatusGatewayTimeout && durationMS > 5000 {
			timeoutPath := strings.Join([]string{"timeout", reqPath}, ":")
			if err := (*model.StatterClient).Inc(timeoutPath, 1, 1); err != nil {
				fmt.Println(err.Error())
			}
		}

	}

}
