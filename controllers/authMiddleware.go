package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"xdpEngine/systemConfig"
)

func LoginRequireMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessUserId := session.Get(systemConfig.SessionKeyUserId)
		sessUserName := session.Get(systemConfig.SessionKeyUserName)
		sessUserRole := session.Get(systemConfig.SessionKeyUserRole)
		if sessUserId == nil || sessUserName == nil || sessUserRole == nil {
			// 未登录
			errlog.Println(c.Request.RequestURI, ": 当前未登录, 拒绝请求")
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "当前未登录",
				"data": "",
			})
			c.Abort()
			return
		} else {
			// 已登录
			c.Set("username", sessUserName)
			c.Set("userId", sessUserName)
			c.Set("userRole", sessUserRole)
		}
	}
}
