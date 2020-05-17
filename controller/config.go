package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"scarlet/common"
	"scarlet/service"
	"strings"
)

var configService = service.NewConfigService()

//@Summary 控制waf开关
//@Tags switch
//@Accept json
//@Produce json
//@Param switchForm body common.SwitchOperation true "可以不填写config_name"
//@Success 200 {object} common.DataResponse true
//@Response 400 common.DataResponse true
//@Router /user/switch/waf [POST]
func WafStatus(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	switchForm := common.SwitchOperation{}
	err := c.ShouldBindJSON(&switchForm)
	if OnJSONError(c, err) {
		return
	}
	switchForm.ConfigName = common.AbbrMap["waf"]
	switchForm.Format()
	err = switchForm.Validate()
	if OnValidateError(c, err) {
		return
	}
	if !serverService.Own(user.ID, switchForm.ServerID) {
		Failure(c, "越权操作", nil)
	} else {
		configService.WafStatus(switchForm)
		Success(c, "修改成功", nil)
	}
}

//@Summary 修改Server Switch或者URI Switch
//@Tags switch
//@Accept json
//@Produce json
//@Param switchForm body common.SwitchOperation true "必须填写config_name；不给uri_id则修改server"
//@Success 200 {object} common.DataResponse true
//@Response 400 common.DataResponse true
//@Router /user/switch/change [POST]
func ChangeSwitch(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	switchForm := common.SwitchOperation{}
	err := c.ShouldBindJSON(&switchForm)
	if OnJSONError(c, err) {
		return
	}
	switchForm.ConfigName = strings.ToLower(switchForm.ConfigName)
	switchForm.Format()
	err = switchForm.Validate()
	if OnValidateError(c, err) {
		return
	}
	if !serverService.Own(user.ID, switchForm.ServerID) {
		Failure(c, "越权操作", nil)
		return
	} else if switchForm.URIID != 0 && !uriService.Own(switchForm.URIID, switchForm.ServerID) {
		Failure(c, "越权操作", nil)
		return
	} else {
		configService.FunctionSwitch(switchForm)
		Success(c, "修改成功", nil)
	}
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
	switchForm := common.GetSwitchForm{}
	err := c.ShouldBindJSON(&switchForm)
	if OnJSONError(c, err) {
		return
	}
	err = switchForm.Validate()
	if OnValidateError(c, err) {
		return
	}
	if !serverService.Own(user.ID, switchForm.ServerID) {
		Failure(c, "越权操作", nil)
		return
	}
	if switchForm.URIID != 0 && !uriService.Own(switchForm.URIID, switchForm.ServerID) {
		Failure(c, "越权操作", nil)
		return
	}
	var res interface{}
	if switchForm.URIID != 0 {
		res = configService.GetCustomSwitch(switchForm.URIID)
	} else {
		res = configService.GetBaseSwitch(switchForm.ServerID)
	}
	Success(c, "获取成功", res)
	return
}
