package service

import (
	"github.com/gomodule/redigo/redis"
	"scarlet/common"
	"scarlet/tool"
)

type ConfigService struct{}

var uriService = NewURIService()

func NewConfigService() *ConfigService {
	return new(ConfigService)
}

// 根据operation 写入WAF开关状态
func (c *ConfigService) WafStatus(operation common.SwitchOperation) {
	server := serverService.Get(operation.ServerID)
	key := tool.BaseConfigKeyGen(server.Domain)
	_, err := redis.Bool(redisPool.Get().Do("hset", key, operation.ConfigName, operation.ConfigStatus))
	if err != nil {
		panic(err)
	}
}

// 根据... 写入各个功能开关状态
func (c *ConfigService) FunctionSwitch(operation common.SwitchOperation) {
	server := serverService.Get(operation.ServerID)
	var key string
	switch {
	case operation.URIID != 0:
		uri := uriService.Get(operation.URIID)
		key = tool.CustomConfigKeyGen(server.Domain, uri.Path)
	case operation.URIID == 0:
		key = tool.BaseConfigKeyGen(server.Domain)
	}
	_, err := redis.Bool(redisPool.Get().Do("hset", key, operation.ConfigName, operation.ConfigStatus))
	if err != nil {
		panic(err)
	}
}

func (c *ConfigService) GetBaseSwitch(server_id uint) common.BaseSwitch {
	server := serverService.Get(server_id)
	key := tool.BaseConfigKeyGen(server.Domain)
	val, err := redis.Values(redisPool.Get().Do("hgetall", key))
	if err != nil {
		panic(err)
	}
	var switc common.BaseSwitch
	err = redis.ScanStruct(val, &switc)
	if err != nil {
		panic(err)
	}
	return switc
}

func (c *ConfigService) GetCustomSwitch(uri_id uint) common.CustomSwitch {
	uri := uriService.Get(uri_id)
	key := tool.CustomConfigKeyGen(uri.Host, uri.Path)
	val, err := redis.Values(redisPool.Get().Do("hgetall", key))
	if err != nil {
		panic(err)
	}
	var switc common.CustomSwitch
	err = redis.ScanStruct(val, &switc)
	if err != nil {
		panic(err)
	}
	return switc
}
