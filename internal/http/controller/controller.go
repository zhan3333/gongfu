package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/client"
	"gongfu/internal/config"
	"gongfu/internal/http/common"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
	"strconv"
)

type UseCase struct {
	Config          *config.Config
	OfficialAccount *client.OfficialAccount
	AuthCode        service.AuthCode
	Token           service.Token
	Store           store.Store
	Storage         service.Storage
	Wechat          service.Wechat
}

func NewController(
	config *config.Config,
	officialAccount *client.OfficialAccount,
	authCode service.AuthCode,
	token service.Token,
	store store.Store,
	storage service.Storage,
	wechat service.Wechat,
) *UseCase {
	return &UseCase{Config: config, OfficialAccount: officialAccount, AuthCode: authCode, Token: token, Store: store, Storage: storage, Wechat: wechat}
}

func (r UseCase) NewUsersQuery(userIds ...uint) (*common.UsersQuery, error) {
	usersMap, err := r.Store.GetUsersMap(context.TODO(), userIds)
	if err != nil {
		return nil, fmt.Errorf("get users map: %w", err)
	}
	return &common.UsersQuery{UsersMap: usersMap, Storage: r.Storage}, nil
}

type CreateCheckInCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (r UseCase) CreateCheckInComment(c *app.Context) result.Result {
	if !c.User.HasAnyRole([]string{model.ROLE_ADMIN, model.ROLE_COACH}) {
		return result.Err(fmt.Errorf("permission denied"))
	}
	checkInId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return result.Err(err)
	}
	if checkInId == 0 {
		return result.Err(fmt.Errorf("invalid check in id"))
	}
	var req CreateCheckInCommentRequest
	if err := c.Bind(&req); err != nil {
		return result.Err(err)
	}
	// create check in comment
	err = r.Store.CreateCheckInComment(c.Request.Context(), uint32(checkInId), uint32(c.User.ID), req.Content)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

func (r UseCase) GetCheckInComments(c *app.Context) result.Result {
	checkInId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return result.Err(err)
	}
	if checkInId == 0 {
		return result.Err(fmt.Errorf("invalid check in id"))
	}
	comments, err := r.Store.GetCheckInComments(c.Request.Context(), uint32(checkInId))
	if err != nil {
		return result.Err(err)
	}
	var userIds []uint
	for _, v := range comments {
		userIds = append(userIds, uint(v.CreateUserId))
	}
	// get users map
	usersMap, err := r.NewUsersQuery(userIds...)
	if err != nil {
		return result.Err(err)
	}
	var res []gin.H
	for _, v := range comments {
		user := usersMap.UsersMap.DefaultGet(uint(v.CreateUserId))
		res = append(res, gin.H{
			"id":        v.ID,
			"content":   v.Content,
			"createdAt": v.CreatedAt,
			"user": gin.H{
				"id":         user.ID,
				"name":       user.Nickname,
				"headImgUrl": r.Storage.GetHeadImageVisitURL(user.HeadImgURL),
				"uuid":       user.UUID,
			},
		})
	}
	return result.Ok(res)
}
