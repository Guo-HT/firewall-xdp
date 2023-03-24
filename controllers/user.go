package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"xdpEngine/db"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
)

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var username string
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("UserLogin: %s", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username, "用户登录", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()

	var json utils.UserLoginForm
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("UserLogin: 请求参数错误")
		db.SetSystemLog(c.ClientIP(), json.Username, "用户登录", false)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误, " + err.Error(),
			"data": "",
		})
		return
	} else {
		username = json.Username
		isRight, isDefaultPwdChg := db.IsUserInfoRight(json.Username, json.Password)
		//fmt.Println(json.Username, json.Password, isRight, isFirst)
		if !isRight {
			logger.Println("用户名或密码错误")
			db.SetSystemLog(c.ClientIP(), json.Username, "用户登录", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "用户名或密码错误",
				"data": gin.H{
					"isFirstLogin": isDefaultPwdChg,
				},
			})
			return
		} else {
			user, err := db.GetUserInfoByUsername(json.Username)
			if err != nil {
				panic(err.Error())
			}
			if !isDefaultPwdChg {
				logger.Printf("用户%s未修改密码", json.Username)
				db.SetSystemLog(c.ClientIP(), json.Username, "用户登录", false)
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "首次登录，需修改密码",
					"data": gin.H{
						"isDefaultPwdChanged": isDefaultPwdChg,
					},
				})
				return
			} else {
				logger.Printf("用户%s登录成功", json.Username)

				session := sessions.Default(c)
				session.Set(systemConfig.SessionKeyUserId, user.ID)
				session.Set(systemConfig.SessionKeyUserName, user.UserName)
				session.Set(systemConfig.SessionKeyUserRole, user.Role)
				session.Set(systemConfig.SessionKeyUserOptTime, time.Now().Unix())
				_ = session.Save()
				db.SetSystemLog(c.ClientIP(), json.Username, "用户登录", true)
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "登录成功",
					"data": gin.H{
						"isDefaultPwdChanged": isDefaultPwdChg,
					},
				})
				return
			}
		}
	}

}

// GetAccessToken 获取访问token(供第三方)
func GetAccessToken(c *gin.Context) {

}

// UserLogout 用户退出登录
func UserLogout(c *gin.Context) {
	var username string
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("UserLogout error: %s\n%s", err, debug.Stack())
			db.SetSystemLog(c.ClientIP(), username, "用户退出登录", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
			return
		}
	}()
	session := sessions.Default(c)
	username = session.Get(systemConfig.SessionKeyUserName).(string)
	logger.Printf("用户 %s 退出登录", username)
	session.Delete(systemConfig.SessionKeyUserId)
	session.Delete(systemConfig.SessionKeyUserName)
	session.Delete(systemConfig.SessionKeyUserRole)
	session.Delete(systemConfig.SessionKeyUserOptTime)
	session.Clear()
	_ = session.Save() // 一定要Save，不然删不掉
	//fmt.Println(session.Get(systemConfig.SessionKeyUserId), session.Get(systemConfig.SessionKeyUserName), session.Get(systemConfig.SessionKeyUserRole))
	db.SetSystemLog(c.ClientIP(), username, "用户退出登录", true)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出登录成功",
		"data": "",
	})
	return
}

// UserInfo 获取当前用户信息
func UserInfo(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("UserInfo error: %s\n%s", err, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()
	session := sessions.Default(c)
	var loginState bool
	sessUserId := session.Get(systemConfig.SessionKeyUserId)
	sessUserName := session.Get(systemConfig.SessionKeyUserName)
	sessUserRole := session.Get(systemConfig.SessionKeyUserRole)
	if sessUserId == nil || sessUserName == nil || sessUserRole == nil {
		loginState = false
	} else {
		loginState = true
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取用户状态成功",
		"data": gin.H{
			"loginState":   loginState,
			"sessUserId":   sessUserId,
			"sessUserName": sessUserName,
			"sessUserRole": sessUserRole,
		},
	})
	return
}

// AddUser 新增用户
func AddUser(c *gin.Context) {
	session := sessions.Default(c)
	var json utils.UserAdd
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("AddUser error: %s\n%s", err, debug.Stack())
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "新增用户"+json.UserName, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []utils.UserInfo{},
			})
			return
		}
	}()
	if err := c.ShouldBindJSON(&json); err != nil {
		db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "新增用户"+json.UserName, false)
		errlog.Printf("AddUser error: %s\n%s", err, debug.Stack())
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []utils.UserInfo{},
		})
		return
	} else {
		logger.Printf("正在添加用户%s...", json.UserName)
		//fmt.Printf("%#v\n", session.Get(systemConfig.SessionKeyUserId))
		//fmt.Printf("%#v\n", session.Get(systemConfig.SessionKeyUserName))
		//fmt.Printf("%#v\n", session.Get(systemConfig.SessionKeyUserRole))
		//fmt.Printf("%#v\n", session.Get(systemConfig.SessionKeyUserOptTime))

		// 密码确认
		if isRight, _ := db.IsUserInfoRight(session.Get(systemConfig.SessionKeyUserName).(string), json.CurUserPassword); !isRight {
			users, total := db.GetUserList(10, 1)
			userList := utils.UserDbToWeb(users)
			errlog.Println("新增用户失败: 密码错误")
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "新增用户"+json.UserName, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "用户新增失败: 密码错误",
				"data": gin.H{
					"total":    total,
					"pageNo":   0,
					"pageSize": 10,
					"data":     userList,
				},
			})
			return
		}

		// 权限判断
		if session.Get(systemConfig.SessionKeyUserRole).(int) != 0 {
			users, total := db.GetUserList(10, 1)
			userList := utils.UserDbToWeb(users)
			errlog.Println("新增用户失败: 权限错误")
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "新增用户"+json.UserName, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "用户新增失败: 权限错误",
				"data": gin.H{
					"total":    total,
					"pageNo":   0,
					"pageSize": 10,
					"data":     userList,
				},
			})
			return
		}

		// 添加用户
		err := db.AddUser(db.User{
			UserName: json.UserName,
			Password: json.Password,
			Email:    json.Email,
			Role:     json.Role,
		})
		if err != nil {
			users, total := db.GetUserList(10, 1)
			userList := utils.UserDbToWeb(users)
			errlog.Println("新增用户失败:", err)
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "新增用户"+json.UserName, false)
			var errStr = ""
			if strings.Contains(err.Error(), "UNIQUE") {
				errStr = ": 用户名已存在"
			}
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "用户新增失败" + errStr,
				"data": gin.H{
					"total":    total,
					"pageNo":   0,
					"pageSize": 10,
					"data":     userList,
				},
			})
			return
		} else {
			logger.Println("新增用户成功")
			users, total := db.GetUserList(10, 1)
			userList := utils.UserDbToWeb(users)
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "新增用户"+json.UserName, true)
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "用户新增成功",
				"data": gin.H{
					"total":    total,
					"pageNo":   0,
					"pageSize": 10,
					"data":     userList,
				},
			})
			return
		}
	}
}

// DelUser 软删除用户
func DelUser(c *gin.Context) {
	var json utils.DelUserCheck
	session := sessions.Default(c)
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("DelUser error: %s\n%s", err, debug.Stack())
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "删除用户"+json.TargetUserName, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []utils.UserInfo{},
			})
			return
		}
	}()

	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Printf("DelUser error: %s\n%s", err, debug.Stack())
		db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "删除用户"+json.TargetUserName, false)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []utils.UserInfo{},
		})
		return
	} else {
		logger.Printf("正在删除用户%s...", json.TargetUserName)

		// 判断密码是否正确
		username := session.Get(systemConfig.SessionKeyUserName).(string)
		isRight, _ := db.IsUserInfoRight(username, json.Password)
		if !isRight {
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "删除用户"+json.TargetUserName, false)
			users, total := db.GetUserList(10, 1)
			userList := utils.UserDbToWeb(users)
			errlog.Println("用户删除失败: 密码错误")
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "用户删除失败: 密码错误",
				"data": gin.H{
					"total":    total,
					"pageNo":   0,
					"pageSize": 10,
					"data":     userList,
				},
			})
			return
		}

		// 权限判断
		if session.Get(systemConfig.SessionKeyUserRole).(int) != 0 {
			users, total := db.GetUserList(10, 1)
			userList := utils.UserDbToWeb(users)
			errlog.Println("删除用户失败: 权限错误")
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "删除用户"+json.TargetUserName, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "用户删除失败: 权限错误",
				"data": gin.H{
					"total":    total,
					"pageNo":   0,
					"pageSize": 10,
					"data":     userList,
				},
			})
			return
		}

		// 执行删除
		err := db.DeleteUserByUsername(json.TargetUserName)
		if err != nil {
			panic(err.Error())
		}
		users, total := db.GetUserList(10, 1)
		userList := utils.UserDbToWeb(users)
		db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "删除用户"+json.TargetUserName, true)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "用户删除成功",
			"data": gin.H{
				"total":    total,
				"pageNo":   0,
				"pageSize": 10,
				"data":     userList,
			},
		})
		return
	}

}

// GetAllUsers 获取用户列表
func GetAllUsers(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("GetAllUsers error: %s\n%s", err, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []utils.UserInfo{},
			})
			return
		}
	}()
	pageNoStr := c.DefaultQuery("page_no", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	pageNo, err := strconv.Atoi(pageNoStr)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		errlog.Printf("GetAllUsers 请求参数不正确: %s", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数不正确",
			"data": []utils.UserInfo{},
		})
		return
	}
	if pageNo < 0 {
		pageNo = 0
	}
	if pageSize < 0 {
		pageSize = 0
	}
	users, total := db.GetUserList(pageSize, pageNo)
	//fmt.Println(users, total)
	userList := utils.UserDbToWeb(users)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取用户状态成功",
		"data": gin.H{
			"total":    total,
			"pageNo":   pageNo,
			"pageSize": pageSize,
			"data":     userList,
		},
	})
	return
}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	session := sessions.Default(c)
	var json utils.ChangeUserPassword
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("ChangePassword error: %s\n%s", err, debug.Stack())
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "用户"+json.Username+"修改密码", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()
	if err := c.ShouldBindJSON(&json); err != nil {
		db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "用户"+json.Username+"修改密码", false)
		errlog.Println("ChangePassword: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误, " + err.Error(),
			"data": "",
		})
		return
	} else {
		isRight, _ := db.IsUserInfoRight(json.Username, json.OldPassword)
		if !isRight {
			logger.Printf("用户%s修改密码失败：信息错误", json.Username)
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "用户"+json.Username+"修改密码", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "信息错误，请检查后重试",
				"data": "",
			})
			return
		} else {
			user, err := db.UpdateUserPassword(json.Username, json.NewPassword)
			if err != nil {
				panic(err.Error())
			}
			logger.Printf("用户%s修改密码成功", user.UserName)
			db.SetSystemLog(c.ClientIP(), session.Get(systemConfig.SessionKeyUserName).(string), "用户"+json.Username+"修改密码", true)
			session.Clear()
			_ = session.Save()
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "密码修改成功",
				"data": "",
			})
			return
		}
	}
}
