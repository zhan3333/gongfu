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
	// 用户角色，可以为 admin/user/coach
	Role string `json:"role" gorm:"column:role;index;default:user"`
	// 用户唯一编码
	UUID string `json:"uuid" gorm:"unique"`
}

func (u User) IsAdmin() bool {
	return u.Role == ROLE_ADMIN
}

func (u User) IsUser() bool {
	return u.Role == ROLE_USER
}

func (u User) IsCoach() bool {
	return u.Role == ROLE_COACH
}

// ROLE_ADMIN 管理员
const ROLE_ADMIN = "admin"

// ROLE_USER 用户
const ROLE_USER = "user"

// ROLE_COACH 教练
const ROLE_COACH = "coach"
