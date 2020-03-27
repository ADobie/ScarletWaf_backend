package service

import (
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"scarlet/tool"
)

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
		DB:   viper.GetInt("redis.db"),
	})

	mysqlClient, err = gorm.Open("mysql", viper.GetString("mysql.user")+":"+viper.GetString("mysql.password")+"@("+viper.GetString("mysql.addr")+")/"+viper.GetString("mysql.database")+"?charset=utf8&parseTime=True&loc=Local")
)

func init() {
	if err != nil {
		tool.GetLogger().Fatal(err)
	}
	mysqlClient.SetLogger(tool.GetLogger().WithField("database", "mysql"))
}

//type RuleService struct{}
//
//func NewRuleService() *RuleService {
//	return new(RuleService)
//}
//
//func GetBaseRulesByUser(user common.User) ([]common.Rule, error) {
//
//}
//
//func GetCustomRulesByUser(user common.User, uri string) ([]common.Rule, error) {
//
//}
//
//func DeleteRules(rules [] common.Rule) error {
//
//}
//
//func DeleteRule(rule common.Rule) error {
//
//}

func Echo(cmd string) {
	fmt.Println(cmd)
}
