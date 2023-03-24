package routers

import (
	"github.com/gin-gonic/gin"
	"xdpEngine/controllers"
)

func SystemLogApiGroup(router *gin.RouterGroup) {
	/**
	* /log
	 */
	// 系统信息配置
	router.GET("/search", controllers.SearchSysLog)

}
