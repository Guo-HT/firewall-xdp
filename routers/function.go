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
		ProtoEngineRG.POST("/start", controllers.StartProtoEngine)     // 开启协议分析
		ProtoEngineRG.POST("/stop", controllers.StopProtoEngine)       // 关闭协议分析
		ProtoEngineRG.GET("/status", controllers.GetProtoEngineStatus) // 获取协议分析开关状态

		ProtoEngineRG.GET("/getProtoIpPort", controllers.GetProtoIpPort) // 新增协议规则

		ProtoEngineRG.GET("/rules", controllers.GetProtoRules)          // 获取协议列表
		ProtoEngineRG.POST("/rules", controllers.SetProtoStatus)        // 配置协议开关
		ProtoEngineRG.POST("/reloadRules", controllers.ReloadProtoRule) // 重载协议列表
		ProtoEngineRG.POST("/addRules", controllers.AddProtoRule)       // 新增协议规则
		ProtoEngineRG.POST("/delRules", controllers.DelProtoRule)       // 新增协议规则

	}
}
