package models

import (
	"github.com/shrewx/ginx/pkg/dbhelper"
	"gorm.io/gorm"
	"time"
)

func init() {
	dbhelper.RegisterTable(&User{})
}

// User 用户信息实体
type User struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`           // 用户唯一标识符
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"` // 用户名
	Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null"`   // 邮箱地址
	Name      string         `json:"name" gorm:"size:100;not null"`                // 真实姓名
	Status    string         `json:"status" gorm:"size:20;default:active"`         // 用户状态
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`             // 创建时间
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`             // 更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`                               // 软删除时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
