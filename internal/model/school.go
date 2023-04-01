package model

import (
	"gorm.io/gorm"
	"time"
)

type School struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Address   string         `json:"address"` // 地址
	Name      string         `json:"name"`    // 学习名称
}
