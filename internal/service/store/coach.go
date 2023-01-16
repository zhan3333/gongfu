package store

import (
	"context"
	"gongfu/internal/model"
)

type Coach interface {
	GetCoach(ctx context.Context, userID uint) (*model.Coach, error)
	InsertOrUpdateCoach(ctx context.Context, coach *model.Coach) error
}

func (s DBStore) GetCoach(ctx context.Context, userID uint) (*model.Coach, error) {
	coach := model.Coach{}
	err := s.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&coach).Error
	if err != nil {
		return nil, err
	}
	if coach.ID == 0 {
		return nil, nil
	}
	return &coach, nil
}

func (s DBStore) InsertOrUpdateCoach(ctx context.Context, coach *model.Coach) error {
	return s.DB.WithContext(ctx).Save(coach).Error
}
