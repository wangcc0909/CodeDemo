package user

import (
	"github.com/gin-gonic/gin"
	"gopractice/cotroller/common"
	"github.com/gin-gonic/gin/binding"
	"fmt"
	"gopractice/util"
	"gopractice/config"
	"gopractice/model"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func Signin(c *gin.Context) {
	sendErrJson := common.SendErrJson

	type EmailLogin struct {
		SigninInput string `json:"signinInput" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6,max=20"`
		LuosimaoRes string `json:"luosimaoRes"`
	}

	type UserNameLogin struct {
		SigninInput string `json:"signinInput" binding:"required,min=4,max=20"`
		Password    string `json:"password" binding:"required,min=6,max=20"`
		LuosimaoRes string `json:"luosimaoRes"`
	}

	var emailLogin EmailLogin
	var userNameLogin UserNameLogin
	var signinInput string
	var password string
	var luosimaoRes string
	var sql string

	if c.Query("loginType") == "email" {
		if err := c.ShouldBindWith(&emailLogin, binding.JSON); err != nil {
			fmt.Println(err.Error())
			sendErrJson("邮箱或密码错误", c)
			return
		}

		signinInput = emailLogin.SigninInput
		password = emailLogin.Password
		luosimaoRes = emailLogin.LuosimaoRes
		sql = "email = ?"
	} else if c.Query("loginType") == "username" {
		if err := c.ShouldBindWith(&userNameLogin, binding.JSON); err != nil {
			fmt.Println(err.Error())
			sendErrJson("用户名或密码错误", c)
			return
		}
		signinInput = userNameLogin.SigninInput
		password = userNameLogin.Password
		luosimaoRes = userNameLogin.LuosimaoRes
		sql = "name = ?"
	}

	err := util.LuosimaoVerify(config.ServerConfig.LuosimaoVertifyURL, config.ServerConfig.LuosimaoAPIKey, luosimaoRes)
	if err != nil {
		sendErrJson(err.Error(), c)
		return
	}

	//先验证账号在验证密码
	var user model.User
	err = model.DB.Where(sql, signinInput).First(&user).Error
	if err != nil {
		sendErrJson("账号不存在", c)
		return
	}

	if user.CheckPassword(password) { //验证密码
		if user.Status == model.UserStatusInActive {
			encodeEmail := base64.StdEncoding.EncodeToString([]byte(user.Email))

			c.JSON(200, gin.H{
				"errNo": model.ErrorCode.InActive,
				"msg":   "账号未激活",
				"data": gin.H{
					"email": encodeEmail,
				},
			})

			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id": user.ID,
		})

		tokenString, err := token.SignedString([]byte(config.ServerConfig.TokenSecret))

		if err != nil {
			fmt.Println(err.Error())
			sendErrJson("内部错误", c)
			return
		}

		if err := model.UserToRedis(user); err != nil {
			fmt.Println(err.Error())
			sendErrJson("内部错误", c)
			return
		}

		c.SetCookie("token", tokenString, config.ServerConfig.TokenMaxAge, "", "", true, true)

		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.SUCCESS,
			"msg":   "success",
			"data": gin.H{
				"token": tokenString,
				"user":  user,
			},
		})
		return
	}

	sendErrJson("用户名或密码错误",c)

}
