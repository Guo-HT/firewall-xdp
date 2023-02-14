package utils

// AppendPortListDeduplicate 去重写入列表: 源列表，待插入列表
func AppendPortListDeduplicate(rawData, insertData []int) (appendData []int) {
	for _, insertValue := range insertData {
		if !IsIntInList(insertValue, rawData) {
			rawData = append(rawData, insertValue)
		}
	}
	return rawData
}

// IsIntInList 判断int是否存在int数组中
func IsIntInList(target int, list []int) (result bool) {
	for _, value := range list {
		if target == value {
			return true
		}
	}
	return false
}
