package controller

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"scarlet/common"
	"scarlet/service"
	"scarlet/tool"
)

var configService = service.NewConfigService()

//@Summary 控制waf开关
//@Tags switch
//@Accept json
//@Produce json
//@Param switchForm body common.SwitchOperation true "可以不填写config_name"
//@Success 200 {object} common.OperationResponse true
//@Response 400 common.DataResponse true
//@Router /user/switch/waf [POST]
func WafStatus(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	switchForm := common.SwitchOperation{}
	err := c.ShouldBindJSON(&switchForm)
	if err != nil {
		c.JSON(400, common.OperationResponse{
			Code: 400,
			Msg:  "Error Binding JSON data " + err.Error(),
		})
		return
	}
	switchForm.ConfigName = common.AbbrMap["waf"]
	switchForm.Format()
	err = switchForm.Validate()
	if err != nil {
		if e, ok := err.(validation.InternalError); ok {
			tool.GetLogger().Fatal("Internal Error in ServerPower: ", e.InternalError())
		} else {
			data, _ := json.Marshal(err)
			c.JSON(400, common.DataResponse{
				Code: 400,
				Msg:  "表单验证错误",
				Data: string(data),
			})
			return
		}
	}
	if !serverService.Own(user.ID, switchForm.ServerID) {
		c.JSON(401, common.OperationResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	} else {
		configService.WafStatus(switchForm)
		c.JSON(200, common.OperationResponse{
			Code: 200,
			Msg:  "修改成功",
		})
		return
	}
}

//@Summary 修改Server Switch或者URI Switch
//@Tags switch
//@Accept json
//@Produce json
//@Param switchForm body common.SwitchOperation true "必须填写config_name；不给uri_id则修改server"
//@Success 200 {object} common.OperationResponse true
//@Response 400 common.DataResponse true
//@Router /user/switch/change [POST]
func ChangeSwitch(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	switchForm := common.SwitchOperation{}
	err := c.ShouldBindJSON(&switchForm)
	if err != nil {
		c.JSON(400, common.OperationResponse{
			Code: 400,
			Msg:  "Error Binding JSON data " + err.Error(),
		})
		return
	}
	switchForm.ConfigName = common.AbbrMap[switchForm.ConfigName]
	switchForm.Format()
	err = switchForm.Validate()
	if err != nil {
		if e, ok := err.(validation.InternalError); ok {
			tool.GetLogger().Fatal("Internal Error in ServerPower: ", e.InternalError())
		} else {
			data, _ := json.Marshal(err)
			c.JSON(400, common.DataResponse{
				Code: 400,
				Msg:  "表单验证错误",
				Data: string(data),
			})
			return
		}
	}
	if !serverService.Own(user.ID, switchForm.ServerID) {
		c.JSON(401, common.OperationResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	}
	if switchForm.URIID != 0 && !uriService.Own(switchForm.URIID, switchForm.ServerID) {
		c.JSON(401, common.OperationResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	}
	configService.FunctionSwitch(switchForm)
	c.JSON(200, common.OperationResponse{
		Code: 200,
		Msg:  "修改成功",
	})
	return

}

//@Summary 修改Server Switch或者URI Switch
//@Tags switch
//@Accept json
//@Produce json
//@Param switchForm body common.SwitchOperation true "server_id必填 uri_id选填"
//@Success 200 {object} common.DataResponse{data=common.CustomSwitch} true
//@Success 200 {object} common.DataResponse{data=common.BaseSwitch} true
//@Response 400 common.DataResponse true
//@Router /user/switch/get [POST]
func GetSwitch(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	switchForm := common.SwitchOperation{}
	err := c.ShouldBindJSON(&switchForm)
	if err != nil {
		c.JSON(400, common.OperationResponse{
			Code: 400,
			Msg:  "Error Binding JSON data " + err.Error(),
		})
		return
	}
	switchForm.ConfigName = common.AbbrMap[switchForm.ConfigName]
	switchForm.Format()
	err = switchForm.Validate()
	if err != nil {
		if e, ok := err.(validation.InternalError); ok {
			tool.GetLogger().Fatal("Internal Error in ServerPower: ", e.InternalError())
		} else {
			data, _ := json.Marshal(err)
			c.JSON(400, common.DataResponse{
				Code: 400,
				Msg:  "表单验证错误",
				Data: string(data),
			})
			return
		}
	}
	if !serverService.Own(user.ID, switchForm.ServerID) {
		c.JSON(401, common.OperationResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	}
	if switchForm.URIID != 0 && !uriService.Own(switchForm.URIID, switchForm.ServerID) {
		c.JSON(401, common.OperationResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	}
	var res interface{}
	if switchForm.URIID != 0 {
		res = configService.GetCustomSwitch(switchForm.URIID)
	} else {
		res = configService.GetBaseSwitch(switchForm.URIID)
	}
	c.JSON(200, common.DataResponse{
		Code: 200,
		Msg:  "修改成功",
		Data: res,
	})
	return
}
