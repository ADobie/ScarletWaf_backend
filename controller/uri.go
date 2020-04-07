package controller

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"scarlet/common"
	"scarlet/tool"
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
	if err != nil {
		tool.GetLogger().WithField("handler", "AddURI").Debug("JSON绑定失败", err)
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "JSON数据不合法" + err.Error(),
			Data: nil,
		})
		return
	} else {
		if err := uri.Validate(); err != nil {
			data, _ := json.Marshal(err)
			c.JSON(400, common.DataResponse{
				Code: 400,
				Msg:  "表单不合法",
				Data: string(data),
			})
			return
		}
		if !serverService.Own(user.ID, uri.ServerID) {
			c.JSON(401, common.OperationResponse{
				Code: 401,
				Msg:  "越权操作",
			})
			return
		} else if !(uriService.Exist(uri)) {
			server := serverService.Get(uri.ServerID)
			uri.Host = server.Domain
			uriService.Add(uri)
			c.JSON(200, common.OperationResponse{
				Code: 200,
				Msg:  "添加成功",
			})
			return
		} else {
			c.JSON(200, common.OperationResponse{
				Code: 400,
				Msg:  "已存在",
			})
		}
	}
}

// @Summary 删除URI
// @Tags uri
// @Accept json
// @Produce json
// @Param uri body common.URI true "server_id、id必填"
// @Success 200 {object} common.OperationResponse
// @Failure 400 {object} common.DataResponse
// @Router /user/uri/delete [post] 'Login required'
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
		if err := uri.Validate(); err != nil {
			data, _ := json.Marshal(err)
			c.JSON(400, common.DataResponse{
				Code: 400,
				Msg:  "表单不合法",
				Data: string(data),
			})
			return
		}
		if serverService.Own(user.ID, uri.ServerID) && uriService.Own(uri.ID, uri.ServerID) {
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
	if err != nil {
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "JSON绑定错误" + err.Error(),
			Data: nil,
		})
		return
	}
	if err = form.Validate(); err != nil {
		data, _ := json.Marshal(err)
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "表单不合法",
			Data: string(data),
		})
	}
	if serverService.Own(user.ID, form.ServerID) {
		uris := uriService.GetByServerID(form.ServerID)
		c.JSON(200, common.DataResponse{
			Code: 200,
			Msg:  "返回成功",
			Data: uris,
		})
		return
	} else {
		c.JSON(401, common.DataResponse{
			Code: 401,
			Msg:  "越权操作",
			Data: nil,
		})
		return
	}
}
