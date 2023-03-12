package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"time"
	"xdpEngine/db"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
)

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("UserLogin: %s", debug.Stack())
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
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误, " + err.Error(),
			"data": "",
		})
		return
	} else {
		isRight, isDefaultPwdChg := db.IsUserInfoRight(json.Username, json.Password)
		//fmt.Println(json.Username, json.Password, isRight, isFirst)
		if !isRight {
			logger.Println("用户名或密码错误")
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
				c.JSON(http.StatusOK, gin.H{
					"code": 200,
					"msg":  "登录成功",
					"data": gin.H{
						"isDefaultPwdChanged": isDefaultPwdChg,
					},
				})
			}

		}
	}

}

// GetAccessToken 获取访问token(供第三方)
func GetAccessToken(c *gin.Context) {

}

// UserLogout 用户退出登录
func UserLogout(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("UserLogout error: %s\n%s", err, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
			return
		}
	}()
	session := sessions.Default(c)
	logger.Printf("用户 %s 退出登录", session.Get(systemConfig.SessionKeyUserName))
	session.Delete(systemConfig.SessionKeyUserId)
	session.Delete(systemConfig.SessionKeyUserName)
	session.Delete(systemConfig.SessionKeyUserRole)
	session.Delete(systemConfig.SessionKeyUserOptTime)
	session.Clear()
	_ = session.Save() // 一定要Save，不然删不掉
	//fmt.Println(session.Get(systemConfig.SessionKeyUserId), session.Get(systemConfig.SessionKeyUserName), session.Get(systemConfig.SessionKeyUserRole))
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
	if sessUserId == nil || sessUserName == nil || sessUserRole == "" {
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

}

// ChangePassword 修改密码
func ChangePassword(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("ChangePassword error: %s\n%s", err, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()

	var json utils.ChangeUserPassword
	if err := c.ShouldBindJSON(&json); err != nil {
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
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "密码修改成功",
				"data": "",
			})
			return
		}
	}
}

// GetSystemLog 获取系统日志
func GetSystemLog(c *gin.Context) {

}
