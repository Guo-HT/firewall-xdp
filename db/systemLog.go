package db

// SetSystemLog 添加系统日志
func SetSystemLog(ip string, user string, option string, optResult bool) {
	dbConn.Create(&SystemLog{
		IP:        ip,
		Username:  user,
		Option:    option,
		OptResult: optResult,
	})
}

// SearchSystemLog 搜索系统日志
func SearchSystemLog(search string, startTime, endTime, pageNo, pageSize int, sort string) (sysLog []SystemLog, total int) {
	var sysLogTmp []SystemLog
	searchLike := "%" + search + "%"
	offset := (pageNo - 1) * pageSize
	_ = dbConn.Model(&SystemLog{}).Where("(ip like ? or option like ? or username like ?) and create_at > ? and create_at < ?", searchLike, searchLike, searchLike, startTime, endTime).Order("create_at " + sort).Limit(pageSize).Offset(offset).Find(&sysLog)
	resultCount := dbConn.Model(&SystemLog{}).Where("(ip like ? or option like ? or username like ?) and create_at > ? and create_at < ?", searchLike, searchLike, searchLike, startTime, endTime).Order("create_at " + sort).Find(&sysLogTmp)
	total = int(resultCount.RowsAffected)
	return
}
