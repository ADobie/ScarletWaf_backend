package tool

import "scarlet/common"

func RuleKeyGen(rule common.Rule) string {
	// validate field
	return rule.Host + "@" + rule.URI
}
