package common

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterForm struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// validatation库默认不支持OR条件，但是我们可以把他们成一个。

func (this *RegisterForm) Validate() error {
	return validation.ValidateStruct(this,
		validation.Field(&this.Name, validation.Required, validation.Length(3, 20)),
		validation.Field(&this.Email, is.Email.Error("请输入合法的邮箱地址")),
	)
}

func (this *RegisterForm) ServerAddressRule() validation.Rule {
	return validation.NewStringRule(func(s string) bool {
		return true
	}, "服务器地址不合法")
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdatePasswordForm struct {
	Password string `json:"password"`
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
