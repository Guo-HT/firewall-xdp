package db

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       int    `gorm:"primaryKey"` // ID
	UserName string `gorm:"unique"`     // 用户名
	Password string // 密码
	Email    string // 邮箱
	//FirstLogin          bool           // 是否首次登录
	Role                int            // 用户角色[0-admin; 1-操作员; 2-访客]
	CreateAt            time.Time      // 创建时间
	DeleteAt            gorm.DeletedAt `gorm:"index"` // 是否删除
	LastErrorTime       time.Time      // 上一次错误时间
	IsLocked            bool           // 是否被锁定
	ErrorTimes          int            // 本轮错误次数
	IsChangedDefualtPwd bool           // 是否已更改初始密码
}

type SystemLog struct {
	ID        int       `gorm:"primaryKey"` // ID
	IP        string    // 操作IP
	User      User      // 操作用户
	UserID    int       // 操作用户ID(外键表示)
	Option    string    // 操作详情
	OptResult bool      // 操作结果
	CreateAt  time.Time // 操作事件
}
