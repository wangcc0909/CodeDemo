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
	"errors"
	"github.com/gomodule/redigo/redis"
)

const (
	activeDuration = 24 * 60 * 60
	resetDuration  = 24 * 60 * 60
)

func sendEmail(action string, title string, curTime int64, user model.User, c *gin.Context) {
	siteName := config.ServerConfig.SiteName
	siteURL := "http://" + config.ServerConfig.Host
	secretStr := fmt.Sprintf("%d%s%s", curTime, user.Email, user.Pass)
	secretStr = fmt.Sprintf("%x", md5.Sum([]byte(secretStr)))
	actionURL := siteURL + action + "%d%s"

	actionURL = fmt.Sprintf(actionURL, user.ID, secretStr)

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
	mail.SendEmail(user.Email, title, content)

}

// ActiveSendMail 发送激活账号的邮件
func ActiveSendMail(c *gin.Context) {
	sendErrJson := common.SendErrJson

	// 接收到的email参数是加密后的，不能加email验证规则
	type ReqData struct {
		Email string `json:"email" binding:"required"`
	}
	var reqData ReqData

	if err := c.ShouldBindJSON(&reqData); err != nil {
		sendErrJson("参数无效", c)
		return
	}

	var user model.User
	user.Email = reqData.Email

	//解密
	var decodeBytes []byte
	var decodeErr error

	if decodeBytes, decodeErr = base64.StdEncoding.DecodeString(user.Email); decodeErr != nil {
		sendErrJson("参数错误", c)
		return
	}

	user.Email = string(decodeBytes)

	if err := model.DB.Where("email = ?", user.Email).Find(&user).Error; err != nil {
		sendErrJson("无效的邮箱", c)
		return
	}

	curTime := time.Now().Unix()

	activeUser := fmt.Sprintf("%s%d", model.ActiveTime, user.ID)

	RedisConn := model.RedisPool.Get()

	defer RedisConn.Close()

	if _, err := RedisConn.Do("SET", activeUser, curTime, "EX", activeDuration); err != nil {
		fmt.Println("Redis set faied ", err)
	}

	go func() {
		sendEmail("/active", "账号激活", curTime, user, c)
	}()

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"email": user.Email,
		},
	})
}

func verifyLink(cacheKey string, c *gin.Context) (model.User, error) {
	var user model.User
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil || userId <= 0 {
		return user, errors.New("无效的链接")
	}

	secret := c.Param("secret")
	if secret == "" {
		return user, errors.New("无效的链接")
	}

	RedisConn := model.RedisPool.Get()
	defer RedisConn.Close()

	emailTime, err := redis.Int64(RedisConn.Do("GET", cacheKey+fmt.Sprintf("%d", userId)))
	if err != nil {
		return user, errors.New("无效的链接")
	}

	if err := model.DB.First(&user, userId).Error; err != nil {
		return user, errors.New("无效的链接")
	}

	secretStr := fmt.Sprintf("%d%s%s", emailTime, user.Email, user.Pass)
	secretStr = fmt.Sprintf("%x", md5.Sum([]byte(secretStr)))

	if secret != secretStr {
		fmt.Println(secret, secretStr)
		return user, errors.New("无效的链接")
	}

	return user, nil
}

func ActiveAccount(c *gin.Context) {
	sendErrJson := common.SendErrJson

	var err error
	var user model.User

	if user, err = verifyLink(model.ActiveTime, c); err != nil {
		sendErrJson("激活链接已失效", c)
		return
	}

	if user.ID <= 0 {
		sendErrJson("激活链接已失效", c)
		return
	}

	updateData := map[string]interface{}{
		"status":       model.UserStatusActived,
		"activated_at": time.Now(),
	}

	if err := model.DB.Model(&user).Update(updateData).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	RedisConn := model.RedisPool.Get()
	defer RedisConn.Close()

	if _, err := RedisConn.Do("DEL", fmt.Sprintf("%s%d", model.ActiveTime, user.ID)); err != nil {
		fmt.Println("redis deleted failed:", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"email": user.Email,
		},
	})
}

func ResetPasswordMail(c *gin.Context) {
	sendErrJson := common.SendErrJson

	type UserReqData struct {
		Email       string `json:"email" binding:"required,email"`
		LuosimaoRes string `json:"luosimaoRes"`
	}

	var userReqData UserReqData

	if err := c.ShouldBindJSON(&userReqData); err != nil {
		sendErrJson("无效的邮箱", c)
		return
	}

	err := util.LuosimaoVerify(config.ServerConfig.LuosimaoVertifyURL, config.ServerConfig.LuosimaoAPIKey, userReqData.LuosimaoRes)
	if err != nil {
		sendErrJson(err.Error(), c)
		return
	}

	var user model.User

	if err := model.DB.Where("email = ?", userReqData.Email).Find(&user).Error; err != nil {
		sendErrJson("没有邮箱"+userReqData.Email+"的用户", c)
		return
	}

	curTime := time.Now().Unix()
	resetUser := fmt.Sprintf("%s%d", model.ResetTime, user.ID)

	RedisConn := model.RedisPool.Get()
	if _, err := RedisConn.Do("SET", resetUser, curTime, "EX", resetDuration); err != nil {
		fmt.Println("redis set failed:", err)
	}

	go func() {
		sendEmail("/ac", "修改密码", curTime, user, c)
	}()

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})

}

func VerifyResetPasswordLink(c *gin.Context) {
	sendErrJson := common.SendErrJson

	if _, err := verifyLink(model.ResetTime, c); err != nil {
		fmt.Println(err.Error())
		sendErrJson("重置密码链接已失效", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})

}

func ResetPassword(c *gin.Context) {
	sendErrJson := common.SendErrJson
	type UserReqData struct {
		Password string `json:"password" binding:"required,min=6,max=20"`
	}

	var userReqData UserReqData

	err := c.ShouldBindJSON(&userReqData)
	if err != nil {
		sendErrJson("参数无效", c)
		return
	}

	var user model.User
	var verifyErr error

	if user, verifyErr = verifyLink(model.ResetTime, c); verifyErr != nil {
		sendErrJson("重置链接已失效", c)
		return
	}

	if user.ID <= 0 {
		sendErrJson("重置链接已失效", c)
		return
	}

	if err := model.DB.Model(&user).Update("pass", user.Pass).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	RedisConn := model.RedisPool.Get()
	if _, err := RedisConn.Do("DEL", fmt.Sprintf("%s%d", model.ResetTime, user.ID)); err != nil {
		fmt.Println("redis delete failed:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})

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

	sendErrJson("用户名或密码错误", c)
}

//注册
func Signup(c *gin.Context) {
	sendErrJson := common.SendErrJson

	type userReqData struct {
		Name     string `json:"name" binding:"required,min=4,max=20"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=20"`
	}

	var userData userReqData

	//这里是获取网络的数据
	if err := c.ShouldBindWith(&userData, binding.JSON); err != nil {
		fmt.Println(err.Error())
		sendErrJson("参数无效", c)
		return
	}

	userData.Name = util.AvoidXss(userData.Name)
	userData.Email = strings.TrimSpace(userData.Email)
	userData.Password = strings.TrimSpace(userData.Password)

	if strings.Index(userData.Email, "@") != -1 {
		sendErrJson("用户名中不能含有@", c)
		return
	}

	var user model.User

	if err := model.DB.Where("email = ? OR name = ?", userData.Email, userData.Name).Find(&user).Error; err == nil {
		if userData.Name == user.Name {
			sendErrJson("用户名"+user.Name+"已被注册", c)
			return
		} else if user.Email == userData.Email {
			sendErrJson("邮箱"+user.Email+"已存在", c)
			return
		}
	}

	var newUser model.User

	newUser.Name = userData.Name
	newUser.Email = userData.Email
	newUser.Pass = newUser.EntryPassword(userData.Password, newUser.Salt())
	newUser.Role = model.UserRoleNormal
	newUser.Status = model.UserStatusInActive
	newUser.Sex = model.UserSexMale
	newUser.AvatarURL = "/images/avatar/" + strconv.Itoa(rand.Intn(2)) + ".png"

	if err := model.DB.Create(&newUser).Error; err != nil {
		sendErrJson(err.Error(), c)
		return
	}

	curTime := time.Now().Unix()

	activeUser := fmt.Sprintf("%s%d", model.ActiveTime, newUser.ID)

	RedisConn := model.RedisPool.Get()

	defer RedisConn.Close()

	if _, err := RedisConn.Do("SET", activeUser, curTime, "EX", activeDuration); err != nil {
		fmt.Println("redis set failed", err)
	}

	go func() {
		sendEmail("/active", "账号激活", curTime, user, c)
	}()

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"email": user.Email,
		},
	})
}

func Signout(c *gin.Context) {
	userInner, exist := c.Get("user")
	var user model.User
	if exist {
		user = userInner.(model.User)
		conn := model.RedisPool.Get()
		defer conn.Close()

		_, err := conn.Do("DEL", fmt.Sprintf("%s%d", model.LoginUser, user.ID))
		if err != nil {
			fmt.Println("redis delete error ", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  gin.H{},
	})
}

func SecretInfo(c *gin.Context) {
	if user, exist := c.Get("user"); exist {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.SUCCESS,
			"msg":   "success",
			"data": gin.H{
				"user": user,
			},
		})
	}
}

// InfoDetail 返回用户详情信息(教育经历、职业经历等)，包含一些私密字段
func InfoDetail(c *gin.Context) {
	sendErrJson := common.SendErrJson

	userInter, _ := c.Get("user")
	user := userInter.(model.User)

	if err := model.DB.First(&user, user.ID).Error; err != nil {
		sendErrJson("error", c)
		return
	}

	if err := model.DB.Model(&user).Related(&user.Schools).Error; err != nil {
		sendErrJson(err.Error(), c)
		return
	}

	if err := model.DB.Model(&user).Related(&user.Careers).Error; err != nil {
		sendErrJson(err.Error(), c)
		return
	}

	if user.Sex == model.UserSexFeMale {
		user.CoverURL = "https://www.golang123.com/upload/img/2017/09/13/d20f62c6-bd11-4739-b79b-48c9fcbce392.jpg"
	} else {
		user.CoverURL = "https://www.golang123.com/upload/img/2017/09/13/e672995e-7a39-4a05-9673-8802b1865c46.jpg"
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"user": user,
		},
	})
}

func PublicInfo(c *gin.Context) {
	sendErrJson := common.SendErrJson
	var userID int
	var err error

	if userID, err = strconv.Atoi(c.Param("id")); err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的ID", c)
		return
	}

	var user model.User

	if err = model.DB.First(&user, userID).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("无效的ID", c)
		return
	}

	if user.Sex == model.UserSexFeMale {
		user.CoverURL = "https://www.golang123.com/upload/img/2017/09/13/d20f62c6-bd11-4739-b79b-48c9fcbce392.jpg"
	} else {
		user.CoverURL = "https://www.golang123.com/upload/img/2017/09/13/e672995e-7a39-4a05-9673-8802b1865c46.jpg"
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": gin.H{
			"user": user,
		},
	})
}

func UploadAvatar(c *gin.Context) {
	sendErrJson := common.SendErrJson

	data, err := common.Upload(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.ERROR,
			"msg":   err.Error(),
			"data":  gin.H{},
		})
		return
	}

	avatarURL := data["url"].(string)
	userInter, _ := c.Get("user")

	user := userInter.(model.User)

	if err := model.DB.Model(&user).Update("avatar_url", avatarURL).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.ERROR,
			"msg":   err.Error(),
			"data":  gin.H{},
		})
		return
	}

	user.AvatarURL = avatarURL

	if model.UserToRedis(user) != nil {
		sendErrJson("error", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data":  data,
	})
}

func Top10(c *gin.Context) {
	TopN(c, 10)
}

func Top100(c *gin.Context) {
	TopN(c, 100)

}

func TopN(c *gin.Context, n int) {
	sendErrJson := common.SendErrJson

	var users []model.User

	if err := model.DB.Order("score DESC").Limit(n).Find(&users).Error; err != nil {
		fmt.Println(err.Error())
		sendErrJson("error", c)
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"errNo": model.ErrorCode.SUCCESS,
			"msg":   "success",
			"data": gin.H{
				"users": users,
			},
		})
	}
}
