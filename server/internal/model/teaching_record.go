package model

import (
	"gorm.io/gorm"
	"time"
)

// TeachingRecord 授课记录
type TeachingRecord struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Date      string         `json:"date" gorm:"index"` // 日期
	Address   string         `json:"address"`           // 授课地点
	UserId    uint           `json:"userId"`            // 记录所属用户
}
