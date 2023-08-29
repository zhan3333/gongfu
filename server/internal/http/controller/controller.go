package controller

import (
	"context"
	"fmt"
	"github.com/silenceper/wechat/v2/officialaccount"
	"gongfu/internal/config"
	"gongfu/internal/http/common"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
)

type UseCase struct {
	Config          *config.Config
	OfficialAccount *officialaccount.OfficialAccount
	AuthCode        service.AuthCode
	Token           service.Token
	Store           store.Store
	Storage         service.Storage
}

func NewController(
	config *config.Config,
	officialAccount *officialaccount.OfficialAccount,
	authCode service.AuthCode,
	token service.Token,
	store store.Store,
	storage service.Storage,
) *UseCase {
	return &UseCase{Config: config, OfficialAccount: officialAccount, AuthCode: authCode, Token: token, Store: store, Storage: storage}
}

func (r UseCase) NewUsersQuery(userIds ...uint) (*common.UsersQuery, error) {
	usersMap, err := r.Store.GetUsersMap(context.TODO(), userIds)
	if err != nil {
		return nil, fmt.Errorf("get users map: %w", err)
	}
	return &common.UsersQuery{UsersMap: usersMap, Storage: r.Storage}, nil
}
