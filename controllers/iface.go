package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"xdpEngine/db"
	"xdpEngine/dpiEngine"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

// AttachNewIface 挂载引擎到新的网卡
func AttachNewIface(c *gin.Context) {
	username, _ := c.Get("username")
	var json utils.IfaceStruct
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("AttachNewIface error: %s", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "挂载引擎到网卡"+json.Iface, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()

	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("AttachNewIface: 请求参数错误")
		db.SetSystemLog(c.ClientIP(), username.(string), "挂载引擎到网卡"+json.Iface, false)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误," + err.Error(),
			"data": "",
		})
		return
	} else {
		if _, ok := xdp.IfaceXdpDict[json.Iface]; !ok {
			// 未绑定网卡，可以挂载
			logger.Printf("[%s]正在挂载到新的网卡", json.Iface)
			xdp.DetachXdp(json.Iface)
			xdp.InitEBpfMap(json.Iface)
			if xdp.IfaceXdpDict[json.Iface].ProtoSwitch {
				logger.Printf("正在开启[%s]的分析功能...", json.Iface)
				dpiEngine.StartIfaceProtoEngine(json.Iface)
			}
			db.SetSystemLog(c.ClientIP(), username.(string), "挂载引擎到网卡"+json.Iface, true)
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "引擎挂载成功",
				"data": "",
			})
			return
		} else {
			errlog.Println("AttachNewIface: 网卡重复绑定")
			db.SetSystemLog(c.ClientIP(), username.(string), "挂载引擎到网卡"+json.Iface, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "网卡已挂载引擎",
				"data": "",
			})
			return
		}
	}
}

// DetachNewIface 卸载引擎
func DetachNewIface(c *gin.Context) {
	username, _ := c.Get("username")
	var json utils.IfaceStruct
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("DetachNewIface error: %s", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "从网卡 "+json.Iface+" 卸载引擎到网卡", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("DetachNewIface: 请求参数错误")
		db.SetSystemLog(c.ClientIP(), username.(string), "从网卡 "+json.Iface+" 卸载引擎到网卡", false)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误," + err.Error(),
			"data": "",
		})
		return
	} else {
		if xdpObj, ok := xdp.IfaceXdpDict[json.Iface]; ok {
			// 绑定网卡，可以挂载
			logger.Printf("[%s]正在卸载引擎", json.Iface)
			xdp.StopXdpEngine(json.Iface, xdpObj)
			db.SetSystemLog(c.ClientIP(), username.(string), "从网卡 "+json.Iface+" 卸载引擎到网卡", true)
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "引擎卸载成功",
				"data": "",
			})
			return
		} else {
			// 未绑定网卡，不卸载
			errlog.Println("DetachNewIface: 网卡未挂载引擎")
			db.SetSystemLog(c.ClientIP(), username.(string), "从网卡 "+json.Iface+" 卸载引擎到网卡", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "网卡未挂载引擎",
				"data": "",
			})
			return
		}
	}
}

// EngineStatus 获取指定网卡上引擎状态
func EngineStatus(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("EngineStatus: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()
	iface := c.Query("iface")
	logger.Printf("[%s]正在获取引擎状态...", iface)

	errStatus := false
	portWhite, err := xdp.GetAllWhitePortMap(iface)
	if err != nil {
		errlog.Println("portWhite error, ", err.Error())
		errStatus = true
	}
	portBlack, err := xdp.GetAllBlackPortMap(iface)
	if err != nil {
		errlog.Println("portBlack error, ", err.Error())
		errStatus = true
	}
	ipWhite, err := xdp.GetAllWhiteIpMap(iface)
	if err != nil {
		errlog.Println("ipWhite error, ", err.Error())
		errStatus = true
	}
	ipBlack, err := xdp.GetAllBlackIpMap(iface)
	if err != nil {
		errlog.Println("ipBlack error, ", err.Error())
		errStatus = true
	}
	protoIpPort := xdp.GetAllProtoIpPortMap()

	protoCode, err := xdp.GetFunctionSwitch("proto", iface)
	if err != nil {
		errlog.Println("protoStatus error, ", err.Error())
		errStatus = true
	}
	if errStatus {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求错误",
			"data": "",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "网卡引擎状态获取成功",
			"data": gin.H{
				"iface":         iface,
				"port_white":    portWhite,
				"port_black":    portBlack,
				"ip_white":      ipWhite,
				"ip_black":      ipBlack,
				"proto_switch":  utils.ConvertProtoCode2Status(protoCode),
				"proto":         dpiEngine.GetStartingProto(),
				"proto_ip_port": protoIpPort,
			},
		})
		return
	}
}

// GetEngineList 获取引擎挂载列表
func GetEngineList(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("EngineStatus: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()
	engineList := make([]string, len(xdp.IfaceXdpDict))
	i := 0
	for iface := range xdp.IfaceXdpDict {
		engineList[i] = iface
		i++
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "引擎挂载列表获取成功",
		"data": engineList,
	})
	return

}

// GetIfaceList http获取所有可用网卡信息
func GetIfaceList(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("GetIfaceList: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()
	logger.Println("正在获取网卡列表")
	netcardList := utils.GetAllNetcard()
	for index, netcard := range netcardList {
		if _, ok := xdp.IfaceXdpDict[netcard.NetcardName]; ok {
			// 已挂载引擎
			netcardList[index].EngineStatus = true
		} else {
			netcardList[index].EngineStatus = false
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取可用网卡成功",
		"data": netcardList,
	})
	return
}

// StartIface up网卡
func StartIface(c *gin.Context) {
	username, _ := c.Get("username")
	var json utils.IfaceStruct
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("StartIface: %s", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "启动网卡 "+json.Iface, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
			return
		}
	}()

	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("StartIface: 请求参数错误")
		db.SetSystemLog(c.ClientIP(), username.(string), "启动网卡 "+json.Iface, false)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误," + err.Error(),
			"data": "",
		})
		return
	} else {
		logger.Println("正在开启网卡:", json.Iface)
		utils.UpIfaceState(json.Iface)
		db.SetSystemLog(c.ClientIP(), username.(string), "启动网卡 "+json.Iface, true)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "网卡开启成功",
			"data": "",
		})
		return
	}
}

// StopIface down网卡
func StopIface(c *gin.Context) {
	username, _ := c.Get("username")
	var json utils.IfaceStruct
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("StopIface: %s", debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "关闭网卡 "+json.Iface, false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
			return
		}
	}()

	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("StopIface: 请求参数错误")
		db.SetSystemLog(c.ClientIP(), username.(string), "关闭网卡 "+json.Iface, false)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误," + err.Error(),
			"data": "",
		})
		return
	} else {
		//fmt.Println(c.ClientIP())
		logger.Println("正在关闭网卡:", json.Iface)
		utils.DownIfaceState(json.Iface)
		db.SetSystemLog(c.ClientIP(), username.(string), "关闭网卡 "+json.Iface, true)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "网卡关闭成功",
			"data": "",
		})
		return
	}
}
