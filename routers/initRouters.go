package routers

import (
	"github.com/gin-gonic/gin"
)

func InitRouters(engine *gin.Engine) {
	XdpApiGroup := engine.Group("/xdp")
	{
		BlackApiGroup(XdpApiGroup)
		WhiteApiGroup(XdpApiGroup)
	}

	FuncApiGroup := engine.Group("/func")
	{
		ProtoApiGroup(FuncApiGroup)
	}

}
