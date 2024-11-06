package model

import (
	"gorm.io/gorm"
	"time"
)

// StudyRecord 学习记录
type StudyRecord struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Date      string         `json:"date" gorm:"index"` // 日期
	Content   string         `json:"content"`           // 学习内容
	UserId    uint           `json:"userId"`            // 记录所属用户
}
