package routers

import "github.com/gin-gonic/gin"
import "xdpEngine/controllers"

func BlackApiGroup(router *gin.RouterGroup) {
	/*
	 * /api
	 */
	BlackListRG := router.Group("/black")
	{
		BlackListRG.GET("/getPort", controllers.GetBlackPort)
		BlackListRG.GET("/getIP", controllers.GetBlackIP)

		BlackListRG.POST("/setPort", controllers.SetBlackPort)
		BlackListRG.POST("/setIP", controllers.SetBlackIP)

		BlackListRG.POST("/delPort", controllers.DelBlackPort)
		BlackListRG.POST("/delIP", controllers.DelBlackIP)
	}
}
