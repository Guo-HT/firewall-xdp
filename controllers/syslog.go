package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"strconv"
	"xdpEngine/db"
	"xdpEngine/utils"
)

func SearchSysLog(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errlog.Println("SearchSysLog error: %s\n%s", err.(string), debug.Stack())
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "查询系统操作日志失败",
				"data": []utils.SystemLogRsp{},
			})
			return
		}
	}()

	startTimeStr := c.DefaultQuery("start_time", "0")
	endTimeStr := c.DefaultQuery("end_time", "9999999999")
	search := c.Query("search")
	pageNoStr := c.Query("page_no")
	pageSizeStr := c.Query("page_size")
	sort := c.DefaultQuery("sort", "desc")
	startTime, err := strconv.Atoi(startTimeStr)
	endTime, err := strconv.Atoi(endTimeStr)
	pageNo, err := strconv.Atoi(pageNoStr)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || !(sort == "asc" || sort == "desc") {
		errlog.Println("SearchSysLog 请求参数错误")
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "请求参数错误",
			"data": []utils.SystemLogRsp{},
		})
		return
	}
	if startTime == 0 && endTime == 0 {
		endTime = 9999999999
	}
	systemLog, total := db.SearchSystemLog(search, startTime, endTime, pageNo, pageSize, sort)
	sysLog := utils.SystemLogDbToWeb(systemLog)
	logger.Println("系统操作日志查询成功")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "系统操作日志查询成功",
		"data": gin.H{
			"page_size": pageSize,
			"page_no":   pageNo,
			"total":     total,
			"log":       sysLog,
		},
	})
	return

}
