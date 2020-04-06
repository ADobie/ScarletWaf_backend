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
		server := serverService.Get(page.ServerID)
		key := tool.BaseRuleKeyGen(server.Domain, page.Type)
		start := (page.Page - 1) * page.Limit
		rules, err := redis.IntMap(redisPool.Get().Do("zrange", key, start, start+page.Limit, "withscores"))
		if err != nil {
			panic(err)
		}
		// debug
		for rule, hit := range rules {
			results = append(results, common.Rule{
				Content:  rule,
				Hit:      hit,
				URI:      "",
				Host:     server.Domain,
				Flag:     "BASE",
				ServerID: page.ServerID,
				Type:     page.Type,
			})
		}
		return results
	case "CUSTOM":
		server := serverService.Get(page.ServerID)
		uri := uriService.Get(page.URIID)
		key := tool.CustomRuleKeyGen(server.Domain, uri.Path, page.Type)
		start := (page.Page - 1) * page.Limit
		rules, err := redis.IntMap(redisPool.Get().Do("zrange", key, start, start+page.Limit, "withscores"))
		if err != nil {
			panic(err)
		}
		// debug
		for rule, hit := range rules {
			results = append(results, common.Rule{
				Content:  rule,
				Hit:      hit,
				URI:      uri.Path,
				Host:     server.Domain,
				Flag:     "CUSTOM",
				ServerID: page.ServerID,
				Type:     page.Type,
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

func (r *RuleService) Delete(rule common.Rule) {
	server := serverService.Get(rule.ServerID)
	uri := uriService.Get(rule.URIID)
	rule.Host = server.Domain
	rule.URI = uri.Path
	var key string
	switch rule.Flag {
	case "BASE":
		key = tool.BaseRuleKeyGen(rule.Host, rule.Type)
	case "CUSTOM":
		key = tool.CustomRuleKeyGen(rule.Host, rule.URI, rule.Type)
	}
	_, err := redis.Int(redisPool.Get().Do("zrem", key, rule.Content))
	if err != nil {
		panic(err)
	}
}

// 开发时使用painc
// 完成后返回error
func (r *RuleService) Add(rules []common.Rule) error {
	if rules == nil {
		panic("rules is nil")
	}
	conn := redisPool.Get()
	for _, rule := range rules {
		var key string
		switch rule.Flag {
		case "BASE":
			key = tool.BaseRuleKeyGen(rule.Host, rule.Type)
		case "CUSTOM":
			key = tool.CustomRuleKeyGen(rule.Host, rule.URI, rule.Type)
		}
		conn.Send("zadd", key, rule.Content)
	}
	conn.Flush()
	_, err := redis.Int(conn.Receive())
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (r *RuleService) MustAdd(rules []common.Rule) {
	if rules == nil {
		panic("rules is nil")
	}
	conn := redisPool.Get()
	for _, rule := range rules {
		var key string
		switch rule.Flag {
		case "BASE":
			key = tool.BaseRuleKeyGen(rule.Host, rule.Type)
		case "CUSTOM":
			key = tool.CustomRuleKeyGen(rule.Host, rule.URI, rule.Type)
		}
		conn.Send("zadd", key, 0, rule.Content)
	}
	conn.Flush()
	_, err := redis.Int(conn.Receive())
	if err != nil {
		panic(err)
	}
}
