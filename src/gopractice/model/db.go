package model

import (
	"github.com/jinzhu/gorm"
	"gopractice/config"
	"fmt"
	"os"
	"github.com/gomodule/redigo/redis"
	"time"
	"gopkg.in/mgo.v2"
)

var DB *gorm.DB

var RedisPool *redis.Pool
var MongoDB *mgo.Database

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
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", config.RedisConfig.URL, redis.DialPassword(config.RedisConfig.Password))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

}

func initMongoDB() {
	if config.MongoConfig.URL == "" {
		return
	}
	session,err := mgo.Dial(config.MongoConfig.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	session.SetMode(mgo.Monotonic,true)
	MongoDB = session.DB(config.MongoConfig.Database)
}

func init() {
	initDB()
	initRedisDB()
	initMongoDB()
}
