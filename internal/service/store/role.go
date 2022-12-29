package store

import (
	"context"
	"gongfu/internal/model"
)

type Role interface {
	GetUserRoles(ctx context.Context, userID uint) ([]string, error)
}

func (s DBStore) GetUserRoles(ctx context.Context, userID uint) ([]string, error) {
	roleNames := []string{}
	err := s.DB.WithContext(ctx).
		Model(&model.UserHasRole{}).
		Where("user_id = ?", userID).
		Pluck("role_name", &roleNames).Error
	if err != nil {
		return nil, err
	}
	return roleNames, nil
}
