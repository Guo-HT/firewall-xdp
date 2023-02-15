package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

func GetWhitePort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("GetWhitePort: %s", debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []int{},
			})
			return
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
	//logger.Println(whitePortList)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Port白名单获取成功",
		"data": whitePortList,
	})
	return

}

func GetWhiteIP(c *gin.Context) {

}

// SetWhitePort 配置Port白名单
func SetWhitePort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("SetWhitePort: %s", debug.Stack())
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
			"msg":  "请求参数错误," + err.Error(),
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

func SetWhiteIP(c *gin.Context) {

}

// DelWhitePort 删除Port白名单
func DelWhitePort(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("DelWhitePort: %s", debug.Stack())
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
			"msg":  "请求参数错误," + err.Error(),
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

func DelWhiteIP(c *gin.Context) {

}
