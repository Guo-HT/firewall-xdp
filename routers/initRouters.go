package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xdpEngine/controllers"
)

func InitRouters(engine *gin.Engine) {
	// 模板文件
	engine.LoadHTMLFiles("web/index.tmpl", "web/login.tmpl")
	// 静态文件
	engine.Static("/web/static", "./web/static")

	// 静态文件路由
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "web/index.tmpl", gin.H{})
	})
	engine.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "web/login.tmpl", gin.H{})
	})

	/********************** API接口 ************************/
	XdpApiGroup := engine.Group("/xdp")
	XdpApiGroup.Use(controllers.LoginRequireMiddleware())
	{
		BlackApiGroup(XdpApiGroup)
		WhiteApiGroup(XdpApiGroup)
	}

	FuncApiGroup := engine.Group("/func")
	FuncApiGroup.Use(controllers.LoginRequireMiddleware())
	{
		ProtoApiGroup(FuncApiGroup)
	}

	NetcardApiGroup := engine.Group("/iface")
	NetcardApiGroup.Use(controllers.LoginRequireMiddleware())
	{
		IfaceApiGroup(NetcardApiGroup)
	}

	UserApiGroup := engine.Group("/user")
	{
		UserOptApiGroup(UserApiGroup)
	}

	StatusApiGroup := engine.Group("/status")
	//StatusApiGroup.Use(controllers.LoginRequireMiddleware())
	{
		EngineStautsApiGroup(StatusApiGroup)
	}

	LogApiGroup := engine.Group("/log")
	LogApiGroup.Use(controllers.LoginRequireMiddleware())
	{
		SystemLogApiGroup(LogApiGroup)
	}

}
