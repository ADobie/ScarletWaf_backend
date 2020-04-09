package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"scarlet/common"
)

// @Summary 获得用户的服务器列表
// @Description 获取用户的注册的服务器列表
// @Tags server
// @Accept json
// @Produce json
// @Success 200 {object} common.DataResponse{data=[]common.Server})
// @Failure 400 {object} common.DataResponse
// @Router /user/server/get [post] 'Login required'
func GetServers(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	c.JSON(200, common.DataResponse{
		Code: 200,
		Msg:  "获取成功🐳",
		Data: serverService.GetByUserID(user.ID),
	})
}

// @Summary 删除服务
// @Tags server
// @Description
// @Accept json
// @Produce json
// @Param servers body  common.GetServerForm true "server_id为必要"
// @Success 200 {object} common.DataResponse
// @Failure 400 {object} common.DataResponse
// @Router /user/server/delete [post] 'Login required'
func DeleteServer(c *gin.Context) {
	var user common.User
	var form common.GetServerForm
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "Error Binding JSON data" + err.Error(),
			Data: nil,
		})
		return
	}
	err = form.Validate()
	OnValidateError(c, err)
	if serverService.Own(user.ID, form.ServerID) {
		serverService.Delete(form.ServerID)
		c.JSON(200, common.DataResponse{
			Code: 200,
			Msg:  "删除成功",
		})
		return
	} else {
		c.JSON(401, common.DataResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	}
}

// @Summary 添加Server
// @Tags server
// @Accept json
// @Produce json
// @Param servers body []common.Server true "服务器列表"
// @Success 200 {object} common.DataResponse{data=[]common.Server}
// @Failure 400 {object} common.DataResponse
// @Router /user/server/get [post] 'Login required'
func AddServer(c *gin.Context) {
	var servers []common.Server
	err := c.ShouldBindJSON(&servers)
	if err != nil {
		logrus.WithField("Handler", "UpdateServer").Fatal("绑定json错误")
	}
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	user.Servers = servers
	userService.UpdateServers(user)
	c.JSON(200, common.DataResponse{
		Code: 200,
		Msg:  "添加成功",
	})
}
