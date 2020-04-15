package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"scarlet/common"
)

// @Summary 增加URI
// @Tags uri
// @Accept json
// @Produce json
// @Param uri body common.URI true "server_id、path必填"
// @Success 200 {object} common.DataResponse{data=[]common.Server}
// @Failure 400 {object} common.DataResponse
// @Router /user/uri/add [post] 'Login required'
func AddURI(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	uri := common.URI{}
	err := c.ShouldBindJSON(&uri)
	if OnJSONError(c, err) {
		return
	}
	err = uri.Validate()
	if OnValidateError(c, err) {
		return
	}
	if !serverService.Own(user.ID, uri.ServerID) {
		Failure(c, "越权操作", nil)
		return
	} else if !(uriService.Exist(uri)) {
		server := serverService.Get(uri.ServerID)
		uri.Host = server.Domain
		uriService.Add(uri)
		Success(c, "添加成功", nil)
		return
	} else {
		Failure(c, "已存在", nil)
	}
}

// @Summary 删除URI
// @Tags uri
// @Accept json
// @Produce json
// @Param uri body common.URI true "server_id、id必填"
// @Success 200 {object} common.DataResponse
// @Failure 400 {object} common.DataResponse
// @Router /user/uri/delete [post] 'Login required'
func DeleteURI(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	uri := common.URI{}
	err := c.ShouldBindJSON(&uri)
	if OnJSONError(c, err) {
		return
	}
	err = uri.Validate()
	if OnValidateError(c, err) {
		return
	}
	if serverService.Own(user.ID, uri.ServerID) && uriService.Own(uri.ID, uri.ServerID) {
		uriService.Delete(uri)
		Success(c, "删除成功", nil)
		return
	} else {
		Failure(c, "越权操作", nil)
		return
	}
}

// @Summary 获取URI
// @Tags uri
// @Accept json
// @Produce json
// @Param uri body common.GetURIForm true "server_id必填"
// @Success 200 {object} common.DataResponse{data=[]common.URI}
// @Failure 400 {object} common.DataResponse
// @Router /user/uri/get [post] 'Login required'
func GetURI(c *gin.Context) {
	var user common.User
	var form common.GetURIForm
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	err := c.ShouldBindJSON(&form)
	if OnJSONError(c, err) {
		return
	}
	err = form.Validate()
	if OnValidateError(c, err) {
		return
	}
	if serverService.Own(user.ID, form.ServerID) {
		uris := uriService.GetByServerID(form.ServerID)
		Success(c, "获取成功", uris)
		return
	} else {
		Failure(c, "越权操作", nil)
		return
	}
}
