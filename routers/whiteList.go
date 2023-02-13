package routers

import (
	"github.com/gin-gonic/gin"
	"xdpEngine/controllers"
)

func WhiteApiGroup(router *gin.RouterGroup) {
	/*
	 * /api
	 */
	WhiteListRG := router.Group("/white")
	{
		WhiteListRG.GET("/getPort", controllers.GetWhitePort)
		WhiteListRG.GET("/getIP", controllers.GetWhiteIP)
		WhiteListRG.POST("/setPort", controllers.SetWhitePort)
		WhiteListRG.POST("/setIP", controllers.SetWhiteIP)
	}
}
