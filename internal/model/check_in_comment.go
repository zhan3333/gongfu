package model

import (
	"gorm.io/gorm"
	"time"
)

type CheckInComment struct {
	ID           uint   `json:"id" gorm:"primarykey"`
	Content      string `json:"content" gorm:"column:content;type:text;not null;comment:评论内容"`
	CheckInId    uint32 `json:"checkInId" gorm:"column:check_in_id;type:bigint(32);comment:签到 id"`
	CreateUserId uint32 `json:"createUserId" gorm:"column:create_user_id;type:bigint(32);not null;default:0;comment:创建用户 id"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func (c CheckInComment) TableName() string {
	return "check_in_comments"
}
