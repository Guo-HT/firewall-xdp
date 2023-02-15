package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

func GetBlackPort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("GetBlackPort: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
		}
	}()
	iface := c.Query("iface")
	blackPortList, err := xdp.GetAllBlackPortMap(iface)
	if err != nil {
		errlog.Println("Port黑名单获取失败,", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "Port黑名单获取失败",
			"data": []int{},
		})
		return
	}
	//logger.Println(blackPortList)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Port黑名单获取成功",
		"data": blackPortList,
	})
	return
}

func GetBlackIP(c *gin.Context) {

}

func SetBlackPort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("SetBlackPort: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
		}
	}()

	var json utils.BlackPortStruct
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("SetBlackPort: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误," + err.Error(),
			"data": []int{},
		})
		return
	} else {
		xdp.IfaceXdpDict[json.Iface].Lock.Lock() // 上写锁
		xdp.IfaceXdpDict[json.Iface].BlackPortList = utils.AppendPortListDeduplicate(xdp.IfaceXdpDict[json.Iface].BlackPortList, json.BlackPortList)
		xdp.IfaceXdpDict[json.Iface].Lock.Unlock() // 解写锁

		err := xdp.InsertBlackPortMap(xdp.IfaceXdpDict[json.Iface].BlackPortList, json.Iface)
		if err != nil {
			errlog.Println("InsertBlackPortMap错误,", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "Port黑名单添加失败",
				"data": []int{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "Port黑名单添加成功",
			"data": xdp.IfaceXdpDict[json.Iface].BlackPortList,
		})
		return
	}
}

func SetBlackIP(c *gin.Context) {

}

func DelBlackPort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("DelBlackPort: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
		}
	}()

	var json utils.BlackPortStruct
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("DelBlackPort: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误," + err.Error(),
			"data": []int{},
		})
		return
	} else {
		xdp.IfaceXdpDict[json.Iface].Lock.Lock() // 上写锁
		// 删除
		xdp.IfaceXdpDict[json.Iface].BlackPortList = utils.DeletePortList(xdp.IfaceXdpDict[json.Iface].BlackPortList, json.BlackPortList)
		xdp.IfaceXdpDict[json.Iface].Lock.Unlock() // 解写锁

		err := xdp.DeleteBlackPortMap(json.BlackPortList, json.Iface)
		if err != nil {
			errlog.Println("DeleteBlackPortMap错误,", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "Port黑名单删除失败",
				"data": []int{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "Port黑名单删除成功",
			"data": xdp.IfaceXdpDict[json.Iface].BlackPortList,
		})
		return
	}
}

func DelBlackIP(c *gin.Context) {

}
