package controller

import (
	"context"
	"encoding/json"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gongfu/internal/service/store"
	"net/http"
	"strconv"
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

// AdminGetUsers 获取 users 分页
func (r Controller) AdminGetUsers(c *app.Context) result.Result {
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

func (r Controller) AdminGetUser(c *app.Context) result.Result {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return result.Err(err)
	}
	user, err := r.Store.GetUser(context.TODO(), uint(userID))
	if err != nil {
		return result.Err(err)
	}
	if user == nil {
		c.String(http.StatusNotFound, "user not found")
		return result.Err(nil)
	}
	return result.Ok(MeResponse{
		ID:         user.ID,
		OpenID:     user.OpenID,
		Phone:      user.Phone,
		Nickname:   user.Nickname,
		HeadImgURL: r.getUserHeadImgUrl(user),
		RoleNames:  user.GetRoleNames(),
		UUID:       user.UUID,
	})
}

func (r Controller) AdminUpdateUser(c *app.Context) result.Result {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return result.Err(err)
	}
	req := struct {
		Nickname string `json:"nickname" binding:"omitempty,max=20"`
		//
		Phone string `json:"phone" binding:"omitempty,max=20"`
		// 等级
		Level string `json:"level"`
		// 任教单位
		TeachingSpace string `json:"teachingSpace"`
		// 任教年限
		TeachingAge string `json:"teachingAge"`
		// 任教经历
		TeachingExperiences []string `json:"teachingExperiences"`
		// 设置角色
		RoleNames []string `json:"roleNames"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	user, err := r.Store.GetUser(context.TODO(), uint(userID))
	if err != nil {
		return result.Err(err)
	}
	if user == nil {
		c.String(http.StatusNotFound, "user not found")
		return result.Err(nil)
	}
	user.Nickname = req.Nickname
	user.Phone = &req.Phone
	if err := r.Store.UpdateUser(context.TODO(), user); err != nil {
		return result.Err(err)
	}

	// 更新 coach 资料
	coach, err := r.Store.GetCoach(context.TODO(), uint(userID))
	if err != nil {
		return result.Err(err)
	}
	if coach == nil {
		coach = &model.Coach{
			UserID:              uint(userID),
			Level:               req.Level,
			TeachingSpace:       "",
			TeachingAge:         "",
			TeachingExperiences: nil,
		}
	}
	coach.Level = req.Level
	coach.TeachingSpace = req.TeachingSpace
	coach.TeachingAge = req.TeachingAge
	exp, _ := json.Marshal(req.TeachingExperiences)
	coach.TeachingExperiences = exp
	if err := r.Store.InsertOrUpdateCoach(context.TODO(), coach); err != nil {
		return result.Err(err)
	}
	// 更新角色
	if err := r.Store.SyncUserRoles(context.TODO(), uint(userID), req.RoleNames); err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

// AdminGetRoleNames 查询所有的角色名称
func (r Controller) AdminGetRoleNames(c *app.Context) result.Result {
	var err error
	var roleNames = []string{}
	if roleNames, err = r.Store.GetRoleNames(context.TODO()); err != nil {
		return result.Err(err)
	}
	return result.Ok(roleNames)
}
