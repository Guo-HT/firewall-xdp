package utils

import (
	"time"
	"xdpEngine/db"
)

// UserDbToWeb 将数据库查询的数据转化为接口数据，隐藏掉不必要的数据
func UserDbToWeb(data []db.User) (userList []UserInfo) {
	for _, user := range data {
		userList = append(userList, UserInfo{
			ID:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
			Role:     user.Role,
			CreateAt: time.Unix(user.CreateAt, 0),
		})
	}
	return
}

// SystemLogDbToWeb 将数据库查询的系统操作日志转化为接口数据, 隐藏掉不必要的数据
func SystemLogDbToWeb(data []db.SystemLog) (logRsp []SystemLogRsp) {
	for _, log := range data {
		logRsp = append(logRsp, SystemLogRsp{
			IP:        log.IP,
			Username:  log.Username,
			Option:    log.Option,
			OptResult: log.OptResult,
			Time:      log.CreateAt,
		})
	}
	return
}
