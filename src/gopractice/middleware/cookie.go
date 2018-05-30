package middleware

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"gopractice/config"
	"gopractice/model"
)

func RefreshTokenCookie(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		fmt.Println(err)
	}

	if token != "" && err == nil {
		c.SetCookie("token", token, config.ServerConfig.TokenMaxAge, "/", "", true, true)
		if user, err := getUser(c); err == nil {
			err = model.UserToRedis(user)
		}
	}

	c.Next()

}
