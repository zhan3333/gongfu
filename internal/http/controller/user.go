package controller

import (
	"context"
	"gongfu/internal/app"
	"gongfu/internal/result"
	"gongfu/internal/service/store"
)

type GetUsersResponseUser struct {
	ID         uint     `json:"id"`
	OpenID     *string  `json:"openid"`
	Phone      *string  `json:"phone"`
	Nickname   string   `json:"nickname"`
	HeadImgURL string   `json:"headimgurl"`
	RoleNames  []string `json:"roleNames"`
	UUID       string   `json:"uuid"`
}

type GetUsersResponse struct {
	Users []GetUsersResponseUser `json:"users"`
	Page  int                    `json:"page"`
	Count int64                  `json:"count"`
	Limit int                    `json:"limit"`
}

// GetUsers 获取 users 分页
func (r Controller) GetUsers(c *app.Context) result.Result {
	req := struct {
		Page    int    `json:"page" form:"page"`
		Limit   int    `json:"limit" form:"limit"`
		Keyword string `json:"keyword" form:"keyword"`
		Desc    bool   `json:"desc" form:"desc"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	userPage, err := r.Store.GetUserPage(context.Background(), store.UserPageQuery{
		Page:    req.Page,
		Limit:   req.Limit,
		Keyword: req.Keyword,
		Desc:    req.Desc,
	})
	if err != nil {
		return result.Err(err)
	}
	ret := GetUsersResponse{
		Users: nil,
		Page:  userPage.Page,
		Count: userPage.Count,
		Limit: userPage.Limit,
	}
	for _, user := range userPage.Users {
		ret.Users = append(ret.Users, GetUsersResponseUser{
			ID:         user.ID,
			OpenID:     user.OpenID,
			Phone:      user.Phone,
			Nickname:   user.Nickname,
			HeadImgURL: r.getUserHeadImgUrl(user),
			RoleNames:  user.GetRoleNames(),
			UUID:       user.UUID,
		})
	}
	return result.Ok(ret)
}
