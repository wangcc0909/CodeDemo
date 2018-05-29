package middleware

import (
	"github.com/gin-gonic/gin"
	"gopractice/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"gopractice/config"
)

func getUser(c *gin.Context) (model.User, error) {
	var user model.User

	tokenString, tokenErr := c.Cookie("token")
	if tokenErr != nil {
		fmt.Println("未登录")
		return user, errors.New("未登录")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpcepte singing method %v", token.Header["alg"])
		}
		return []byte(config.ServerConfig.TokenSecret), nil
	})

	if err != nil {
		fmt.Println(err)
		return user, errors.New("未登录")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := uint(claims["id"].(float64))
		user, err = model.UserFromRedis(userId)

		if err != nil {
			return user, errors.New("未登录")
		}

		return user, nil
	}

	return user, errors.New("未登录")
}
