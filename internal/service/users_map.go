package service

import (
	"gongfu/internal/model"
)

type UsersMap map[uint]*model.User

const defaultUserName = "no user name"
const defaultHeadImgUrl = ""

func (m UsersMap) GetInfo(userID uint) (userName string, headImgUrl string) {
	if u, ok := m[userID]; ok {
		return u.Nickname, u.HeadImgURL
	} else {
		return defaultUserName, defaultHeadImgUrl
	}
}

func (m UsersMap) DefaultGet(userID uint) *model.User {
	if u, ok := m[userID]; ok {
		return u
	} else {
		return &model.User{
			Nickname:   defaultUserName,
			HeadImgURL: defaultHeadImgUrl,
		}
	}
}
