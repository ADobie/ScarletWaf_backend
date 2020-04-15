package service

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"scarlet/tool"
	"time"
)

var (
	redisPool        = newPool()
	mysqlLink        = fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.Mysql.UserName, conf.Mysql.Password, conf.Mysql.Addr, conf.Mysql.Database)
	mysqlClient, err = gorm.Open("mysql", mysqlLink)
	serverService    = NewServerService()
	conf             = tool.GetConfig()
)

func newPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (c redis.Conn, err error) {
			return redis.Dial("tcp", conf.Redis.Addr, redis.DialPassword(conf.Redis.Password))
		},
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
	}
}

func init() {
	if err != nil {
		tool.GetLogger().Fatal(err)
	}
	mysqlClient.SetLogger(tool.GetLogger().WithField("database", "mysql"))
}
