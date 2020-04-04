package common

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterForm struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (this *RegisterForm) Validate() error {
	return validation.ValidateStruct(this,
		validation.Field(&this.Name, validation.Required, validation.Length(3, 20)),
		validation.Field(&this.Email, is.Email.Error("请输入合法的邮箱地址")),
	)
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordForm struct {
	Password string `json:"password"`
}

type FormatForm interface {
	validation.Validatable
	Format()
}

// Page 第几页
// Limit 一页最多
// Flag base 或者 custom之一
// Type get/post/header/cookie/ua/black_ip/white_ip 大小写无关
type RulePage struct {
	Page     uint   `json:"page"`
	Limit    uint   `json:"limit"`
	Flag     string `json:"flag"`
	Type     string `json:"type"`
	ServerID uint   `json:"server_id"`
	URIID    uint   `json:"uri_id"`
}

// 必要字段 "Content" "ServerID" "Flag" "Type"
// 可选 "URIID"
type DeleteRuleForm struct {
	Rule
	ServerID uint `json:"server_id"`
	URIID    uint `json:"uri_id"`
}

type GetURIForm struct {
	ServerID uint `json:"server_id"`
}

func (g GetURIForm) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.ServerID, validation.Required.Error("server_id是必填的")),
	)
}

type AddRuleForm struct {
	ServerID uint   `json:"server_id"`
	URIID    uint   `json:"uri_id"`
	Rules    []Rule `json:"rules"`
}

func (a AddRuleForm) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ServerID, validation.Required.Error("必须指定服务器ID")),
		validation.Field(&a.Rules),
	)
}

func (a *AddRuleForm) Format() {
	for index, rule := range a.Rules {
		rule.Format()
		a.Rules[index] = rule
	}
}

type GetServerForm struct {
	ServerID uint `json:"server_id"`
}

func (g GetServerForm) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.ServerID, validation.Required.Error("server_id 是必填的")),
	)
}

type SwitchOperation struct {
	ServerID     uint   `json:"server_id"`
	URIID        uint   `json:"uri_id"`
	ConfigName   string `json:"config_name"`
	ConfigStatus bool   `json:"config_value"`
}

var AbbrMap = map[string]string{"waf": "waf_status", "get": "get_args_check", "post": "post_args_check", "cookie": "cookie_check", "ua": "ua_check", "blacklist": "ip_blacklist", "whitelist": "ip_whitelist", "cc": "cc_defense", "sql": "libsqli_token_check"}

func (s SwitchOperation) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ServerID, validation.Required.Error("必须声明服务器ID")),
		validation.Field(&s.ConfigName, validation.By(RuleLimit)),
	)
}

func RuleLimit(item interface{}) error {
	configName := item.(string)
	if configName == "" {
		return nil
	}
	for _, val := range AbbrMap {
		if configName == val {
			return nil
		}
	}
	return errors.New("未定义规则")
}

func (s *SwitchOperation) Format() {
	if val, ok := AbbrMap[s.ConfigName]; ok {
		s.ConfigName = val
	}
	return
}
