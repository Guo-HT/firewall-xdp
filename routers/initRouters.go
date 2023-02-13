package routers

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(engine *gin.Engine) {
	ApiGroup := engine.Group("/xdp")
	{
		BlackApiGroup(ApiGroup)
		WhiteApiGroup(ApiGroup)
	}
}
