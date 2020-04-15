package tool

import (
	"bufio"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"path"
	"scarlet/common"
)

var temp = `
[redis]
addr = "127.0.0.1:6379"
password = "123456"

[mysql]
addr = "127.0.0.1:3306"
database = "scarlet"
password = "123456"
username = "scarlet"`

var PROJECT_DIR, _ = os.Getwd()
var Conf common.Conf

func init() {
	if CheckInstall() {
		fmt.Printf("[+] 加载配置\n")
		Conf = readConf()
		fmt.Printf("%+v\n", Conf)
	} else {
		Install()
	}
}

func CheckInstall() bool {
	if _, err := os.Stat(path.Join(PROJECT_DIR, "scarlet.toml")); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func Install() {
	conf := common.Conf{}
	var input string
	var useDefault bool
	if _, err := toml.Decode(temp, &conf); err != nil {
		log.Fatal("ERR INSTALL", err)
	}
	fmt.Printf("开始配置数据库信息~(留空使用默认)\n")
	fmt.Printf("[+] Mysql 地址(默认为: 127.0.0.1:3306):")
	if input, useDefault = readFromStdin(); !useDefault {
		conf.Mysql.Addr = input
	}
	fmt.Printf("[+] Mysql 用户名(默认为: scarlet):")
	if input, useDefault = readFromStdin(); !useDefault {
		conf.Mysql.UserName = input
	}
	fmt.Printf("[+] Mysql 密码(默认为: 123456):")
	if input, useDefault = readFromStdin(); !useDefault {
		conf.Mysql.Password = input
	}
	fmt.Printf("[+] Mysql 数据库(默认为: scarlet):")
	if input, useDefault = readFromStdin(); !useDefault {
		conf.Mysql.Database = input
	}
	fmt.Printf("[+] Redis 地址(默认为: 127.0.0.1:6379):")
	if input, useDefault = readFromStdin(); !useDefault {
		conf.Redis.Addr = input
	}
	fmt.Printf("[+] Redis 密码(默认为: 123456):")
	if input, useDefault = readFromStdin(); !useDefault {
		conf.Redis.Password = input
	}
	writeConf(conf)
	fmt.Printf("[+] 配置写入成功\n")
}

func readConf() common.Conf {
	data, err := ioutil.ReadFile(path.Join(PROJECT_DIR, "scarlet.toml"))
	if err != nil {
		log.Fatal(err)
	}
	tempData := string(data)
	if _, err := toml.Decode(tempData, &Conf); err != nil {
		fmt.Printf("[+] 读取配置失败 :%s 是否重装(y/n):", err)
		input, _ := readFromStdin()
		if input == "y" {
			Install()
		} else {
			log.Fatal("请重新配置")
		}
	}
	return Conf
}

/*
	os.O_RDWR  open the file read-write
	os.O_CREATE  create file if not exist
*/
func writeConf(c common.Conf) {
	file, err := os.OpenFile(path.Join(PROJECT_DIR, "scarlet.toml"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	encoder := toml.NewEncoder(file)
	err = encoder.Encode(c)
	if err != nil {
		panic(err)
	}
}

// if input if blank ,which means USE DEFAULT
func readFromStdin() (string, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		data := scanner.Text()
		if data == "" {
			return "", true
		} else {
			return data, false
		}
	} else {
		panic("error or timeout")
	}
}
func GetConfig() common.Conf {
	return Conf
}
