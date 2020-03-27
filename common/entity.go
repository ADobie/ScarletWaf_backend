package common

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jinzhu/gorm"
)

type Rule struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Hit     int    `json:"hit"`
	URI     string `json:"uri"`
	Host    string `json:"host"`
	Flag    string `json:"flag"`
}

type User struct {
	gorm.Model
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Servers  []Server `json:"servers"`
}

type URI struct {
	gorm.Model
	Switch   Switch
	Path     string `json:"path"`
	Host     string `json:"host"`
	ServerID uint
}

type Server struct {
	gorm.Model
	Domain string `json:"domain"`
	IP     string `json:"ip"`
	Switch Switch
	UserID uint
	URI    []URI
}

type Switch struct {
	gorm.Model
	IpBlacklist   bool
	IpWhitelist   bool
	GetArgsCheck  bool
	PostArgsCheck bool
	CookieCheck   bool
	UaCheck       bool
	CCDefense     bool
	SqlTokenCheck bool
	URIID         uint
}

// TODO:关于rune和普通字符的长度问题 Register的验证
func (this *User) Validate() error {
	return validation.ValidateStruct(this,
		validation.Field(&this.Name, validation.Required, validation.Length(3, 20).Error("用户名长度需不小于3但不大于20")),
		validation.Field(&this.Email, is.Email.Error("请输入合法的邮箱地址")),
	)
}
