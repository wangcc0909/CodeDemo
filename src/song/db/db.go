package db

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"os"
	"song/config"
)

var DB *gorm.DB

func initDB() {
	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if config.ServerConfig.Env == "development" {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(config.DBConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.DBConfig.MaxOpenConns)

	DB = db
}

func init() {
	initDB()
}
