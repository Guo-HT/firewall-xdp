package db

import (
	"gorm.io/gorm"
	"time"
)

/****************************************************/
//	添加表后，记得修改init中的迁移、初始化
/****************************************************/

type User struct {
	ID                  int            `gorm:"primaryKey;autoIncrement"` // ID
	UserName            string         `gorm:"unique"`                   // 用户名
	Password            string         // 密码
	Email               string         // 邮箱
	Role                int            // 用户角色[0-admin; 1-操作员; 2-访客]
	CreateAt            time.Time      // 创建时间
	DeleteAt            gorm.DeletedAt `gorm:"index"` // 是否删除
	LastErrorTime       time.Time      // 上一次错误时间
	IsLocked            bool           `gorm:"default:false;"` // 是否被锁定
	ErrorTimes          int            `gorm:"default:0;"`     // 本轮错误次数
	IsChangedDefualtPwd bool           `gorm:"default:false;"` // 是否已更改初始密码
	//FirstLogin          bool           // 是否首次登录
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

type SystemSetting struct {
	ID      int    `gorm:"primaryKey"`
	Title   string // 系统名称
	Icon    string // 系统图标
	RunTime int64  // 运行时间
}
