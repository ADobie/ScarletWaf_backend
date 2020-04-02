package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
	"scarlet/common"
	"scarlet/service"
	"scarlet/tool"
	"strconv"
)

var userService = service.NewUserService()
var serverService = service.NewServerService()
var uriService = service.NewURIService()
var ruleService = service.NewRuleService()

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

func UserLogin(c *gin.Context) {
	var user common.User
	c.ShouldBindJSON(&user)
	val, ok := c.Get("session")
	if !ok {
		tool.GetLogger().WithField("Handler", "UserLogin").Fatal("Error getting Session obj")
	}
	session := val.(jwt.MapClaims)
	user, ok = userService.Auth(user)
	if !ok {
		c.JSON(400, common.OperationResponse{
			Code: 400,
			Msg:  "用户名或密码错误",
		})
	} else {
		session["login"] = true
		session["user_id"] = user.ID
		saveSession(c, session)
		c.Set("session", session)
		c.JSON(200, common.DataResponse{
			Code: 200,
			Msg:  "登录成功",
		})
	}

}

// @Summary 获得用户的服务器列表
// @Description 获取用户的注册的服务器列表
// @Accept json
// @Produce json
// @Router /user/server [GET] 'Login required'
func GetServers(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	c.JSON(200, common.DataResponse{
		Code: 200,
		Msg:  "获取成功🐳",
		Data: serverService.GetServersByUserID(user.ID),
	})
}

func DeleteServer(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	val, err := strconv.ParseUint(c.Param("serverid"), 10, 64)
	serverId := uint(val)
	if err != nil {
		tool.GetLogger().WithField("Handler", "DeleteServer").Debug("serverID参数不合法", err)
		c.JSON(406, common.DataResponse{
			Code: 406,
			Msg:  "serverID参数不合法",
			Data: nil,
		})
		return
	}
	if serverService.Own(user.ID, serverId) {
		serverService.DeleteServerByServerID(serverId)
		c.JSON(200, common.OperationResponse{
			Code: 200,
			Msg:  "删除成功",
		})
		return
	} else {
		c.JSON(401, common.OperationResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	}
}

func UpdateServer(c *gin.Context) {
	var servers []common.Server
	err := c.ShouldBindJSON(&servers)
	if err != nil {
		logrus.WithField("Handler", "UpdateServer").Fatal("绑定json错误")
	}
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	user.Servers = servers
	userService.UpdateUserServers(user)
	c.JSON(200, common.OperationResponse{
		Code: 200,
		Msg:  "添加成功",
	})
}

// 给某个服务器添加URI
// TODO:并添加默认配置
func AddURI(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	uri := common.URI{}
	err := c.ShouldBindJSON(&uri)
	if err != nil {
		tool.GetLogger().WithField("handler", "AddURI").Debug("JSON绑定失败", err)
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "JSON数据不合法" + err.Error(),
			Data: nil,
		})
		return
	} else {
		if serverService.Own(user.ID, uri.ServerID) {
			// TODO:在执行数据库操作之前做好检查  待会写validate
			uriService.Add(uri)
			c.JSON(200, common.OperationResponse{
				Code: 200,
				Msg:  "添加成功",
			})

			return
		} else {
			c.JSON(401, common.OperationResponse{
				Code: 401,
				Msg:  "越权操作",
			})
			return
		}
	}
}

func DeleteURI(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	uri := common.URI{}
	err := c.ShouldBindJSON(&uri)
	if err != nil {
		tool.GetLogger().WithField("handler", "AddURI").Debug("JSON绑定失败", err)
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "JSON数据不合法" + err.Error(),
			Data: nil,
		})
		return
	} else {
		if serverService.Own(user.ID, uri.ServerID) && uriService.Own(uri.ID, uri.ServerID) {
			// TODO:在执行数据库操作之前做好检查  待会写validate
			uriService.Delete(uri)
			c.JSON(200, common.OperationResponse{
				Code: 200,
				Msg:  "删除成功",
			})

			return
		} else {
			c.JSON(401, common.OperationResponse{
				Code: 401,
				Msg:  "越权操作",
			})
			return
		}
	}
}

// 范围是 user-> server -> base -> Type
// user -> server -> custom -> Type
func GetRules(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	rulePage := common.RulePage{}
	err := c.ShouldBindJSON(&rulePage)
	if err != nil {
		logrus.WithField("Handler", "UpdateServer").Debug("绑定json错误")
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "绑定json发生错误" + err.Error(),
			Data: nil,
		})
		return
	}
	if serverService.Own(user.ID, rulePage.ServerID) && uriService.Own(rulePage.URIID, rulePage.ServerID) {
		rules := ruleService.GetRulePage(rulePage)
		c.JSON(200, common.DataResponse{
			Code: 200,
			Msg:  "获取成功",
			Data: rules,
		})
		return
	} else {
		c.JSON(401, common.DataResponse{
			Code: 401,
			Msg:  "越权操作",
			Data: nil,
		})
	}

}
