package main

import (
	"github.com/gin-gonic/gin"
	"xdpEngine/routers"
	"xdpEngine/systemConfig"
	"xdpEngine/xdp"
)

func main() {
	systemConfig.PrintBanner()
	xdp.InitEBpfMap()   // 获取ebpf maps
	go xdp.ListenExit() // 监听退出信号
	engine := gin.Default()

	routers.InitRouters(engine)

	if err := engine.Run(":1888"); err != nil {
		systemConfig.Logger.Println("Gin start error:", err)
	}
}
