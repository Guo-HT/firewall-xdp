package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"xdpEngine/dpiEngine"
	"xdpEngine/routers"
	"xdpEngine/systemConfig"
	"xdpEngine/xdp"
)

var iface = flag.String("i", systemConfig.DefaultIface, "绑定的网卡名, eg: ens33")
var runMode = flag.String("M", "debug", "是否调试模式, debug->debug模式; release->release模式; test->test模式")

func main() {
	flag.Parse() // 解析命令行参数

	systemConfig.RunMode = *runMode
	if *runMode == "release" {
		gin.SetMode(gin.ReleaseMode) // 生产模式
	} else if *runMode == "debug" {
		gin.SetMode(gin.DebugMode) // 调试模式
	} else if *runMode == "test" {
		gin.SetMode(gin.TestMode) // 测试模式
	} else {
		systemConfig.Errlog.Fatalln("-d, 参数错误, 退出...")
	}

	xdp.InitEBpfMap(*iface) // 获取ebpf maps
	go xdp.ListenExit()     // 监听退出信号
	dpiEngine.StartProtoEngine()

	engine := gin.Default()
	routers.InitRouters(engine)

	if err := engine.Run(":" + systemConfig.ServerPortStr); err != nil {
		systemConfig.Errlog.Println("Gin start error:", err)
	}
	systemConfig.Logger.Println("============================")
}
