package store

import (
	"context"
	"gongfu/internal/model"
)

var _ School = (*DBStore)(nil)

type School interface {
	GetSchools(ctx context.Context) ([]*model.School, error)
	GetSchool(ctx context.Context, id uint) (*model.School, error)
}

// GetSchools 获取学校列表
func (s DBStore) GetSchools(ctx context.Context) ([]*model.School, error) {
	var schools = []*model.School{}
	if err := s.DB.WithContext(ctx).Find(&schools).Error; err != nil {
		return nil, err
	}
	return schools, nil
}

func (s DBStore) GetSchool(ctx context.Context, id uint) (*model.School, error) {
	var school = model.School{}
	if err := s.DB.WithContext(ctx).Where("id = ?", id).Find(&school).Error; err != nil {
		return nil, err
	}
	if school.ID == 0 {
		return nil, nil
	}
	return &school, nil
}
