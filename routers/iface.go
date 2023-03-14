package routers

import (
	"github.com/gin-gonic/gin"
	"xdpEngine/controllers"
)

func IfaceApiGroup(router *gin.RouterGroup) {
	/*
	 * /iface
	 */
	NetcardEngineRG := router.Group("/engine")
	{
		NetcardEngineRG.POST("/attach", controllers.AttachNewIface)
		NetcardEngineRG.POST("/detach", controllers.DetachNewIface)
		NetcardEngineRG.GET("/getStatus", controllers.EngineStatus)
		NetcardEngineRG.GET("/getEngineList", controllers.GetEngineList)
		NetcardEngineRG.GET("/getIfaceList", controllers.GetIfaceList)
	}
}
