package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"sort"
	"strconv"
	"xdpEngine/utils"
	"xdpEngine/xdp"
)

// ******************** Port 操作 **************************

// GetBlackPort 获取Port黑名单
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
	pageNoStr := c.Query("page_no")
	pageSizeStr := c.Query("page_size")
	pageNo, err := strconv.Atoi(pageNoStr)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		errlog.Println("GetBlackPort error: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []int{},
		})
		return
	}
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
	//sort.Ints(blackPortList) // 先排序，后分页
	sort.SliceStable(blackPortList, func(i, j int) bool {
		if blackPortList[i].Hit > blackPortList[j].Hit {
			return true
		}
		return false
	})
	data, pNo, pSize := utils.IntIntStructListLimit(blackPortList, pageNo, pageSize)
	//logger.Println(blackPortList)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Port黑名单获取成功",
		"data": gin.H{
			"page_no":   pNo,
			"page_size": pSize,
			"total":     len(blackPortList),
			"data":      data,
		},
	})
	return
}

// SetBlackPort 配置Port黑名单
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

// DelBlackPort 删除Port黑名单
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

// ********************* IP 操作 ***************************

// GetBlackIP 删除IP黑名单
func GetBlackIP(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("GetBlackIP: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []string{},
			})
		}
	}()
	iface := c.Query("iface")
	pageNoStr := c.Query("page_no")
	pageSizeStr := c.Query("page_size")
	pageNo, err := strconv.Atoi(pageNoStr)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		errlog.Println("GetBlackIP error: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []int{},
		})
		return
	}
	blackIpList, err := xdp.GetAllBlackIpMap(iface)
	if err != nil {
		errlog.Println("IP黑名单获取失败,", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "IP黑名单获取失败",
			"data": []string{},
		})
		return
	}
	//sort.Strings(blackIpList)
	sort.SliceStable(blackIpList, func(i, j int) bool {
		if blackIpList[i].Hit > blackIpList[j].Hit {
			return true
		}
		return false
	})
	data, pNo, pSize := utils.StringIntStructListLimit(blackIpList, pageNo, pageSize)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "IP黑名单获取成功",
		"data": gin.H{
			"page_no":   pNo,
			"page_size": pSize,
			"total":     len(blackIpList),
			"data":      data,
		},
	})
	return
}

// SetBlackIP 配置IP黑名单
func SetBlackIP(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("SetBlackIP: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []string{},
			})
		}
	}()
	var json utils.BlackIpStruct
	if err := c.ShouldBindJSON(&json); err != nil || !utils.IsIpListRight(json.BlackIpList) {
		errlog.Println("SetBlackIP: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []string{},
		})
		return
	} else {
		xdp.IfaceXdpDict[json.Iface].Lock.Lock()
		xdp.IfaceXdpDict[json.Iface].BlackIpList = utils.AppendIPListDeduplicate(xdp.IfaceXdpDict[json.Iface].BlackIpList, json.BlackIpList)
		xdp.IfaceXdpDict[json.Iface].Lock.Unlock()
		err := xdp.InsertBlackIpMap(xdp.IfaceXdpDict[json.Iface].BlackIpList, json.Iface)
		if err != nil {
			errlog.Println("InsertBlackIpMap,", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "IP黑名单添加失败",
				"data": []string{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "IP黑名单添加成功",
			"data": xdp.IfaceXdpDict[json.Iface].BlackIpList,
		})
		return
	}
}

// DelBlackIP 获取IP黑名单
func DelBlackIP(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			errlog.Printf("DelBlackIP: %s \n %s", e, debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
				"data": []string{},
			})
		}
	}()

	var json utils.BlackIpStruct
	if err := c.ShouldBindJSON(&json); err != nil || !utils.IsIpListRight(json.BlackIpList) {
		errlog.Println("DelBlackIP: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []string{},
		})
		return
	} else {
		xdp.IfaceXdpDict[json.Iface].Lock.Lock()
		xdp.IfaceXdpDict[json.Iface].BlackIpList = utils.DeleteIpList(xdp.IfaceXdpDict[json.Iface].BlackIpList, json.BlackIpList)
		xdp.IfaceXdpDict[json.Iface].Lock.Unlock()

		err := xdp.DeleteBlackIpMap(json.BlackIpList, json.Iface)
		if err != nil {
			errlog.Println("DeleteBlackIpMap error: ", err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "IP黑名单删除失败",
				"data": []int{},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "IP黑名单删除成功",
			"data": xdp.IfaceXdpDict[json.Iface].BlackIpList,
		})
	}
}
