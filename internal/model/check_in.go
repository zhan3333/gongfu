package model

import (
	"gorm.io/gorm"
	"time"
)

type CheckIn struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	// 保存路径
	Key      string `json:"key" gorm:"unique;column:key"`
	UserID   uint   `json:"userID" gorm:"column:user_id;index"`
	FileName string `json:"fileName" gorm:"column:file_name"`
}
