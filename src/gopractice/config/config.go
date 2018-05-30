package config

import (
	"gopractice/util"
	"io/ioutil"
	"fmt"
	"os"
	"regexp"
	"encoding/json"
	"unicode/utf8"
	"strings"
)

var jsonData map[string]interface{}

//这里通过读取配置文件来启动
func initJson() {
	result, err := ioutil.ReadFile("src/config.json")
	if err != nil {
		fmt.Println("ReadFile error ", err.Error())
		os.Exit(-1)
	}

	configStr := string(result[:])

	reg := regexp.MustCompile(`/\*.*\*/`)//这里是将所有的注释都删掉

	configStr = reg.ReplaceAllString(configStr, "")
	result = []byte(configStr)

	if err := json.Unmarshal(result, &jsonData); err != nil {
		fmt.Println("invalid Config ", err.Error())
		os.Exit(-1)
	}
}

type dBConfig struct {
	Dialect      string
	Database     string
	User         string
	Password     string
	Host         string
	Port         int
	charset      string
	URL          string
	MaxIdleConns int
	MaxOpenConns int
}

var DBConfig dBConfig

func initDB() {
	util.SetStructByJson(&DBConfig, jsonData["database"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database, DBConfig.charset)
	DBConfig.URL = url
}

type redisConfig struct {
	Host      string
	Port      int
	Password  string
	URL       string
	MaxIdle   int
	MaxActive int
}

var RedisConfig redisConfig

func initRedis() {
	util.SetStructByJson(&RedisConfig,jsonData["redis"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%d",RedisConfig.Host,RedisConfig.Port)
	RedisConfig.URL = url
}

type mongoConfig struct {
	URL string
	Database string
}

var MongoConfig mongoConfig

func initMongo() {
	util.SetStructByJson(&MongoConfig,jsonData["mongo"].(map[string]interface{}))
	//MongoConfig.URL = fmt.Sprintf("")
}

type serverConfig struct {
	APIPoweredBy       string
	SiteName           string
	Host               string
	ImgHost            string
	Env                string
	LogDir             string
	LogFile            string
	APIPrefix          string
	UploadImgDir       string
	ImgPath            string
	MaxMultipartMemory int
	Port               int
	StatsEnable        bool
	TokenSecret        string
	TokenMaxAge        int
	PassSalt           string
	LuosimaoVertifyURL string
	LuosimaoAPIKey     string
	CrawlerName        string
	MailUser           string
	MailPass           string
	MailHost           string
	MailPort           int
	MailFrom           string
	Github             string
	BaiduPushLink      string
}

var ServerConfig serverConfig

func initServer() {
	util.SetStructByJson(&ServerConfig, jsonData["go"].(map[string]interface{}))
	sep := string(os.PathSeparator)
	execPath, _ := os.Getwd()
	length := utf8.RuneCountInString(execPath)
	lastChar := execPath[length-1:]
	if lastChar != sep {
		execPath = execPath + sep
	}

	if ServerConfig.UploadImgDir == "" {
		pathArr := []string{"website", "static", "upload", "img"}
		uploadImgDir := execPath + strings.Join(pathArr, sep)
		ServerConfig.UploadImgDir = uploadImgDir
	}

	ymdStr := util.GetTodayYMD("-")

	if ServerConfig.LogDir == "" {
		ServerConfig.LogDir = execPath
	} else {
		length = utf8.RuneCountInString(ServerConfig.LogDir)
		lastChar = ServerConfig.LogDir[length-1:]
		if lastChar != sep {
			ServerConfig.LogDir = ServerConfig.LogDir + sep
		}
	}

	ServerConfig.LogFile = ServerConfig.LogDir + ymdStr + ".log"
}

type statsDConfig struct {
	URL string
	Prefix string
}

var StatsDConfig statsDConfig
func initStatsD() {
	util.SetStructByJson(&StatsDConfig,jsonData["statsd"].(map[string]interface{}))
}

func init() {
	initJson()
	initDB()
	initRedis()
	initMongo()
	initServer()
	initStatsD()
}
