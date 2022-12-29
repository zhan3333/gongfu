package store

import "gorm.io/gorm"
import _ "gorm.io/driver/mysql"

type Store interface {
	// CheckIn 打卡相关查询
	CheckIn
	// CheckInCount 打卡统计相关查询
	CheckInCount
	// User 用户相关查询
	User
	// Role 角色相关查询
	Role
	// Coach 教练信息相关查询
	Coach
}

var _ Store = (*DBStore)(nil)

type DBStore struct {
	DB *gorm.DB
}

func NewDBStore(DB *gorm.DB) Store {
	return &DBStore{DB: DB}
}
