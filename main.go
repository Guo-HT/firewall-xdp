package main

import (
	"flag"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
	_ "xdpEngine/db"
	"xdpEngine/dpiEngine"
	"xdpEngine/routers"
	"xdpEngine/systemConfig"
	"xdpEngine/xdp"
)

var iface = flag.String("i", systemConfig.DefaultIface, "绑定的网卡名, eg: ens33/\"ens33 ens34\"")
var runMode = flag.String("m", "debug", "是否调试模式, debug->debug模式; release->release模式; test->test模式")

func main() {
	defer systemConfig.SayBye()
	flag.Parse() // 解析命令行参数

	systemConfig.RunMode = *runMode
	if *runMode == "release" {
		gin.SetMode(gin.ReleaseMode) // 生产模式
	} else if *runMode == "debug" {
		gin.SetMode(gin.DebugMode) // 调试模式
	} else if *runMode == "test" {
		gin.SetMode(gin.TestMode) // 测试模式
	} else {
		systemConfig.Errlog.Fatalln("-m, 参数错误, 退出...")
	}

	for _, iface := range strings.Split(*iface, " ") {
		xdp.InitEBpfMap(iface) // 获取ebpf maps
		dpiEngine.StartProtoEngine()
	}
	go xdp.ListenExit() // 监听退出信号

	engine := gin.Default()
	engine.Use(sessions.Sessions(systemConfig.SessionName, systemConfig.SessStore))
	routers.InitRouters(engine)

	go func() {
		if err := engine.Run(":" + systemConfig.ServerPortStr); err != nil {
			systemConfig.Errlog.Println("Gin start error:", err)
			xdp.StopAllXdpEngine()
			os.Exit(-1)
		}
	}()

	select {
	case <-systemConfig.CtrlC:
		xdp.StopAllXdpEngine()
		systemConfig.Logger.Println("再见!")
	}
}
