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
	"strings"
	"strconv"
	"math/rand"
	"time"
	"crypto/md5"
	"gopractice/cotroller/mail"
)

const (
	activeDuration = 24 * 60 * 60
)

func sendEmail(action string,title string,curTime int64,user model.User,c *gin.Context) {
	siteName := config.ServerConfig.SiteName
	siteURL := "http://" + config.ServerConfig.Host
	secretStr := fmt.Sprintf("%d%s%s",curTime,user.Email,user.Pass)
	secretStr = fmt.Sprintf("%x",md5.Sum([]byte(secretStr)))
	actionURL := siteURL + action + "%d%s"

	actionURL = fmt.Sprintf(actionURL,user.ID,secretStr)

	fmt.Println(actionURL)

	content := "<p><b>亲爱的" + user.Name + ":</b></p>" +
		"<p>我们收到您在 " + siteName + " 的注册信息, 请点击下面的链接, 或粘贴到浏览器地址栏来激活帐号.</p>" +
		"<a href=\"" + actionURL + "\">" + actionURL + "</a>" +
		"<p>如果您没有在 " + siteName + " 填写过注册信息, 说明有人滥用了您的邮箱, 请删除此邮件, 我们对给您造成的打扰感到抱歉.</p>" +
		"<p>" + siteName + " 谨上.</p>"

	if action == "/reset" {
		content = "<p><b>亲爱的" + user.Name + ":</b></p>" +
			"<p>你的密码重设要求已经得到验证。请点击以下链接, 或粘贴到浏览器地址栏来设置新的密码: </p>" +
			"<a href=\"" + actionURL + "\">" + actionURL + "</a>" +
			"<p>感谢你对" + siteName + "的支持，希望你在" + siteName + "的体验有益且愉快。</p>" +
			"<p>(这是一封自动产生的email，请勿回复。)</p>"
	}
	content += "<p><img src=\"" + siteURL + "/images/logo.png\" style=\"height: 42px;\"/></p>"
	//fmt.Println(content)
	mail.SendEmail(user.Email,title,content)

}

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

//注册
func Signup(c *gin.Context) {
	sendErrJson := common.SendErrJson

	type userReqData struct {
		Name string `json:"name" binding:"required,min=4,max=20"`
		Email string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=20"`
	}

	var userData userReqData

	//这里是获取网络的数据
	if err := c.ShouldBindWith(&userData,binding.JSON);err != nil{
		fmt.Println(err.Error())
		sendErrJson("参数无效",c)
		return
	}

	userData.Name = util.AvoidXss(userData.Name)
	userData.Email = strings.TrimSpace(userData.Email)
	userData.Password = strings.TrimSpace(userData.Password)

	if strings.Index(userData.Email, "@") != -1 {
		sendErrJson("用户名中不能含有@",c)
		return
	}

	var user model.User

	if err := model.DB.Where("email = ? OR name = ?", userData.Email, userData.Name).Find(&user).Error; err == nil {
		if userData.Name == user.Name {
			sendErrJson("用户名" + user.Name + "已被注册",c)
			return
		}else if user.Email == userData.Email {
			sendErrJson("邮箱" + user.Email + "已存在",c)
			return
		}
	}

	var newUser model.User

	newUser.Name = userData.Name
	newUser.Email = userData.Email
	newUser.Pass = newUser.EntryPassword(userData.Password,newUser.Salt())
	newUser.Role = model.UserRoleNormal
	newUser.Status = model.UserStatusInActive
	newUser.Sex = model.UserSexMale
	newUser.AvatarURL = "/images/avatar/" + strconv.Itoa(rand.Intn(2)) + ".png"

	if err := model.DB.Create(&newUser).Error; err != nil {
		sendErrJson(err.Error(),c)
		return
	}

	curTime := time.Now().Unix()

	activeUser := fmt.Sprintf("%s%d",model.ActiveTime,newUser.ID)

	RedisConn := model.RedisPool.Get()

	defer RedisConn.Close()

	if _,err := RedisConn.Do("SET",activeUser,curTime,"EX",activeDuration);err != nil{
		fmt.Println("redis set failed",err)
	}

	go func() {
		sendEmail("/active","账号激活",curTime,user,c)
	}()

	c.JSON(http.StatusOK,gin.H{
		"errNo":model.ErrorCode.SUCCESS,
		"msg":"success",
		"data":gin.H{
			"email":user.Email,
		},
	})

}
