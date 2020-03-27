package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
	"scarlet/common"
	"scarlet/service"
)

var userService = service.NewUserService()

// AddUser godoc
// @Summary 用户注册
// @Description 用户注册
// @Tags User
// @Accept json
// @Produce json
// @Param user body common.OperationResponse true "注册的表单"
// @Success 200 {object} common.OperationResponse true
// @Failure 400 {object} common.OperationResponse true
// @Router /user [post]
func AddUser(c *gin.Context) {
	user := common.User{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		logrus.WithField("Handler", "AddUser").Fatal("绑定json错误")
	}
	errs := user.Validate()
	fmt.Printf("%v", errs)
	if errs != nil {
		if e, ok := errs.(validation.InternalError); ok {
			logrus.WithField("Handler", "AddUser").Fatal("规则错误", e.InternalError())
		} else {
			c.JSON(400, common.DataResponse{
				Code: 400,
				Msg:  "用户不合法",
				Data: e,
			})
			return
		}
	}
	userService.AddUser(user)
	c.JSON(200, common.DataResponse{
		Code: 200,
		Msg:  "注册成功",
		Data: nil,
	})

}

// UpdateUserPssword godoc
// @Summary 用户密码更改
// @Description 用户密码更改
// @Tags User
// @Accept json
// @Produce json
// @Param user body common.UpdatePasswordForm true "修改密码的表单"
// @Success 200 {object} common.OperationResponse true
// @Failure 400 {object} common.OperationResponse true
// @Router /user [put]
func UpdateUser(c *gin.Context) {

}

func UpdateServer(c *gin.Context) {
	var servers []common.Server
	err := c.ShouldBindJSON(&servers)
	if err != nil {
		logrus.WithField("Handler", "UpdateServer").Fatal("绑定json错误")
	}
	user, _ := c.Get("user")
	val, _ := user.(common.User)
	val.Servers = servers
	userService.UpdateUserServers(val)
	c.JSON(200, common.OperationResponse{
		Code: 200,
		Msg:  "添加成功",
	})
}

func UserLogin(c *gin.Context) {
	var user common.User
	c.ShouldBindJSON(&user)
	var ok bool
	user, ok = userService.Auth(user)
	if !ok {
		c.JSON(400, common.OperationResponse{
			Code: 400,
			Msg:  "用户名或密码错误",
		})
	} else {
		token := tokenService.Login(user)
		c.JSON(200, common.DataResponse{
			Code: 200,
			Msg:  "登录成功",
			Data: gin.H{
				"token": token,
			},
		})
	}

}
