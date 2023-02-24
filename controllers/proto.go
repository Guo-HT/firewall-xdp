package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"xdpEngine/dpiEngine"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

// StartProtoEngine 开启所有网口的协议分析功能
func StartProtoEngine(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("开启协议分析功能失败: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
		}
	}()
	for iface, xdpObj := range xdp.IfaceXdpDict {
		if xdpObj.ProtoSwitch {
			logger.Printf("[%s]协议分析功能已开启，跳过", iface)
			continue
		}
		logger.Printf("[%s]正在开启协议分析功能", iface)
		go dpiEngine.GetPacketFromChannel(iface) // 启动消费者
		go dpiEngine.PacketCapture(iface)        // 开始抓包
		xdpObj.ProtoSwitch = true
		err := xdp.SetFunctionSwitch("proto", "start")
		if err != nil {
			errlog.Println("SetFunctionSwitch start 'proto' error:", err.Error())
			xdpObj.CancelP() // 结束开启的相关协程
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "网卡协议识别功能开启失败",
				"data": "",
			})
			return
		}
		xdpObj.SessionFlow = make(map[string]*utils.SessionTuple) // 初始化会话流表
		xdpObj.CtxP, xdpObj.CancelP = context.WithCancel(context.Background())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "网卡协议识别功能开启成功",
		"data": "",
	})
	return
}

// StopProtoEngine 关闭所有网口的协议分析功能
func StopProtoEngine(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("关闭协议分析功能失败: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
		}
	}()
	for iface, xdpObj := range xdp.IfaceXdpDict {
		logger.Printf("[%s]正在关闭协议分析功能", iface)
		err := xdp.SetFunctionSwitch("proto", "stop")
		if err != nil {
			errlog.Println("SetFunctionSwitch stop 'proto' error:", err.Error())
			xdpObj.CancelP() // 结束开启的相关协程
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "网卡协议识别功能开启失败",
				"data": "",
			})
			return
		}
		// 清空当前的缓冲池
		xdpObj.ProtoPoolChannel = make([]chan utils.FiveTuple, systemConfig.DefaultChanNum)
		for i := 0; i < systemConfig.DefaultChanNum; i++ {
			xdpObj.ProtoPoolChannel[i] = make(chan utils.FiveTuple, 10000)
		}
		xdpObj.SessionFlow = make(map[string]*utils.SessionTuple) // 清空会话流表
		xdpObj.CancelP()                                          // 结束开启的相关协程
		xdpObj.ProtoSwitch = false
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "网卡协议识别功能关闭成功",
		"data": "",
	})
	return
}
