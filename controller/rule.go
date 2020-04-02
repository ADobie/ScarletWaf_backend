package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"scarlet/common"
)

// AddRule godoc
// @Summary addRule
// @Description addRule
// @Accept json
// @Produce json
// @Param rule body common.Account true "test"
// @Param token header string true "test"
// @Success 200 {object} common.Account "test"
// @Failure 400 {object} common.Account "test"
// @Router /user/rule [POST]
func AddRule(c *gin.Context) {
	rule := common.Rule{}
	err := c.ShouldBindJSON(&rule)
	if err != nil {
		logrus.Fatal("Error binding json to rule struct :", err)
	}
	// rule dao service 检查
	// 给出其他响应
	//val,err:=redis.Zset(tool.RuleKeyGen(rule),rule.Content)
	//if err!=nil{this.logger.Fatal("Err set rule :",err)}
	c.JSON(200, common.OperationResponse{
		Code: 200,
		Msg:  "操作成功",
	})
}

//func GetRule(c *gin.Context) {
//	service.Echo("233")
//	c.String(200, "应该创建好了")
//}
