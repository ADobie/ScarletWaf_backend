package controller

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
	"scarlet/common"
	"scarlet/service"
)

var userService = service.NewUserService()
var serverService = service.NewServerService()
var uriService = service.NewURIService()
var ruleService = service.NewRuleService()

// Add godoc
// @Summary 用户注册
// @Description 用户注册
// @Tags User
// @Accept json
// @Produce json
// @Param user body common.DataResponse true "注册的表单"
// @Success 200 {object} common.DataResponse
// @Failure 400 {object} common.DataResponse
// @Router /user [post]
func AddUser(c *gin.Context) {
	user := common.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil && common.DEVELOP {
		logrus.WithField("Handler", "Add").Debug("绑定json错误")
		Failure(c, "绑定JSON错误"+err.Error(), nil)
		return
	}
	errs := user.Validate()
	fmt.Printf("%v", errs)
	if errs != nil {
		if e, ok := errs.(validation.InternalError); ok {
			logrus.WithField("Handler", "Add").Debug("规则错误", e.InternalError())
		} else {
			data, _ := json.Marshal(e)
			Failure(c, "用户不合法", string(data))
			return
		}
	}
	if !userService.Exist(user) {
		userService.Add(user)
		Success(c, "注册成功", nil)
		return
	} else {
		Failure(c, "用户已存在", nil)
		return
	}
}

func UpdateUser(c *gin.Context) {

}

// Add godoc
// @Summary 用户登录
// @Description 用户登录
// @Tags User
// @Accept json
// @Produce json
// @Param user body common.User true "邮箱 密码必填"
// @Success 200 {object} common.DataResponse
// @Failure 400 {object} common.DataResponse
// @Router /login [post]
func UserLogin(c *gin.Context) {
	var user common.User
	err := c.ShouldBindJSON(&user)
	val, ok := c.Get("session")
	if OnJSONError(c, err) {
		return
	}
	session := val.(jwt.MapClaims)
	user, ok = userService.Auth(user)
	if !ok {
		Failure(c, "用户名或密码错误", nil)
		return
	} else {
		session["login"] = true
		session["user_id"] = user.ID
		saveSession(c, session)
		c.Set("session", session)
		Success(c, "登录成功", nil)
	}

}
