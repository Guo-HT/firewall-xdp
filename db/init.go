package db

import (
	"gorm.io/gorm"
	"log"
	"os"
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
}
