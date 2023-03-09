package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouters(engine *gin.Engine) {
	engine.LoadHTMLFiles("web/index.tmpl")
	engine.Static("/web/static", "./web/static")

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "web/index.tmpl", gin.H{})
	})

	XdpApiGroup := engine.Group("/xdp")
	{
		BlackApiGroup(XdpApiGroup)
		WhiteApiGroup(XdpApiGroup)
	}

	FuncApiGroup := engine.Group("/func")
	{
		ProtoApiGroup(FuncApiGroup)
	}

	NetcardApiGroup := engine.Group("/iface")
	{
		IfaceApiGroup(NetcardApiGroup)
	}

	UserApiGroup := engine.Group("/user")
	{
		UserOptApiGroup(UserApiGroup)
	}

}
