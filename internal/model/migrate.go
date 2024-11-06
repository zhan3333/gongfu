package model

import (
	"time"
)

type Migrate struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	Action    string    `json:"action" gorm:"index"`
}
