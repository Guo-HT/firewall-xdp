package utils

import "xdpEngine/db"

// UserDbToWeb 将数据库查询的数据转化为接口数据，隐藏掉不必要的数据
func UserDbToWeb(data []db.User) (userList []UserInfo) {
	for _, user := range data {
		userList = append(userList, UserInfo{
			ID:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
			Role:     user.Role,
			CreateAt: user.CreateAt,
		})
	}
	return
}
