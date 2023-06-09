package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"strconv"
	"xdpEngine/db"
	"xdpEngine/dpiEngine"
	"xdpEngine/systemConfig"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

// StartProtoEngine 开启所有网口的协议分析功能
func StartProtoEngine(c *gin.Context) {
	username, _ := c.Get("username")
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("开启协议分析功能失败: %s", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "开启协议阻断", false)
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
		xdp.IfaceXdpDict[iface].Lock.RLock()
		//go dpiEngine.GetPacketFromChannel(iface) // 启动消费者
		//go dpiEngine.PacketCapture(iface)        // 开始抓包
		dpiEngine.StartIfaceProtoEngine(iface)
		xdpObj.ProtoSwitch = true
		err := xdp.SetFunctionSwitch("proto", "start")
		if err != nil {
			errlog.Println("SetFunctionSwitch start 'proto' error:", err.Error())
			xdpObj.CancelP() // 结束开启的相关协程
			db.SetSystemLog(c.ClientIP(), username.(string), "开启协议阻断", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "网卡协议识别功能开启失败",
				"data": "",
			})
			return
		}
		xdpObj.SessionFlow = make(map[string]*utils.SessionTuple) // 初始化会话流表
		xdpObj.CtxP, xdpObj.CancelP = context.WithCancel(context.Background())
		xdp.IfaceXdpDict[iface].Lock.RUnlock()
	}
	systemConfig.ProtoEngineStatus = true
	db.SetSystemLog(c.ClientIP(), username.(string), "开启协议阻断", true)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "网卡协议识别功能开启成功",
		"data": "",
	})
	return
}

// StopProtoEngine 关闭所有网口的协议分析功能
func StopProtoEngine(c *gin.Context) {
	username, _ := c.Get("username")
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("关闭协议分析功能失败: %s", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "关闭协议阻断", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
		}
	}()
	for iface, xdpObj := range xdp.IfaceXdpDict {
		logger.Printf("[%s]正在关闭协议分析功能", iface)
		xdp.IfaceXdpDict[iface].Lock.RLock()
		xdpObj.ProtoSwitch = false
		err := xdp.SetFunctionSwitch("proto", "stop")
		if err != nil {
			errlog.Println("SetFunctionSwitch stop 'proto' error:", err.Error())
			xdpObj.CancelP() // 结束开启的相关协程
			db.SetSystemLog(c.ClientIP(), username.(string), "关闭协议阻断", false)
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
			xdpObj.ProtoPoolChannel[i] = make(chan utils.FiveTuple, systemConfig.DefaultChanLength)
		}
		xdpObj.SessionFlow = make(map[string]*utils.SessionTuple) // 清空会话流表
		_ = xdp.UpdateProtoIpPortMap()

		xdpObj.CancelP() // 结束开启的相关协程
		xdp.IfaceXdpDict[iface].Lock.RUnlock()
	}
	systemConfig.ProtoEngineStatus = false
	db.SetSystemLog(c.ClientIP(), username.(string), "关闭协议阻断", true)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "网卡协议识别功能关闭成功",
		"data": "",
	})
	return
}

// GetProtoEngineStatus 获取协议分析引擎开关状态
func GetProtoEngineStatus(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("获取协议分析引擎开关状态, %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取协议识别引擎状态成功",
		"data": systemConfig.ProtoEngineStatus,
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
	username, _ := c.Get("username")
	var json utils.ProtoStatusConf
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("配置协议开关状态失败, %s", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "配置 "+json.ProtoName+" 识别开关为 "+strconv.FormatBool(json.Status), false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": dpiEngine.ProtoRuleList,
			})
			return
		}
	}()

	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("SetProtoStatus: 请求参数错误")
		db.SetSystemLog(c.ClientIP(), username.(string), "配置 "+json.ProtoName+" 识别开关为 "+strconv.FormatBool(json.Status), false)
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
		db.SetSystemLog(c.ClientIP(), username.(string), "配置 "+json.ProtoName+" 识别开关为 "+strconv.FormatBool(json.Status), true)
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
	username, _ := c.Get("username")
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("ReloadProtoRule error : ", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "重载协议规则", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": dpiEngine.ProtoRuleList,
			})
			return
		}
	}()
	dpiEngine.InitProtoRules() // 初始化协议规则列表
	db.SetSystemLog(c.ClientIP(), username.(string), "重载协议规则", true)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "协议规则文件重载完成",
		"data": dpiEngine.ProtoRuleList,
	})
	return
}

// AddProtoRule 导入新的规则，保存配置
func AddProtoRule(c *gin.Context) {
	username, _ := c.Get("username")
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("AddProtoRule error : ", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "添加协议规则", false)
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
		db.SetSystemLog(c.ClientIP(), username.(string), "添加协议规则", false)
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
		dpiEngine.ProtoRuleList = utils.ResetRulesId(dpiEngine.ProtoRuleList)
		// 写入配置文件
		_ = dpiEngine.WriteProtoRuleFile() // 协议规则写入文件
		dpiEngine.InitProtoRules()         // 初始化协议规则列表
		db.SetSystemLog(c.ClientIP(), username.(string), "添加协议规则", true)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "协议规则添加完成",
			"data": dpiEngine.ProtoRuleList,
		})
		return
	}
}

// GetProtoIpPort 获取所有IP-Port策略
func GetProtoIpPort(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("获取协议阻断IP-Port策略失败, %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []utils.IpPort{},
			})
			return
		}
	}()
	logger.Println("获取协议阻断IP-Port策略")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取协议规则状态成功",
		"data": xdp.GetAllProtoIpPortMap(),
	})
	return
}

// DelProtoRule 删除指定协议规则
func DelProtoRule(c *gin.Context) {
	username, _ := c.Get("username")
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("DelProtoRule error : ", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "删除协议规则", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": dpiEngine.ProtoRuleList,
			})
			return
		}
	}()

	var json utils.ProtoId
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("DelProtoRule: 请求参数错误")
		db.SetSystemLog(c.ClientIP(), username.(string), "删除协议规则", false)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []utils.ProtoRule{},
		})
		return
	} else {
		dpiEngine.ProtoRuleList = utils.DeleteRuleById(dpiEngine.ProtoRuleList, json.Id)
		// 写入配置文件
		_ = dpiEngine.WriteProtoRuleFile() // 协议规则写入文件
		dpiEngine.InitProtoRules()         // 初始化协议规则列表
		db.SetSystemLog(c.ClientIP(), username.(string), "删除协议规则", true)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "协议规则删除完成",
			"data": dpiEngine.ProtoRuleList,
		})
		return
	}

}
