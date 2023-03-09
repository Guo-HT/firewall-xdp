package utils

// IntListLimit 返回分页数据
func IntListLimit(l []int, pageNo, pageSize int) (resultList []int, rightPageNo, rightPageSize int) {
	start, end, pNo := calcSliceDuring(len(l), pageNo, pageSize)
	if end == 0 {
		return l[start:], pNo, pageSize
	} else {
		return l[start:end], pNo, pageSize
	}
}

// StringListLimit 返回分页数据
func StringListLimit(l []string, pageNo, pageSize int) (resultList []string, rightPageNo, rightPageSize int) {
	start, end, pNo := calcSliceDuring(len(l), pageNo, pageSize)
	if end == 0 {
		return l[start:], pNo, pageSize
	} else {
		return l[start:end], pNo, pageSize
	}
}

// IntIntStructListLimit 返回分页数据
func IntIntStructListLimit(l []PortHitCount, pageNo, pageSize int) (resultList []PortHitCount, rightPageNo, rightPageSize int) {
	start, end, pNo := calcSliceDuring(len(l), pageNo, pageSize)
	if end == 0 {
		return l[start:], pNo, pageSize
	} else {
		return l[start:end], pNo, pageSize
	}
}

// StringIntStructListLimit 返回分页数据
func StringIntStructListLimit(l []IpHitCount, pageNo, pageSize int) (resultList []IpHitCount, rightPageNo, rightPageSize int) {
	start, end, pNo := calcSliceDuring(len(l), pageNo, pageSize)
	if end == 0 {
		return l[start:], pNo, pageSize
	} else {
		return l[start:end], pNo, pageSize
	}
}

// calcSliceDuring 计算分页
func calcSliceDuring(length, pageNo, pageSize int) (startIndex, endIndex, resultPageNo int) {
	if pageNo < 0 || pageSize < 0 {
		// 返回第一页数据，默认10条
		return 0, 9, 1
	}
	if length == 0 {
		return 0, 0, 1
	}
	if length < pageNo*pageSize { // 范围超过最后一页
		rightPageSize := length % pageSize
		rightPageNo := length / pageSize
		if rightPageSize == 0 {
			rightPageNo = rightPageNo - 1
			rightPageSize = pageSize
		}
		return length - rightPageSize, 0, rightPageNo + 1
	} else {
		startIndex = pageSize * (pageNo - 1)
		endIndex = pageSize * pageNo
		resultPageNo = pageNo
		return
	}
}
