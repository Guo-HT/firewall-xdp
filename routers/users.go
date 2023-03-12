package routers

import (
	"github.com/gin-gonic/gin"
	"xdpEngine/controllers"
)

func UserOptApiGroup(router *gin.RouterGroup) {
	/*
	 * /user
	 */
	UserEngineRG := router.Group("/status")
	{
		UserEngineRG.POST("/login", controllers.UserLogin)
		UserEngineRG.POST("/getToken", controllers.GetAccessToken)
		UserEngineRG.POST("/logout", controllers.UserLogout)
		UserEngineRG.GET("/info", controllers.UserInfo)
		UserEngineRG.POST("/addUser", controllers.AddUser)
		UserEngineRG.POST("/changePwd", controllers.ChangePassword)
	}

	OptLogRG := router.Group("/log")
	{
		OptLogRG.GET("/systemLog", controllers.GetSystemLog)
	}
}
