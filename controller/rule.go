package controller

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"
	"scarlet/common"
)

//@Summary 获取规则
//@Tags rule
//@Accept json
//@Produce json
//@Param rulePageForm body common.RulePage true "page为页号，limit为一页的最大数量，类型为get\post\ua\header\cookie之一"
//@Success 200 {object} common.DataResponse{data=[]common.Rule} true
//@Response 400 common.DataResponse true
//@Router /user/rule/get [POST]
func GetRules(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	rulePage := common.RulePage{}
	err := c.ShouldBindJSON(&rulePage)
	if err != nil {
		logrus.WithField("Handler", "GetRules").Debug("绑定json错误")
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "绑定json发生错误" + err.Error(),
			Data: nil,
		})
		return
	}
	if serverService.Own(user.ID, rulePage.ServerID) {
		if rulePage.URIID != 0 && !uriService.Own(rulePage.URIID, rulePage.ServerID) {
			c.JSON(401, common.DataResponse{
				Code: 401,
				Msg:  "越权操作",
				Data: nil,
			})
			return
		} else {
			if rulePage.URIID == 0 {
				rulePage.Flag = "BASE"
			} else {
				rulePage.Flag = "CUSTOM"
			}
			rules := ruleService.GetRulePage(rulePage)
			c.JSON(200, common.DataResponse{
				Code: 200,
				Msg:  "获取成功",
				Data: rules,
			})
			return
		}
	} else {
		c.JSON(401, common.DataResponse{
			Code: 401,
			Msg:  "越权操作",
			Data: nil,
		})
		return
	}
}

//@Summary 删除规则
//@Tags rule
//@Accept json
//@Produce json
//@Param rulePageForm body common.RulePage true "必须给定server_id ，uri_id可选 type为get\post\ua\header\cookie之一"
//@Success 200 {object} common.OperationResponse true
//@Response 400 common.DataResponse true
//@Router /user/rule/delete [POST]
func DeleteRule(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	rule := common.Rule{}
	err := c.ShouldBindJSON(&rule)
	if err != nil {
		logrus.WithField("Handler", "DeleteRule").Debug("绑定json错误")
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "绑定json发生错误" + err.Error(),
			Data: nil,
		})
		return
	}
	// 表单验证
	if rule.URIID == 0 {
		rule.Flag = "BASE"
	} else {
		rule.Flag = "CUSTOM"
	}
	rule.Format()
	if err := rule.Validate(); err != nil {
		data, _ := json.Marshal(err)
		c.JSON(406, common.DataResponse{
			Code: 406,
			Msg:  "表单不合法",
			Data: string(data),
		})
		return
	} else {
		ruleService.Delete(rule)
		c.JSON(200, common.OperationResponse{
			Code: 200,
			Msg:  "删除成功",
		})
		return
	}
	// 权限验证
	if ok := serverService.Own(user.ID, rule.ServerID) && uriService.Own(rule.URIID, rule.ServerID); !ok {
		c.JSON(401, common.OperationResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	} else {
		ruleService.Delete(rule)
		c.JSON(200, common.OperationResponse{
			Code: 200,
			Msg:  "删除成功",
		})
		return
	}

}

//@Summary 增加规则
//@Tags rule
//@Accept json
//@Produce json
//@Param rulePageForm body common.RulePage true "必须给定server_id 、content，uri_id可选 type为get\post\ua\header\cookie之一"
//@Success 200 {object} common.OperationResponse true
//@Response 400 common.DataResponse true
//@Router /user/rule/add [POST]
func AddRule(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	addRuleForm := common.AddRuleForm{}
	err := c.ShouldBindJSON(&addRuleForm)
	rules := addRuleForm.Rules
	if err != nil {
		logrus.WithField("Handler", "DeleteRule").Debug("绑定json错误")
		c.JSON(400, common.DataResponse{
			Code: 400,
			Msg:  "绑定json发生错误" + err.Error(),
			Data: nil,
		})
		return
	}
	// 表单验证
	err = addRuleForm.Validate()
	if err != nil {
		if e, ok := err.(validation.InternalError); ok {
			c.JSON(400, common.DataResponse{
				Code: 400,
				Msg:  "表单不合法",
				Data: e,
			})
			return
		} else {
			errs, _ := json.Marshal(err)
			c.JSON(400, common.DataResponse{
				Code: 400,
				Msg:  "表单不合法",
				Data: string(errs),
			})
			return
		}
	}

	// 权限验证
	if ok := serverService.Own(user.ID, addRuleForm.ServerID) && uriService.Own(addRuleForm.URIID, addRuleForm.ServerID); !ok {
		c.JSON(401, common.OperationResponse{
			Code: 401,
			Msg:  "越权操作",
		})
		return
	}

	// 一条规则属于 某个server （BASE rule） 或者 某个server的某个URI （CUSTOM rule）
	// 用户提供必要的content 和 server_id 和 可选的uri_id
	// 入库前查出Domain和Path 如 waf.heyao.top和/login 填入
	server := serverService.Get(addRuleForm.ServerID)
	for index, _ := range rules {
		rules[index].Host = server.Domain
	}
	if addRuleForm.URIID != 0 {
		uri := uriService.Get(addRuleForm.URIID)
		for index, _ := range rules {
			rules[index].URI = uri.Path
		}
	}
	ruleService.MustAdd(rules)
	c.JSON(200, common.OperationResponse{
		Code: 200,
		Msg:  "添加成功",
	})
}
