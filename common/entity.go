package common

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jinzhu/gorm"
	"strings"
)

type Rule struct {
	Content  string `json:"content"`
	Hit      int    `json:"hit"`
	URI      string `json:"uri"`
	Host     string `json:"host"`
	Flag     string `json:"flag"`
	Type     string `json:"type"`
	ServerID uint   `json:"server_id"`
	URIID    uint   `json:"uri_id"`
}

func (r *Rule) Format() {
	r.Flag = strings.ToUpper(r.Flag)
	r.Type = strings.ToUpper(r.Type)
}

type User struct {
	gorm.Model
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Servers  []Server `json:"servers" swaggerignore:"true"`
}

type URI struct {
	gorm.Model
	Path     string       `json:"path"`
	Host     string       `json:"host"`
	ServerID uint         `json:"server_id"`
	Switch   CustomSwitch `json:"-" gorm:"-"`
	Option   Option       `gorm:"-" json:"-"`
}

func (u URI) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ServerID, validation.Required.Error("server_id是必填的")),
	)
}

type Server struct {
	gorm.Model
	Domain string     `json:"domain"`
	IP     string     `json:"ip"`
	Switch BaseSwitch `gorm:"-" json:"-"`
	Option Option     `gorm:"-"`
	UserID uint       `json:"user_id"`
	URI    []URI      `gorm:"-"`
	ServerSwitch
}

type ServerSwitch struct {
	wafStatus bool `json:"waf_status"`
}

type BaseSwitch struct {
	IpBlacklist   bool `json:"ip_blacklist" redis:"ip_blacklist"`
	IpWhitelist   bool `json:"ip_whitelist" redis:"ip_whitelist"`
	GetArgsCheck  bool `json:"get_args_check" redis:"get_args_check"`
	PostArgsCheck bool `json:"post_args_check" redis:"post_args_check"`
	CookieCheck   bool `json:"cookie_check" redis:"cookie_check"`
	UaCheck       bool `json:"ua_check" redis:"ua_check"`
	CCDefense     bool `json:"cc_defense" redis:"cc_defense"`
	SqlTokenCheck bool `json:"sql_token_check" redis:"libsqli_token_check"`
}

type CustomSwitch struct {
	IpBlacklist   bool `json:"ip_blacklist" redis:"ip_blacklist"`
	IpWhitelist   bool `json:"ip_whitelist" redis:"ip_white_list"`
	GetArgsCheck  bool `json:"get_args_check" redis:"get_args_check"`
	PostArgsCheck bool `json:"post_args_check" redis:"post_args_check"`
	CookieCheck   bool `json:"cookie_check" redis:"cookie_check"`
}

type Option struct {
	CCRate    string
	ProxyPass string
}

// TODO:关于rune和普通字符的长度问题 Register的验证
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(3, 20).Error("用户名长度需不小于3但不大于20")),
		validation.Field(&u.Email, is.Email.Error("请输入合法的邮箱地址")),
	)
}

func (r Rule) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Type, validation.Required.Error("规则类型是必须声明的"), validation.In("GET", "POST", "COOKIE", "HEADER", "UA", "BLACKIP", "WHITEIP").Error("必须是GET POST COOKIE HEADER UA BLACKIP WHITEIP 之一")),
		validation.Field(&r.Flag, validation.Required.Error("必须声明是否为自定义的"), validation.In("BASE", "CUSTOM").Error("必须是BASE CUSTOM之一")),
	)
}
