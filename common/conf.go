package common

type MysqlConf struct {
	Addr     string `json:"addr"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type RedisConf struct {
	Addr     string `json:"addr"`
	Database string `json:"database"`
}

type Conf struct {
	Mysql MysqlConf `json:"mysql"`
	Redis RedisConf `json:"redis"`
}

var DEVELOP = true
