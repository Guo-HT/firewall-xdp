package routers

import "github.com/gin-gonic/gin"

func UserOptApiGroup(router *gin.RouterGroup) {
	/*
	 * /user
	 */
	UserEngineRG := router.Group("/status")
	{
		UserEngineRG.POST("/login")
		UserEngineRG.POST("/logout")
		UserEngineRG.GET("/info")
	}

	OptLogRG := router.Group("/log")
	{
		OptLogRG.GET("/systemLog")
	}
}
