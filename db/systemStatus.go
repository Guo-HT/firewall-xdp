package db

// GetServerRuntime 获取系统运行时间
func GetServerRuntime() (serverRuntime int64) {
	var sysSet SystemSetting
	dbConn.Model(&SystemSetting{}).First(&sysSet)
	serverRuntime = sysSet.RunTime
	return
}

// ServerRuntimeIncrease 系统运行时间增加
func ServerRuntimeIncrease() {
	var sysSet SystemSetting
	dbConn.Model(&SystemSetting{}).Where(&SystemSetting{}).First(&sysSet)
	serverRunTime := sysSet.RunTime
	dbConn.Model(&SystemSetting{}).Where(&sysSet).Updates(SystemSetting{RunTime: serverRunTime + 1})
}

// GetSystemSetting 获取系统信息
func GetSystemSetting() (sysSet SystemSetting) {
	dbConn.Model(&SystemSetting{}).First(&sysSet)
	return
}

// SetSystemSetting 配置系统信息
func SetSystemSetting(title, icon string) (newSysSet SystemSetting, e error) {
	var sysSet SystemSetting
	dbConn.First(&sysSet)
	sysSet.Title = title
	sysSet.Icon = icon
	dbConn.Save(sysSet)
	dbConn.First(&newSysSet)
	return
}
