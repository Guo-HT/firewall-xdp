package db

import (
	"gorm.io/gorm"
	"log"
	"os"
	"time"
	"xdpEngine/systemConfig"
)

var (
	logger *log.Logger
	errlog *log.Logger
	CtrlC  chan os.Signal
	dbConn *gorm.DB
)

func init() {
	logger = systemConfig.Logger
	errlog = systemConfig.Errlog
	CtrlC = systemConfig.CtrlC
	dbConn = systemConfig.DB
	MigrateAll()
	DataInit()
	go AddRuntime()
}

// MigrateAll 自动迁移数据库
func MigrateAll() {
	logger.Println("正在迁移数据库...")
	err := dbConn.AutoMigrate(&User{})
	if err != nil {
		errlog.Fatalln("table [User] migrate error:", err.Error())
	}
	err = dbConn.AutoMigrate(&SystemLog{})
	if err != nil {
		errlog.Fatalln("table [SystemLog] migrate error:", err.Error())
	}
	err = dbConn.AutoMigrate(&SystemSetting{})
	if err != nil {
		errlog.Fatalln("table [SystemSetting] migrate error:", err.Error())
	}
}

// DataInit 初始化数据库
func DataInit() {
	// 1. 用户表
	var countUser int64
	dbConn.Model(&User{}).Count(&countUser)
	logger.Printf("当前有[ %d ]个用户", countUser)
	if countUser == 0 {
		dbConn.Create(&User{
			ID:                  0,
			UserName:            "admin",
			Password:            "admin",
			Email:               "",
			Role:                0,
			IsLocked:            false,
			ErrorTimes:          0,
			IsChangedDefualtPwd: false,
		})
	}
	// 2. 系统信息表
	var countServerInfo int64
	dbConn.Model(&SystemSetting{}).Count(&countServerInfo)
	if countServerInfo == 0 {
		logger.Println("当前无系统信息配置，正在导入默认配置")
		dbConn.Create(&SystemSetting{
			Title:   "高性能基础防火墙",
			Icon:    "firewall.png",
			RunTime: 0,
		})
	} else {
		logger.Println("已配置系统信息")
	}
}

func AddRuntime() {
	logger.Println("开启引擎运行计时...")
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			ServerRuntimeIncrease()
		}
	}
}
