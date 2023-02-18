package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

// ******************** Port 操作 **************************

// GetWhitePort 获取Port白名单
func GetWhitePort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("GetWhitePort: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
		}
	}()
	iface := c.Query("iface")
	whitePortList, err := xdp.GetAllWhitePortMap(iface)
	if err != nil {
		errlog.Println("Port白名单获取失败,", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "Port白名单获取失败",
			"data": []int{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Port白名单获取成功",
		"data": whitePortList,
	})
	return
}

// SetWhitePort 配置Port白名单
func SetWhitePort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("SetWhitePort: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
		}
	}()

	var json utils.WhitePortStruct
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("SetWhitePort: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []int{},
		})
		return
	} else {
		xdp.IfaceXdpDict[json.Iface].Lock.Lock() // 上写锁
		xdp.IfaceXdpDict[json.Iface].WhitePortList = utils.AppendPortListDeduplicate(xdp.IfaceXdpDict[json.Iface].WhitePortList, json.WhitePortList)
		xdp.IfaceXdpDict[json.Iface].Lock.Unlock() // 解写锁

		err := xdp.InsertWhitePortMap(xdp.IfaceXdpDict[json.Iface].WhitePortList, json.Iface)
		if err != nil {
			errlog.Println("InsertWhitePortMap错误,", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "Port白名单添加失败",
				"data": []int{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "Port白名单添加成功",
			"data": xdp.IfaceXdpDict[json.Iface].WhitePortList,
		})
		return
	}
}

// DelWhitePort 删除Port白名单
func DelWhitePort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("DelWhitePort: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
		}
	}()

	var json utils.WhitePortStruct
	if err := c.ShouldBindJSON(&json); err != nil {
		errlog.Println("DelWhitePort: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []int{},
		})
		return
	} else {
		xdp.IfaceXdpDict[json.Iface].Lock.Lock() // 上写锁
		// 删除
		xdp.IfaceXdpDict[json.Iface].WhitePortList = utils.DeletePortList(xdp.IfaceXdpDict[json.Iface].WhitePortList, json.WhitePortList)
		xdp.IfaceXdpDict[json.Iface].Lock.Unlock() // 解写锁

		err := xdp.DeleteWhitePortMap(json.WhitePortList, json.Iface)
		if err != nil {
			errlog.Println("DeleteWhitePortMap错误,", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "Port白名单删除失败",
				"data": []int{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "Port白名单删除成功",
			"data": xdp.IfaceXdpDict[json.Iface].WhitePortList,
		})
		return
	}
}

// ********************* IP 操作 ***************************

// GetWhiteIP 获取IP白名单
func GetWhiteIP(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("GetWhiteIP: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []string{},
			})
		}
	}()
	iface := c.Query("iface")
	whiteIpList, err := xdp.GetAllWhiteIpMap(iface)
	if err != nil {
		errlog.Println("IP白名单获取失败,", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "IP白名单获取失败",
			"data": []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "IP白名单获取成功",
		"data": whiteIpList,
	})
	return
}

// SetWhiteIP 配置IP白名单
func SetWhiteIP(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("SetWhiteIP: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []string{},
			})
		}
	}()

	var json utils.WhiteIpStruct
	if err := c.ShouldBindJSON(&json); err != nil || !utils.IsIpListRight(json.WhiteIpList) {
		errlog.Println("SetWhiteIP: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []string{},
		})
		return
	} else {
		xdp.IfaceXdpDict[json.Iface].Lock.Lock()
		xdp.IfaceXdpDict[json.Iface].WhiteIpList = utils.AppendIPListDeduplicate(xdp.IfaceXdpDict[json.Iface].WhiteIpList, json.WhiteIpList)
		xdp.IfaceXdpDict[json.Iface].Lock.Unlock()
		err := xdp.InsertWhiteIpMap(xdp.IfaceXdpDict[json.Iface].WhiteIpList, json.Iface)
		if err != nil {
			errlog.Println("InsertWhiteIpMap错误,", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "IP白名单添加失败",
				"data": []string{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "IP白名单添加成功",
			"data": xdp.IfaceXdpDict[json.Iface].WhiteIpList,
		})
		return
	}
}

// DelWhiteIP 删除IP白名单
func DelWhiteIP(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("DelWhiteIP: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []string{},
			})
		}
	}()
	var json utils.WhiteIpStruct
	if err := c.ShouldBindJSON(&json); err != nil || !utils.IsIpListRight(json.WhiteIpList) {
		errlog.Println("DelWhiteIP: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []string{},
		})
		return
	} else {
		xdp.IfaceXdpDict[json.Iface].Lock.Lock()
		xdp.IfaceXdpDict[json.Iface].WhiteIpList = utils.DeleteIpList(xdp.IfaceXdpDict[json.Iface].WhiteIpList, json.WhiteIpList)
		xdp.IfaceXdpDict[json.Iface].Lock.Unlock()

		err := xdp.DeleteWhiteIpMap(json.WhiteIpList, json.Iface)
		if err != nil {
			errlog.Println("DeleteWhiteIpMap error: ", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "IP白名单删除失败",
				"data": []int{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "IP白名单删除成功",
			"data": xdp.IfaceXdpDict[json.Iface].WhiteIpList,
		})
	}
}
