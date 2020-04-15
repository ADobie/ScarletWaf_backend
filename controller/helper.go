package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"scarlet/common"
)

func OnJSONError(c *gin.Context, err error) bool {
	if err != nil {
		Failure(c, "Error Binding JSON Data"+err.Error(), nil)
		return true
	}
	return false
}

func OnValidateError(c *gin.Context, err error) bool {
	if val, ok := err.(validation.InternalError); ok {
		panic(val)
		return true
	}
	if err != nil {
		data, _ := json.Marshal(err)
		c.JSON(406, common.DataResponse{
			Code: 406,
			Msg:  "表单不合法",
			Data: string(data),
		})
		c.Abort()
		return true
	}
	return false
}

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, common.DataResponse{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
	c.Abort()
}

func Failure(c *gin.Context, msg string, data interface{}) {
	c.JSON(400, common.DataResponse{
		Code: 400,
		Msg:  msg,
		Data: data,
	})
	c.Abort()
}
