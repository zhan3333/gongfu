package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	OpenID     *string        `json:"openid" gorm:"column:open_id;unique"`
	UnionID    *string        `json:"unionid" gorm:"column:union_id"`
	Phone      *string        `json:"phone" gorm:"column:phone;unique"`
	Nickname   string         `json:"nickname" gorm:"column:nickname"`
	Sex        int32          `json:"sex" gorm:"column:sex"`
	Province   string         `json:"province" gorm:"column:province"`
	City       string         `json:"city" gorm:"column:city"`
	Country    string         `json:"country" gorm:"column:country"`
	HeadImgURL string         `json:"headimgurl" gorm:"column:head_img_url"`
}
