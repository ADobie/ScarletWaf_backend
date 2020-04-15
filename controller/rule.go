package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	if OnJSONError(c, err) {
		return
	}
	if serverService.Own(user.ID, rulePage.ServerID) {
		if rulePage.URIID != 0 && !uriService.Own(rulePage.URIID, rulePage.ServerID) {
			Failure(c, "越权操作", nil)
			return
		} else {
			if rulePage.URIID == 0 {
				rulePage.Flag = "BASE"
			} else {
				rulePage.Flag = "CUSTOM"
			}
			rules := ruleService.GetRulePage(rulePage)
			Success(c, "获取成功", rules)
			return
		}
	} else {
		Failure(c, "越权操作", nil)
		return
	}
}

//@Summary 删除规则
//@Tags rule
//@Accept json
//@Produce json
//@Param rulePageForm body common.RulePage true "必须给定server_id ，uri_id可选 type为get\post\ua\header\cookie之一"
//@Success 200 {object} common.DataResponse true
//@Response 400 common.DataResponse true
//@Router /user/rule/delete [POST]
func DeleteRule(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	rule := common.Rule{}
	err := c.ShouldBindJSON(&rule)
	if OnJSONError(c, err) {
		return
	}
	// 表单验证
	if rule.URIID == 0 {
		rule.Flag = "BASE"
	} else {
		rule.Flag = "CUSTOM"
	}
	rule.Format()
	err = rule.Validate()
	if OnValidateError(c, err) {
		return
	} else if ok := serverService.Own(user.ID, rule.ServerID) && uriService.Own(rule.URIID, rule.ServerID); !ok {
		Failure(c, "越权操作", nil)
		return
	} else {
		ruleService.Delete(rule)
		Success(c, "删除成功", nil)
		return
	}

}

//@Summary 增加规则
//@Tags rule
//@Accept json
//@Produce json
//@Param rulePageForm body common.RulePage true "必须给定server_id 、content，uri_id可选 type为get\post\ua\header\cookie之一"
//@Success 200 {object} common.DataResponse true
//@Response 400 common.DataResponse true
//@Router /user/rule/add [POST]
func AddRule(c *gin.Context) {
	var user common.User
	session := c.MustGet("session").(jwt.MapClaims)
	user = session["user"].(common.User)
	addRuleForm := common.AddRuleForm{}
	err := c.ShouldBindJSON(&addRuleForm)
	rules := addRuleForm.Rules
	if OnJSONError(c, err) {
		return
	}
	// 表单验证
	err = addRuleForm.Validate()
	if OnValidateError(c, err) {
		return
	}
	// 权限验证
	if ok := serverService.Own(user.ID, addRuleForm.ServerID) && uriService.Own(addRuleForm.URIID, addRuleForm.ServerID); !ok {
		Failure(c, "越权操作", nil)
		return
	} else {
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
		c.JSON(200, common.DataResponse{
			Code: 200,
			Msg:  "添加成功",
		})
	}
}
