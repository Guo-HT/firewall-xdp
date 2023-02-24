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
		ProtoEngineRG.POST("/start", controllers.StartProtoEngine)
		ProtoEngineRG.POST("/stop", controllers.StopProtoEngine)
	}
}
