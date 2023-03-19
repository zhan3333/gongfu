package admin

import (
	"context"
	"fmt"
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

type UsersQuery struct {
	usersMap store.UsersMap
	storage  service.Storage
}

type Coach struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (q UsersQuery) GetCoach(id *uint) *Coach {
	if id == nil {
		return nil
	}
	if user := q.usersMap[*id]; user != nil {
		return &Coach{
			ID:     user.ID,
			Name:   user.Nickname,
			Avatar: q.storage.GetHeadImageVisitURL(user.HeadImgURL),
		}
	}
	return nil
}

func (q UsersQuery) GetCoaches(ids ...uint) []Coach {
	var items = []Coach{}
	for _, id := range ids {
		if coach := q.GetCoach(&id); coach != nil {
			items = append(items, *coach)
		}
	}
	return items
}

func (r UseCase) NewUsersQuery(userIds ...uint) (*UsersQuery, error) {
	usersMap, err := r.Store.GetUsersMap(context.TODO(), userIds)
	if err != nil {
		return nil, fmt.Errorf("get users map: %w", err)
	}
	return &UsersQuery{usersMap: usersMap, storage: r.Storage}, nil
}
