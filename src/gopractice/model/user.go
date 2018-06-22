package model

import (
	"time"
	"strconv"
	"fmt"
	"crypto/md5"
	"gopractice/config"
	"github.com/gomodule/redigo/redis"
	"encoding/json"
	"errors"
)

type User struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdateAt     time.Time  `json:"updateAt"`
	DeleteAt     *time.Time `sql:"index" json:"deleteAt"`
	ActivateAt   time.Time  `json:"activateAt"`
	Name         string     `json:"name"`
	Pass         string     `json:"-"`
	Email        string     `json:"email"`
	Sex          uint       `json:"sex"`
	Location     string     `json:"location"`
	Introduce    string     `json:"introduce"`
	Phone        string     `json:"phone"`
	Score        uint       `json:"score"`
	ArticleCount uint       `json:"articleCount"`
	CommentCount uint       `json:"commentCount"`
	CollectCount uint       `json:"collectCount"`
	Signature    string     `json:"signature"`
	Role         int        `json:"role"`
	AvatarURL    string     `json:"avatarURL"`
	CoverURL     string     `json:"coverURL"`
	Status       int        `json:"status"`
	Schools      []School   `json:"schools"`
	Careers      []Career   `json:"careers"`
}

func (user User) CheckPassword(password string) bool {
	if password == "" || user.Pass == "" {
		return false
	}

	return user.EntryPassword(password, user.Salt()) == user.Pass
}

//每个用户都有一个不同的盐
func (user User) Salt() string {
	var userSalt string
	if user.Pass == "" {
		userSalt = strconv.Itoa(int(time.Now().Unix()))
	} else {
		userSalt = user.Pass[0:10]
	}

	return userSalt
}

//给密码加密
func (user User) EntryPassword(password string, salt string) (hash string) {
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	hash = salt + password + config.ServerConfig.PassSalt
	hash = salt + fmt.Sprintf("%x", md5.Sum([]byte(hash)))
	return
}

func UserFromRedis(userId uint) (User, error) {
	loginUser := fmt.Sprintf("%s%d", LoginUser, userId)

	conn := RedisPool.Get()

	defer conn.Close()

	result, err := redis.Bytes(conn.Do("GET", loginUser))
	if err != nil {
		fmt.Println(err.Error())
		return User{}, err
	}

	var user User
	err = json.Unmarshal(result, &user)
	if err != nil {
		return User{}, errors.New("未登录  no login")
	}

	return user, nil
}

func UserToRedis(user User) error {
	result, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	loginUserKey := fmt.Sprintf("%s%d", LoginUser, user.ID)

	conn := RedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", loginUserKey, result, "EX", config.ServerConfig.TokenMaxAge)
	if err != nil {
		fmt.Println("redis set filed: ", err.Error())
		return err
	}

	return nil
}

const (
	//普通用户
	UserRoleNormal = 1

	//网站编辑
	UserRoleEditor = 2

	//管理员
	UserRoleAdmin = 3

	//超级管理员
	UserRoleSuperAdmin = 4

	//爬虫
	UserRoleCrawler = 5
)

const (
	//user inactive  用户未激活
	UserStatusInActive = 1

	//user actived  用户已激活
	UserStatusActived = 2

	//user frozen 用户已冻结
	UserStatusFrozen = 3
)

const (
	//user sex male  男
	UserSexMale = 0

	//user sex female 女
	UserSexFeMale = 1

	//maxUserNameLen 用户名的最大长度
	MaxUserNameLen = 20

	//minUserNameLen 用户名的最小长度
	MinUserNameLen = 4

	//maxPasswordLen 密码的最大长度
	MaxPassLen = 20

	//minPasswordLen 密码的最小长度
	MinPassLen = 6

	//MaxSignatureLen 个性签名的最大长度
	MaxSignatureLen = 200

	//MaxLocationLen  居住地的最大长度
	MaxLocationLen = 200

	//MaxIntroduceLen 个性签名的最大长度
	MaxIntroduceLen = 500
)
