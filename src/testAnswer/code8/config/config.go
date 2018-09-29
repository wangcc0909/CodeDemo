package config

import (
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"gopractice/util"
)

var jsonData map[string]interface{}

//这里通过读取配置文件来启动
func initJson() {
	result, err := ioutil.ReadFile("src/config.json")
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

type serverConfig struct {
	SiteName           string
	Env                string
}

var ServerConfig serverConfig

func initServer() {
	util.SetStructByJson(&ServerConfig, jsonData["go"].(map[string]interface{}))
}

func init() {
	initJson()
	initDB()
	initServer()
}
