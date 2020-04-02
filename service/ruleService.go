package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"scarlet/common"
	"scarlet/tool"
	"time"
)

var (
	redisPool        = newPool()
	mysqlClient, err = gorm.Open("mysql", "scarlet:scarlet@(127.0.0.1:3306)/scarlet?charset=utf8&parseTime=True&loc=Local")
	serverService    = NewServerService()
)

func newPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (c redis.Conn, err error) {
			return redis.Dial("tcp", "localhost:6379")
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

type RuleService struct{}

func NewRuleService() *RuleService {
	return new(RuleService)
}

func (r *RuleService) GetRulePage(page common.RulePage) []common.Rule {
	var results []common.Rule
	switch page.Flag {
	case "BASE":
		server := common.Server{
			Model: gorm.Model{ID: page.ServerID},
		}
		mysqlClient.First(&server)
		key := tool.BaseRuleKeyGen(server.Domain, page.Type)
		start := (page.Page - 1) * page.Limit
		rules, err := redis.IntMap(redisPool.Get().Do("zrange", key, start, start+page.Limit, "withscores"))
		if err != nil {
			panic(err)
		}
		// debug
		for rule, hit := range rules {
			results = append(results, common.Rule{
				Content: rule,
				Hit:     hit,
				URI:     "",
				Host:    server.Domain,
				Flag:    "BASE",
			})
		}
		return results
	case "CUSTOM":
		server := common.Server{
			Model: gorm.Model{ID: page.ServerID},
		}
		mysqlClient.First(&server)
		uri := common.URI{
			Model: gorm.Model{ID: page.URIID},
		}
		mysqlClient.First(&uri)
		key := tool.CustomRuleKeyGen(server.Domain, uri.Path, page.Type)
		start := (page.Page - 1) * page.Limit
		rules, err := redis.IntMap(redisPool.Get().Do("zrange", key, start, start+page.Limit, "withscores"))
		if err != nil {
			panic(err)
		}
		// debug
		for rule, hit := range rules {
			results = append(results, common.Rule{
				Content: rule,
				Hit:     hit,
				URI:     uri.Path,
				Host:    server.Domain,
				Flag:    "CUSTOM",
			})
		}
		return results
	}
	return nil
}

//func GetBaseRulesByServer(server common.Server) ([]common.Rule, error) {
//
//}
//
//func GetCustomRulesByUser(service ServerService,uri common.URI) ([]common.Rule, error) {
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
//
//func Echo(cmd string) {
//	fmt.Println(cmd)
//}
