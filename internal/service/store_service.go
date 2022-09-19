package service

import (
	"context"
	"errors"
	"fmt"
	"gongfu/internal/model"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"time"
)
import _ "gorm.io/driver/mysql"

type Store interface {
	UserIDExists(ctx context.Context, userID uint) (bool, error)
	OpenIDExists(ctx context.Context, openID string) (bool, error)
	GetUserByOpenID(ctx context.Context, openID string) (*model.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	GetUser(ctx context.Context, userID uint) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	CreateCheckIn(ctx context.Context, checkIn *model.CheckIn) error
	GetTodayCheckIn(ctx context.Context, userID uint) (*model.CheckIn, error)
	GetCheckIn(ctx context.Context, id uint) (*model.CheckIn, error)
	GetCheckInByKey(ctx context.Context, key string) (*model.CheckIn, error)
	GetCheckInTop(ctx context.Context, startAt int64, endAt int64) ([]*model.CheckIn, error)
	GetUsersMap(ctx context.Context, userIDs []uint) (map[uint]*model.User, error)
}

var _ Store = (*DBStore)(nil)

type DBStore struct {
	DB *gorm.DB
}

func (s DBStore) GetUsersMap(ctx context.Context, userIDs []uint) (map[uint]*model.User, error) {
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

func (s DBStore) GetCheckInTop(ctx context.Context, startAt int64, endAt int64) ([]*model.CheckIn, error) {
	var checkIns = []*model.CheckIn{}
	start := time.Unix(startAt, 0).Format("2006-01-02 15:01:05")
	end := time.Unix(endAt, 0).Format("2006-01-02 15:01:05")
	fmt.Println(start, end)
	err := s.DB.WithContext(ctx).Where("id in (?)",
		s.DB.Model(&model.CheckIn{}).Select("max(id)").
			Where("created_at between ? and ?", start, end).
			Group("user_id"),
	).Order("created_at").Find(&checkIns).Error
	if err != nil {
		return nil, err
	}
	if len(checkIns) == 0 {
		return []*model.CheckIn{}, nil
	}
	return checkIns, nil
}

func (s DBStore) GetCheckIn(ctx context.Context, id uint) (*model.CheckIn, error) {
	var checkIn = &model.CheckIn{}
	err := s.DB.WithContext(ctx).Where("id = ?", id).Find(&checkIn).Error
	if err != nil {
		return nil, err
	}
	if checkIn.ID == 0 {
		return nil, nil
	}
	return checkIn, nil
}

func (s DBStore) GetCheckInByKey(ctx context.Context, key string) (*model.CheckIn, error) {
	var checkIn = &model.CheckIn{}
	err := s.DB.WithContext(ctx).Where("`key` = ?", key).Find(&checkIn).Error
	if err != nil {
		return nil, err
	}
	if checkIn.ID == 0 {
		return nil, nil
	}
	return checkIn, nil
}

func (s DBStore) CreateCheckIn(ctx context.Context, checkIn *model.CheckIn) error {
	return s.DB.WithContext(ctx).Create(checkIn).Error
}

func (s DBStore) GetTodayCheckIn(ctx context.Context, userID uint) (*model.CheckIn, error) {
	var checkIn = &model.CheckIn{}
	now := time.Now()
	startOfDay := fmt.Sprintf("%s 00:00:00", now.Format("2006-01-02"))
	endOfDay := fmt.Sprintf("%s 23:59:59", now.Format("2006-01-02"))
	err := s.DB.WithContext(ctx).Where("user_id = ?", userID).Where("created_at between ? and ?", startOfDay, endOfDay).Last(&checkIn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return checkIn, nil
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

func (s DBStore) UpdateUser(ctx context.Context, user *model.User) error {
	return s.DB.WithContext(ctx).Save(user).Error
}

func (s DBStore) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	var user = &model.User{}
	err := s.DB.WithContext(ctx).Where("id = ?", userID).Find(user).Error
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

func NewDBStore(DB *gorm.DB) Store {
	return &DBStore{DB: DB}
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
