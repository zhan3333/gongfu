package controller

import (
	"github.com/silenceper/wechat/v2/officialaccount"
	"gongfu/internal/config"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
)

type Controller struct {
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
) *Controller {
	return &Controller{Config: config, OfficialAccount: officialAccount, AuthCode: authCode, Token: token, Store: store, Storage: storage}
}
