package admin

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
func (r UseCase) AdminGetUsers(c *app.Context) result.Result {
	req := struct {
		Page    int    `json:"page" form:"page"`
		Limit   int    `json:"limit" form:"limit"`
		Keyword string `json:"keyword" form:"keyword"`
		Desc    bool   `json:"desc" form:"desc"`
		RoleIds []int  `json:"roleIds" form:"roleIds"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	userPage, err := r.Store.GetUserPage(context.Background(), store.UserPageQuery{
		Page:    req.Page,
		Limit:   req.Limit,
		Keyword: req.Keyword,
		Desc:    req.Desc,
		RoleIds: req.RoleIds,
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
			HeadImgURL: r.Storage.GetHeadImageVisitURL(user.HeadImgURL),
			RoleNames:  user.GetRoleNames(),
			UUID:       user.UUID,
		})
	}
	return result.Ok(ret)
}

type MeResponse struct {
	ID              uint                   `json:"id"`
	OpenID          *string                `json:"openid"`
	Phone           *string                `json:"phone"`
	Nickname        string                 `json:"nickname"`
	HeadImgURL      string                 `json:"headimgurl"`
	RoleNames       []string               `json:"roleNames"`
	UUID            string                 `json:"uuid"`
	TeachingRecords []model.TeachingRecord `json:"teachingRecords"`
	StudyRecords    []model.StudyRecord    `json:"studyRecords"`
}

func (r UseCase) AdminGetUser(c *app.Context) result.Result {
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
	teachingRecords, err := r.Store.GetTeachingRecords(c.Request.Context(), user.ID)
	if err != nil {
		return result.Err(err)
	}
	studyRecords, err := r.Store.GetStudyRecords(c.Request.Context(), user.ID)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(MeResponse{
		ID:              user.ID,
		OpenID:          user.OpenID,
		Phone:           user.Phone,
		Nickname:        user.Nickname,
		HeadImgURL:      r.Storage.GetHeadImageVisitURL(user.HeadImgURL),
		RoleNames:       user.GetRoleNames(),
		UUID:            user.UUID,
		TeachingRecords: teachingRecords,
		StudyRecords:    studyRecords,
	})
}

func (r UseCase) AdminUpdateUser(c *app.Context) result.Result {
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
func (r UseCase) AdminGetRoleNames(c *app.Context) result.Result {
	var err error
	var roleNames = []string{}
	if roleNames, err = r.Store.GetRoleNames(context.TODO()); err != nil {
		return result.Err(err)
	}
	return result.Ok(roleNames)
}

// AdminEditTeachingRecord 编辑授课记录
func (r UseCase) AdminEditTeachingRecord(c *app.Context) result.Result {
	req := struct {
		Id      uint   `json:"id" binding:"omitempty"`
		Date    string `json:"date" binding:"required"`
		Address string `json:"address" binding:"required"`
		UserId  uint   `json:"userId" binding:"required"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	if req.Id == 0 {
		err := r.Store.CreateTeachingRecord(c.Request.Context(), &store.CreateTeachingRecordInput{
			Date:    req.Date,
			Address: req.Address,
			UserId:  req.UserId,
		})
		if err != nil {
			return result.Err(err)
		}
	} else {
		err := r.Store.UpdateTeachingRecord(c.Request.Context(), &store.UpdateTeachingRecordInput{
			Id:      req.Id,
			Date:    req.Date,
			Address: req.Address,
		})
		if err != nil {
			return result.Err(err)
		}
	}
	return result.Ok(nil)
}

func (r UseCase) AdminDeleteTeachingRecord(c *app.Context) result.Result {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return result.Err(err)
	}
	err = r.Store.DeleteTeachingRecord(c.Request.Context(), uint(id))
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

func (r UseCase) AdminEditStudyRecord(c *app.Context) result.Result {
	req := struct {
		Id      uint   `json:"id" binding:"omitempty"`
		Date    string `json:"date" binding:"required"`
		Content string `json:"content" binding:"required"`
		UserId  uint   `json:"userId" binding:"required"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	if req.Id == 0 {
		err := r.Store.CreateStudyRecord(c.Request.Context(), &store.CreateStudyRecordInput{
			Date:    req.Date,
			Content: req.Content,
			UserId:  req.UserId,
		})
		if err != nil {
			return result.Err(err)
		}
	} else {
		err := r.Store.UpdateStudyRecord(c.Request.Context(), &store.UpdateStudyRecordInput{
			Id:      req.Id,
			Date:    req.Date,
			Content: req.Content,
		})
		if err != nil {
			return result.Err(err)
		}
	}
	return result.Ok(nil)
}

func (r UseCase) AdminDeleteStudyRecord(c *app.Context) result.Result {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return result.Err(err)
	}
	err = r.Store.DeleteStudyRecord(c.Request.Context(), uint(id))
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}
