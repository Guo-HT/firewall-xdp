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

// GetProtoRules 获取协议规则状态
func GetProtoRules(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("获取协议规则状态失败, %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []utils.ProtoRule{},
			})
			return
		}
	}()
	logger.Println("获取协议规则列表")
	rules := dpiEngine.ProtoRuleList
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取协议规则状态成功",
		"data": rules,
	})
	return
}

// SetProtoStatus 配置指定协议开关状态
func SetProtoStatus(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("配置协议开关状态失败, %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": dpiEngine.ProtoRuleList,
			})
			return
		}
	}()

	var json utils.ProtoStatusConf
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("SetProtoStatus: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []int{},
		})
		return
	} else {
		for index, rule := range dpiEngine.ProtoRuleList {
			if json.ProtoName == rule.ProtocolName {
				statusString := ""
				if json.Status {
					statusString = "开启"
				} else {
					statusString = "关闭"
				}
				logger.Printf("正在%s%s的分析", statusString, rule.ProtocolName)
				dpiEngine.ProtoRuleList[index].IsEnable = json.Status
			}
		}
		// 保存文件
		_ = dpiEngine.WriteProtoRuleFile()
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "配置" + json.ProtoName + "协议开关成功",
			"data": dpiEngine.ProtoRuleList,
		})
		return
	}
}

// ReloadProtoRule 重载协议规则文件
func ReloadProtoRule(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("ReloadProtoRule error : ", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": dpiEngine.ProtoRuleList,
			})
			return
		}
	}()
	dpiEngine.InitProtoRules() // 初始化协议规则列表
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "协议规则文件重载完成",
		"data": dpiEngine.ProtoRuleList,
	})
	return
}

// AddProtoRule 导入新的规则，保存配置
func AddProtoRule(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("AddProtoRule error : ", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": dpiEngine.ProtoRuleList,
			})
			return
		}
	}()

	var json utils.ProtoRule
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("AddProtoRule: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []utils.ProtoRule{},
		})
		return
	} else {
		logger.Printf("正在添加协议规则...")
		// 加入协议规则列表
		dpiEngine.ProtoRuleList = append(dpiEngine.ProtoRuleList, json)
		// 写入配置文件
		_ = dpiEngine.WriteProtoRuleFile() // 协议规则写入文件
		dpiEngine.InitProtoRules()         // 初始化协议规则列表
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "协议规则添加完成",
			"data": dpiEngine.ProtoRuleList,
		})
		return
	}
}
