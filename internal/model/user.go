package model

import (
	"gorm.io/gorm"
	"time"
)

type CoachStatus string

const (
	CoachStatusRegular    CoachStatus = "regular"
	CoachStatusInternship CoachStatus = "internship"
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
	// 用户唯一编码
	UUID        string `json:"uuid" gorm:"unique"`
	Roles       []UserHasRole
	CoachStatus CoachStatus `json:"coachStatus" gorm:"column:coach_status"`
}

func (u User) GetRoleNames() []string {
	ret := []string{}
	for _, v := range u.Roles {
		ret = append(ret, v.RoleName)
	}
	return ret
}

func (u User) HasRole(role string) bool {
	for _, v := range u.Roles {
		if v.RoleName == role {
			return true
		}
	}
	return false
}

func (u User) HasAnyRole(roles []string) bool {
	for _, v := range u.Roles {
		for _, v2 := range roles {
			if v.RoleName == v2 {
				return true
			}
		}
	}
	return false
}

// ROLE_ADMIN 管理员
const ROLE_ADMIN = "admin"

// ROLE_USER 用户
const ROLE_USER = "user"

// ROLE_COACH 教练
const ROLE_COACH = "coach"

const ROLE_MEMBER = "member"
