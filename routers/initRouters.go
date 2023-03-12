package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouters(engine *gin.Engine) {
	engine.LoadHTMLFiles("web/index.tmpl", "web/login.tmpl")
	engine.Static("/web/static", "./web/static")

	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "web/index.tmpl", gin.H{})
	})
	engine.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "web/login.tmpl", gin.H{})
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
