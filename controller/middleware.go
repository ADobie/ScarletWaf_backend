package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"scarlet/common"
	"scarlet/tool"
)

var (
	secret  = tool.SecretGen()
	JWTNAME = "SCARLET"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("session")
		session := val.(jwt.MapClaims)
		if !ok {
			tool.GetLogger().WithField("middleware", "LoginRequired").Debug("找不到session对象")
			c.Next()
		} else {
			if session["login"].(bool) != false {
				user := userService.Get(uint(session["user_id"].(float64)))
				session["user"] = user
				c.Next()
			} else {
				c.JSON(401, common.DataResponse{
					Code: 401,
					Msg:  "未登录",
					Data: nil,
				})
				c.Abort()
				return
			}
			return
		}
	}
}

// 未设置JWT -- > 分配一个未登录的
// 已设置JWT -- > 解析验证
// 验证失败  -- > json: hacker
// 验证成功  -- > 解析后存入上下文
// c.Next() -- > 等待Handler调用
// 取session、加密、放到Cookie中
// 注意session类型为jwt.MapClaims
// 在取值时需要自己类型转换一下
// example 看controller.GetRules

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var val *http.Cookie
		var err error
		if val, err = c.Request.Cookie("SCARLET"); err != nil {
			session := jwt.MapClaims{}
			session["login"] = false
			c.Set("session", session)
			c.Next()
			return
			// set cookie with session
		} else {
			token, err := jwt.Parse(val.Value, func(token *jwt.Token) (i interface{}, err error) {
				return secret, nil
			})
			if err != nil {
				tool.GetLogger().WithField("JWT", err).Warn("签名不合法、解析发生错误")
				c.JSON(400, common.DataResponse{
					Code: 400,
					Msg:  "Hacker!",
					Data: nil,
				})
				c.Abort()
				return
			} else {
				var session = token.Claims.(jwt.MapClaims)
				c.Set("session", session)
				c.Next()
				// set cookie with session
			}
		}
	}
}

func JWT_TOKEN() gin.HandlerFunc {
	return func(c *gin.Context) {
		val := c.Request.Header.Get(JWTNAME)
		if val == "" {
			session := jwt.MapClaims{}
			session["login"] = false
			c.Set("session", session)
			c.Next()
			return
			// set cookie with session
		} else {
			token, err := jwt.Parse(val, func(token *jwt.Token) (i interface{}, err error) {
				return secret, nil
			})
			if err != nil {
				tool.GetLogger().WithField("JWT", err).Warn("签名不合法、解析发生错误")
				c.JSON(400, common.DataResponse{
					Code: 400,
					Msg:  "Hacker!",
					Data: nil,
				})
				c.Abort()
				return
			} else {
				var session = token.Claims.(jwt.MapClaims)
				c.Set("session", session)
				c.Next()
				// set cookie with session
			}
		}
	}
}

func saveSession(c *gin.Context, session jwt.MapClaims) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, session)
	data, _ := token.SignedString(secret)
	//c.SetCookie(JWTNAME, data, 3600, "/", ".localhost", 0, true, true)
	c.Writer.Header().Set(JWTNAME, string(data))
	return
}
