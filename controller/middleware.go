package controller

import (
	"github.com/gin-gonic/gin"
	"scarlet/common"
	"scarlet/service"
)

var (
	tokenService = service.NewTokenService()
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if userId, err := tokenService.Get(token); err != nil {
			c.JSON(400, common.OperationResponse{
				Code: 400,
				Msg:  "用户未登录",
			})
		} else {
			user := userService.GetUserByID(userId)
			c.Set("user", user)
		}
		c.Next()
	}
}
