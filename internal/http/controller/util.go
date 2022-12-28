package controller

import (
	"context"
	"gongfu/internal/model"
	"strings"
)

func (r Controller) getUserHeadImgUrl(user *model.User) string {
	avatarURL := user.HeadImgURL
	if !strings.HasPrefix(avatarURL, "http") {
		avatarURL = r.Storage.GetPublicVisitURL(context.Background(), avatarURL) + "!avatar"
	}
	return avatarURL
}
