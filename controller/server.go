package controller

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"scarlet/common"
)

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
		Data: serverService.GetByUserID(user.ID),
	})
}

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
	if err != nil {
		data, _ := json.Marshal(err)
		c.JSON(406, common.DataResponse{
			Code: 406,
			Msg:  "表单不合法",
			Data: string(data),
		})
		return
	}
	if serverService.Own(user.ID, form.ServerID) {
		serverService.Delete(form.ServerID)
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
	c.JSON(200, common.OperationResponse{
		Code: 200,
		Msg:  "添加成功",
	})
}
