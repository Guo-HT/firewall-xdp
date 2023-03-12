package db

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

func UpdateUserPassword(username string, newPassword string) (user User, err error) {
	user = User{}
	result := dbConn.Where(&User{UserName: username}).First(&user)
	//fmt.Println(result.RowsAffected)
	if result.RowsAffected > 0 {
		// 可查出数据，存在正确的用户
		//result.First(&user)
		user.IsChangedDefualtPwd = true
		user.Password = newPassword
		dbConn.Save(&user)
		return user, nil
	}
	return user, result.Error
}
