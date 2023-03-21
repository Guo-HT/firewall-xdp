package routers

import (
	"github.com/gin-gonic/gin"
	"xdpEngine/controllers"
)

func EngineStautsApiGroup(router *gin.RouterGroup) {
	/**
	* /status
	 */

	// 数据概览
	router.GET("/overview", controllers.StatusOverview)

	// 系统信息配置
	StatusRG := router.Group("/setting")
	{
		StatusRG.POST("/systemTitle", controllers.SetSystemBanner) // 配置系统显示信息
		StatusRG.GET("/systemTitle", controllers.GetSystemBanner)  // 获取系统显示信息
	}
}
