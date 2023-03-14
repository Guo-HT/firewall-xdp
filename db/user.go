package db

import "errors"

// IsUserInfoRight 检查用户名、密码是否匹配，默认密码是否修改
func IsUserInfoRight(userName, password string) (isRight bool, isPwdChanged bool) {
	user := User{}
	result := dbConn.Where(&User{UserName: userName, Password: password}).First(&user)
	//fmt.Println(result.RowsAffected)
	if result.RowsAffected > 0 {
		// 可查出数据，存在正确的用户
		isRight = true
		//result.First(&user)
		if user.IsChangedDefualtPwd {
			isPwdChanged = true
		} else {
			isPwdChanged = false
		}
	} else {
		// 查不出数据，用户名或密码错误
		isRight = false
		isPwdChanged = true
	}
	return
}

// GetUserInfoById 通过ID查找用户
func GetUserInfoById(id int) (user User, err error) {
	result := dbConn.Where(&User{ID: id}).First(&user)
	if result.Error != nil {
		errlog.Println("GetUserInfoById error: ", result.Error.Error())
		return User{}, result.Error
	}
	return user, nil
}

// GetUserInfoByUsername 通过用户名查找用户
func GetUserInfoByUsername(username string) (user User, err error) {
	result := dbConn.Where(&User{UserName: username}).First(&user)
	if result.Error != nil {
		errlog.Println("GetUserInfoByUsername error: ", result.Error.Error())
		return User{}, result.Error
	}
	return user, nil
}

// GetUserList 获取用户列表
func GetUserList(pageSize, pageNo int) (userList []User, total int) {
	offset := (pageNo - 1) * pageSize
	//fmt.Printf("offset: %d, pageSize: %d", offset, pageSize)
	result := dbConn.Model(&User{}).Limit(pageSize).Offset(offset).Find(&userList)
	//result := dbConn.Model(&User{}).Find(&userList)
	total = int(result.RowsAffected)
	return
}

// UpdateUserPassword 用户修改密码
func UpdateUserPassword(username string, newPassword string) (user User, err error) {
	user = User{}
	result := dbConn.Where(&User{UserName: username}).First(&user)
	//fmt.Println(result.RowsAffected)
	if result.RowsAffected > 0 {
		// 可查出数据，存在正确的用户
		dbConn.Model(&user).Where(&User{UserName: username}).Updates(User{IsChangedDefualtPwd: true, Password: newPassword})
		return user, nil
	}
	return user, result.Error
}

// DeleteUserByUsername 删除用户（用户名）
func DeleteUserByUsername(username string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			errlog.Println("DeleteUserByUsername error: ", r)
			err = errors.New("error")
		}
	}()
	user := User{UserName: username}
	dbConn.Where(&user).Delete(&user)
	return nil
}

// AddUser 添加用户
func AddUser(user User) (err error) {
	newUser := User{
		UserName: user.UserName,
		Password: user.Password,
		Email:    user.Email,
		Role:     user.Role,
	}
	result := dbConn.Create(&newUser)

	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
