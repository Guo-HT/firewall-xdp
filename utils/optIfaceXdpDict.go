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

// AppendIPListDeduplicate 去重写入列表: 源列表，待插入列表
func AppendIPListDeduplicate(rawData, insertData []string) (appendData []string) {
	for _, insertValue := range insertData {
		if !IsStrInList(insertValue, rawData) {
			rawData = append(rawData, insertValue)
		}
	}
	return rawData
}

// IsStrInList 判断string是否存在string数组中
func IsStrInList(target string, list []string) (result bool) {
	for _, value := range list {
		if target == value {
			return true
		}
	}
	return false
}

// DeletePortList 删除母列表中的子列表
func DeletePortList(rawData, deleteData []int) (resultData []int) {
	for _, eachData := range deleteData {
		index := GetPortListIndex(eachData, rawData)
		if index >= 0 {
			rawData = append(rawData[:index], rawData[index+1:]...)
		}
	}
	return rawData
}

// GetPortListIndex 获取元素在列表中的索引
func GetPortListIndex(target int, rawData []int) (index int) {
	for idx, each := range rawData {
		if each == target {
			return idx
		}
	}
	return -1
}

// DeleteIpList 删除母列表中的子列表
func DeleteIpList(rawData, deleteData []string) (resultData []string) {
	for _, eachData := range deleteData {
		index := GetIpListIndex(eachData, rawData)
		if index >= 0 {
			rawData = append(rawData[:index], rawData[index+1:]...)
		}
	}
	return rawData
}

// GetIpListIndex 获取元素在列表中的索引
func GetIpListIndex(target string, rawData []string) (index int) {
	for idx, each := range rawData {
		if each == target {
			return idx
		}
	}
	return -1
}
