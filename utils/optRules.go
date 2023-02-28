package utils

// ResetRulesId 整理协议ID, 从1开始
func ResetRulesId(protoRuleList []ProtoRule) (resetedProtoRuleList []ProtoRule) {
	for index, rule := range protoRuleList {
		rule.Id = index
		resetedProtoRuleList = append(resetedProtoRuleList, rule)
	}
	return
}

// DeleteRuleById 通过ID删除协议规则
func DeleteRuleById(protoRuleList []ProtoRule, id int) (deletedProtoRuleList []ProtoRule) {
	for _, value := range protoRuleList {
		if value.Id == id {
			continue
		} else {
			deletedProtoRuleList = append(deletedProtoRuleList, value)
		}
	}
	return
}
