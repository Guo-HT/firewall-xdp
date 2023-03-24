package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"time"
	"xdpEngine/db"
	"xdpEngine/utils"
)

// StatusOverview 数据概览页
func StatusOverview(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("StatusOverview error: %s, %s", err, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
			return
		}
	}()
	//logger.Println("正在获取系统概览数据")
	ChanCpuLoad := make(chan float64)
	ChanSpeedIO := make(chan []utils.IOSpeed)
	// 系统时间
	serverTime := time.Now().Unix()
	// 运行时间
	serverRuntime := db.GetServerRuntime()
	// CPU占用
	go utils.GetCpuPercent(ChanCpuLoad)
	// Mem占用
	memPercent := utils.GetMemPercent()
	// 磁盘占用
	diskPercent := utils.GetDiskPercent()
	// 温度
	temperature := utils.GetCpuTemperature()
	// 网络IO
	go utils.GetAllNetcardIOSpeed(ChanSpeedIO)

	// 获取数据
	cpuLoad := <-ChanCpuLoad
	speedIO := <-ChanSpeedIO

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "概览数据获取成功",
		"data": gin.H{
			"server_time":    serverTime,
			"system_runtime": serverRuntime,
			"cpu_percent":    cpuLoad,
			"mem_percent":    memPercent,
			"disk_percent":   diskPercent,
			"temperature":    temperature,
			"speed_io":       speedIO,
		},
	})
	return
}

// SetSystemBanner 配置系统名称、图标
func SetSystemBanner(c *gin.Context) {
	username, _ := c.Get("username")
	defer func() {
		if err := recover(); err != nil {
			errlog.Printf("SetSystemBanner error: %s, %s", err, debug.Stack())
			db.SetSystemLog(c.ClientIP(), username.(string), "配置系统信息", false)
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
			return
		}
	}()

	var json utils.SystemSetting
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("SetSystemBanner: 请求参数错误")
		db.SetSystemLog(c.ClientIP(), username.(string), "配置系统信息", false)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误," + err.Error(),
			"data": "",
		})
		return
	} else {
		newSet, err := db.SetSystemSetting(json.Title, json.Icon)
		if err != nil {
			panic(err.Error())
		}
		db.SetSystemLog(c.ClientIP(), username.(string), "配置系统信息", true)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "系统信息配置成功",
			"data": gin.H{
				"title": newSet.Title,
				"icon":  newSet.Icon,
			},
		})
		return
	}
}

// GetSystemBanner 获取系统名称、图标
func GetSystemBanner(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("GetSystemBanner error: %s, %s", err, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": "",
			})
			return
		}
	}()
	sysSet := db.GetSystemSetting()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "系统信息获取成功",
		"data": gin.H{
			"title": sysSet.Title,
			"icon":  sysSet.Icon,
		},
	})
	return

}
