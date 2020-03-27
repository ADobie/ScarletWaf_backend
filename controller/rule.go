package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"scarlet/common"
	"scarlet/service"
)

// AddRule godoc
// @Summary 添加规则
// @Description 给某个用户添加一条规则
// @Accept json
// @Produce json
// @Param rule body common.Rule true "所需要添加的规则，注意需要携带"
// @Param token header string true "用户登录后获得的token"
// @Success 200 {object} common.OperationResponse "能显示么"
// @Failure 400 {object} common.OperationResponse "可以显示失败么"
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

func GetRule(c *gin.Context) {
	service.Echo("233")
	c.String(200, "应该创建好了")
}
