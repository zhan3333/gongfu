package common

import (
	"gongfu/internal/service"
	"gongfu/internal/service/store"
)

type UsersQuery struct {
	UsersMap store.UsersMap
	Storage  service.Storage
}

type Coach struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (q UsersQuery) Name(id uint) string {
	if user := q.UsersMap[id]; user != nil {
		return user.Nickname
	}
	return "unknown"
}

func (q UsersQuery) Names(ids ...uint) []string {
	var names []string
	for _, id := range ids {
		names = append(names, q.Name(id))
	}
	return names
}

func (q UsersQuery) GetCoach(id uint) *Coach {
	if id == 0 {
		return nil
	}
	if user := q.UsersMap[id]; user != nil {
		return &Coach{
			ID:     user.ID,
			Name:   user.Nickname,
			Avatar: q.Storage.GetHeadImageVisitURL(user.HeadImgURL),
		}
	}
	return nil
}

func (q UsersQuery) GetCoaches(ids ...uint) []Coach {
	var items []Coach
	for _, id := range ids {
		if coach := q.GetCoach(id); coach != nil {
			items = append(items, *coach)
		}
	}
	return items
}
