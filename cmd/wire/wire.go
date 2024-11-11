//go:generate go run github.com/google/wire/cmd/wire

//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"gongfu/internal/client"
	"gongfu/internal/config"
	"gongfu/internal/http"
	"gongfu/internal/http/admin"
	"gongfu/internal/http/controller"
	"gongfu/internal/http/middlewares"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
)

func NewApp(
	config *config.Config,
) (*http.Route, error) {
	panic(wire.Build(
		http.NewRoute,
		controller.NewController,
		admin.NewUseCase,
		middlewares.NewMiddlewares,
		client.NewRedis,
		service.NewOfficialAccount,
		service.NewAuthCodeService,
		service.NewTokenService,
		service.NewStorageService,
		store.NewStore,
		service.NewWechat,
	))
}
