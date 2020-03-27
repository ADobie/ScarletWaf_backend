package common

import "github.com/nicksnyder/go-i18n/v2/i18n"

var (
	NameERR = i18n.Message{
		ID:          "name_err",
		Description: "default form name err msg",
		Other:       "用户名是必填的",
	}
	AgeERR = i18n.Message{
		ID:          "age_err",
		Description: "default form age err msg",
		Other:       "年龄是必填的",
	}
	ServerAddressERR = i18n.Message{
		ID:          "server_address_err",
		Description: "default form server address err msg",
		Other:       "服务器地址必须是一个IP地址或者是域名",
	}
	PhoneERR = i18n.Message{
		ID:          "phone_err",
		Description: "default form phone err msg",
		Other:       "手机格式不合法",
	}
	EmailERR = i18n.Message{
		ID:          "email_err",
		Description: "default form email err msg",
		Other:       "邮箱地址不合法",
	}
	WelcomeMsg = i18n.Message{
		ID:          "welcome_msg",
		Description: "Default Welome MSG",
		Other:       "你好",
	}
)
