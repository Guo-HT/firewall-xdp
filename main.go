package main

import (
	"github.com/gin-gonic/gin"
	"xdpEngine/routers"
	"xdpEngine/systemConfig"
	"xdpEngine/xdp"
)

func main() {
	systemConfig.PrintBanner()
	xdp.InitEBpfMap()
	go xdp.ListenExit()
	engine := gin.Default()
	// 获取ebpf maps

	routers.InitRouters(engine)

	if err := engine.Run(":1888"); err != nil {
		systemConfig.Logger.Println("Gin start error:", err)
	}
}
