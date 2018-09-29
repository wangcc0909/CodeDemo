package config

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"unicode/utf8"
	"song/util"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var jsonData map[string]interface{}

//这里通过读取配置文件来启动
func initJson() {
	result, err := ioutil.ReadFile("src/song/config.json")
	if err != nil {
		fmt.Println("ReadFile error ", err.Error())
		os.Exit(-1)
	}

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
	Charset      string
	URL          string
	MaxIdleConns int
	MaxOpenConns int
}

var DBConfig dBConfig

func initDB() {
	util.SetStructByJSON(&DBConfig, jsonData["database"].(map[string]interface{}))
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database, DBConfig.Charset)
	DBConfig.URL = url
}

type serverConfig struct {
	SiteName           string
	Env                string
	LogDir             string
	LogFile            string
	MaxMultipartMemory int
	Port               int
}

var ServerConfig serverConfig

func initServer() {
	util.SetStructByJSON(&ServerConfig, jsonData["go"].(map[string]interface{}))
	sep := string(os.PathSeparator)
	execPath, _ := os.Getwd()
	length := utf8.RuneCountInString(execPath)
	lastChar := execPath[length-1:]
	if lastChar != sep {
		execPath = execPath + sep
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

func init() {
	initJson()
	initDB()
	initServer()
}
