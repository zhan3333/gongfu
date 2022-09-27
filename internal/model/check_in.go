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

type CheckInDay struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// 用于查询指定用户的打卡历史
	UserID uint `json:"userID" gorm:"column:user_id;index:idx_user_id"`
	// 用于查询指定日期下的打卡历史: 20220921
	Date      string `json:"date" gorm:"column:date;index:idx_date_check_in_id"`
	CheckInID uint   `json:"checkInID" gorm:"column:check_in_id;index:idx_date_check_in_id"`
}

type CheckInCount struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint      `json:"userID" gorm:"column:user_id;index:idx_user_id;unique"`
	// 总打卡次数
	CheckInCount uint `json:"checkInCount" gorm:"column:check_in_count;index:idx_count"`
	// 连续打卡次数
	CheckInContinuous uint `json:"checkInContinuous" gorm:"column:check_in_continuous;index:idx_continuous"`
}
