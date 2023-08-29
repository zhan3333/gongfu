package admin

import (
	"context"
	"fmt"
	"gongfu/internal/http/common"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
)

type UseCase struct {
	Store   store.Store
	Storage service.Storage
}

func NewUseCase(store store.Store, storage service.Storage) UseCase {
	return UseCase{
		Store:   store,
		Storage: storage,
	}
}

func (r UseCase) NewUsersQuery(userIds ...uint) (*common.UsersQuery, error) {
	usersMap, err := r.Store.GetUsersMap(context.TODO(), userIds)
	if err != nil {
		return nil, fmt.Errorf("get users map: %w", err)
	}
	return &common.UsersQuery{UsersMap: usersMap, Storage: r.Storage}, nil
}
