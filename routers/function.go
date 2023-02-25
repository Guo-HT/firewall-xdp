package routers

import (
	"github.com/gin-gonic/gin"
	"xdpEngine/controllers"
)

func ProtoApiGroup(router *gin.RouterGroup) {
	/*
	 * /func
	 */
	ProtoEngineRG := router.Group("/proto")
	{
		ProtoEngineRG.POST("/start", controllers.StartProtoEngine) // 开启协议分析
		ProtoEngineRG.POST("/stop", controllers.StopProtoEngine)   // 关闭协议分析

		ProtoEngineRG.GET("/rules", controllers.GetProtoRules)          // 获取协议列表
		ProtoEngineRG.POST("/rules", controllers.SetProtoStatus)        // 配置协议列表
		ProtoEngineRG.POST("/reloadRules", controllers.ReloadProtoRule) // 配置协议列表
		ProtoEngineRG.POST("/addRules", controllers.AddProtoRule)       // 配置协议列表
	}
}
