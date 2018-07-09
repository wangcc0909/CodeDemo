package middleware

import (
	"github.com/gin-gonic/gin"
	"gopractice/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"errors"
	"gopractice/config"
	"gopractice/cotroller/common"
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

func SigninRequired(c *gin.Context) {
	sendErrJson := common.SendErrJson

	var user model.User
	var err error

	if user, err = getUser(c);err != nil {
		sendErrJson("未登录",model.ErrorCode.LoginTimeOut)
		return
	}

	c.Set("user",user)
	c.Next()
}

func EditorRequired(c *gin.Context) {
	sendErrJson := common.SendErrJson

	var user model.User
	var err error

	if user, err = getUser(c);err != nil {
		sendErrJson("未登录",model.ErrorCode.LoginTimeOut)
		return
	}

	if user.Role == model.UserRoleEditor || user.Role == model.UserRoleAdmin ||
		user.Role == model.UserRoleSuperAdmin || user.Role == model.UserRoleCrawler {
		c.Set("user",user)
		c.Next()
	}else {
		sendErrJson("没有权限",c)
	}

}

//给context设置user
func SetContextUser(c *gin.Context) {
	var user model.User
	var err error

	if user,err = getUser(c);err != nil {
		c.Set("user",nil)
		c.Next()
		return
	}

	c.Set("user",user)
	c.Next()
}

//必须是管理员
func AdminRequired(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var user model.User
	var err error
	if user, err = getUser(c);err != nil {
		sendErrJson("未登录",model.ErrorCode.LoginTimeOut,c)
		return
	}

	if user.Role == model.UserRoleCrawler || user.Role == model.UserRoleAdmin || user.Role == model.UserRoleSuperAdmin {
		c.Set("user",user)
		c.Next()
	}else {
		sendErrJson("没有权限",c)
	}
}
