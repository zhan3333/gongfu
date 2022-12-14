package store

import (
	"context"
	"gongfu/internal/model"
)

type User interface {
	UserIDExists(ctx context.Context, userID uint) (bool, error)
	OpenIDExists(ctx context.Context, openID string) (bool, error)
	GetUserByOpenID(ctx context.Context, openID string) (*model.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	GetUserByUUID(ctx context.Context, uuid string) (*model.User, error)
	GetUser(ctx context.Context, userID uint) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	GetUsersMap(ctx context.Context, userIDs []uint) (UsersMap, error)
	GetUserPage(ctx context.Context, query UserPageQuery) (*UserPage, error)
}

type UserPageQuery struct {
	Page    int
	Limit   int
	Keyword string
	Desc    bool
}

type UserPage struct {
	Users []*model.User
	Count int64 // 总数
	Page  int
	Limit int
}

func (s DBStore) GetUserPage(ctx context.Context, query UserPageQuery) (*UserPage, error) {
	q := s.DB.WithContext(ctx).Model(&model.User{}).Preload("Roles")
	if query.Desc {
		q = q.Order("id desc")
	}
	if query.Keyword != "" {
		q = q.Where("nickname like ?", query.Keyword).Or("phone like ?", query.Keyword)
	}
	count := int64(0)
	if err := q.Count(&count).Error; err != nil {
		return nil, err
	}
	q = q.Limit(query.Limit).Offset(query.Page * query.Limit)
	users := []*model.User{}
	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	return &UserPage{
		Users: users,
		Count: count,
		Page:  query.Page,
		Limit: query.Limit,
	}, nil
}

func (s DBStore) UserIDExists(ctx context.Context, userID uint) (bool, error) {
	count := int64(0)
	err := s.DB.WithContext(ctx).Where("id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s DBStore) OpenIDExists(ctx context.Context, openID string) (bool, error) {
	count := int64(0)
	err := s.DB.WithContext(ctx).Where("open_id = ?", openID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s DBStore) CreateUser(ctx context.Context, user *model.User) error {
	return s.DB.WithContext(ctx).Create(&user).Error
}

func (s DBStore) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user = &model.User{}
	err := s.DB.WithContext(ctx).Where("phone = ?", phone).Find(user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func (s DBStore) GetUserByUUID(ctx context.Context, uuid string) (*model.User, error) {
	var user = &model.User{}
	err := s.DB.WithContext(ctx).Preload("Roles").Where("uuid = ?", uuid).Find(user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func (s DBStore) UpdateUser(ctx context.Context, user *model.User) error {
	return s.DB.WithContext(ctx).Save(user).Error
}

func (s DBStore) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	var user = &model.User{}
	err := s.DB.WithContext(ctx).Preload("Roles").Where("id = ?", userID).Find(user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func (s DBStore) GetUserByOpenID(ctx context.Context, openID string) (*model.User, error) {
	var user = &model.User{}
	err := s.DB.WithContext(ctx).Where("open_id = ?", openID).Find(user).Error
	if err != nil {
		return nil, err
	}
	if user.ID == 0 {
		return nil, nil
	}
	return user, nil
}

func (s DBStore) GetUsersMap(ctx context.Context, userIDs []uint) (UsersMap, error) {
	var users = []*model.User{}
	err := s.DB.WithContext(ctx).Where("id in (?)", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	ret := map[uint]*model.User{}
	for _, user := range users {
		ret[user.ID] = user
	}
	return ret, nil
}
