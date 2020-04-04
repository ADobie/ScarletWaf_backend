package tool

import "scarlet/common"

func RuleKeyGen(rule common.Rule) string {
	// validate field
	return rule.Host + "@" + rule.URI
}

func BaseRuleKeyGen(host string, tp string) string {
	return host + "@BASE@" + tp
}

func CustomRuleKeyGen(host string, path string, tp string) string {
	return host + "@" + path + "@" + tp
}

func BaseConfigKeyGen(host string) string {
	return "CONFIG" + "@" + host + "@BASE"
}

func CustomConfigKeyGen(host string, path string) string {
	return "CONFIG" + "@" + host + "@" + path
}
