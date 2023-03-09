package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"sort"
	"strconv"
	"xdpEngine/systemConfig"
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
	//sort.Ints(whitePortList) // 先排序，后分页
	sort.SliceStable(whitePortList, func(i, j int) bool {
		if whitePortList[i].Hit > whitePortList[j].Hit {
			return true
		}
		return false
	})
	data, pNo, pSize := utils.IntIntStructListLimit(whitePortList, pageNo, pageSize)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Port白名单获取成功",
		"data": gin.H{
			"page_no":   pNo,
			"page_size": pSize,
			"total":     len(whitePortList),
			"data":      data,
		},
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
		json.WhitePortList = append(json.WhitePortList, systemConfig.ServerPort) // 默认将本服务端口写入白名单，保证服务正常
		xdp.IfaceXdpDict[json.Iface].Lock.Lock()                                 // 上写锁
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
	pageNoStr := c.Query("page_no")
	pageSizeStr := c.Query("page_size")
	pageNo, err := strconv.Atoi(pageNoStr)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		errlog.Println("GetWhiteIP error: 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []int{},
		})
		return
	}
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
	//sort.Sort(whiteIpList)
	sort.SliceStable(whiteIpList, func(i, j int) bool {
		if whiteIpList[i].Hit > whiteIpList[j].Hit {
			return true
		}
		return false
	})
	data, pNo, pSize := utils.StringIntStructListLimit(whiteIpList, pageNo, pageSize)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "IP白名单获取成功",
		"data": gin.H{
			"page_no":   pNo,
			"page_size": pSize,
			"total":     len(whiteIpList),
			"data":      data,
		},
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
