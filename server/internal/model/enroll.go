package model

import (
	"gorm.io/gorm"
	"time"
)

type Enroll struct {
	ID            uint `gorm:"primarykey"`
	ActivityId    string
	Status        string
	Amount        int64
	UserId        uint
	OutTradeNo    string
	TransactionId string
	Description   string
	Attach        string

	Username  string
	Sex       string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
