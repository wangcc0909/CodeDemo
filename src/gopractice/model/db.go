package model

import (
	"github.com/jinzhu/gorm"
	"gopractice/config"
	"fmt"
	"os"
	"github.com/gomodule/redigo/redis"
	"time"
)

var DB *gorm.DB

var RedisPool *redis.Pool

func initDB() {
	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if config.ServerConfig.Env == DevelopmentMode {
		db.LogMode(true)
	}

	db.DB().SetMaxIdleConns(config.DBConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.DBConfig.MaxOpenConns)

	DB = db
}

func initRedisDB() {
	RedisPool = &redis.Pool{
		MaxIdle:     config.RedisConfig.MaxIdle,
		MaxActive:   config.RedisConfig.MaxActive,
		IdleTimeout: 240 * time.Second,
		Wait:true,
		Dial: func() (redis.Conn, error) {
			conn,err := redis.Dial("tcp",config.RedisConfig.URL,redis.DialPassword(config.RedisConfig.Password))
			if err != nil {
				return nil, err
			}
			return conn,nil
		},
	}

}

func init() {
	initDB()
	initRedisDB()
}
