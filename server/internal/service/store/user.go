package store

import (
	"context"
	"fmt"
	"gongfu/internal/model"
)

type CreateTeachingRecordInput struct {
	Date    string
	Address string
	UserId  uint
}

type UpdateTeachingRecordInput struct {
	Id      uint
	Date    string
	Address string
}

type CreateStudyRecordInput struct {
	Date    string
	Content string
	UserId  uint
}

type UpdateStudyRecordInput struct {
	Id      uint
	Date    string
	Content string
}

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
	// GetCoachUserIds 获取教练 userIds
	GetCoachUserIds(ctx context.Context) ([]uint, error)
	GetCoaches(ctx context.Context) ([]model.User, error)
	// GetTeachingRecords 获取授课记录
	GetTeachingRecords(ctx context.Context, userId uint) ([]model.TeachingRecord, error)
	// GetStudyRecords 获取学习记录
	GetStudyRecords(ctx context.Context, userId uint) ([]model.StudyRecord, error)
	CreateTeachingRecord(ctx context.Context, in *CreateTeachingRecordInput) error
	UpdateTeachingRecord(ctx context.Context, in *UpdateTeachingRecordInput) error
	DeleteTeachingRecord(ctx context.Context, recordId uint) error

	CreateStudyRecord(ctx context.Context, in *CreateStudyRecordInput) error
	UpdateStudyRecord(ctx context.Context, in *UpdateStudyRecordInput) error
	DeleteStudyRecord(ctx context.Context, recordId uint) error
}

func (s DBStore) CreateTeachingRecord(ctx context.Context, in *CreateTeachingRecordInput) error {
	record := model.TeachingRecord{
		Date:    in.Date,
		Address: in.Address,
		UserId:  in.UserId,
	}
	return s.DB.WithContext(ctx).Create(&record).Error
}

func (s DBStore) UpdateTeachingRecord(ctx context.Context, in *UpdateTeachingRecordInput) error {
	return s.DB.WithContext(ctx).Save(&model.TeachingRecord{
		ID:      in.Id,
		Date:    in.Date,
		Address: in.Address,
	}).Error
}

func (s DBStore) DeleteTeachingRecord(ctx context.Context, recordId uint) error {
	return s.DB.WithContext(ctx).Delete(&model.TeachingRecord{}, recordId).Error
}

func (s DBStore) CreateStudyRecord(ctx context.Context, in *CreateStudyRecordInput) error {
	record := model.StudyRecord{
		Date:    in.Date,
		Content: in.Content,
		UserId:  in.UserId,
	}
	return s.DB.WithContext(ctx).Create(&record).Error
}

func (s DBStore) UpdateStudyRecord(ctx context.Context, in *UpdateStudyRecordInput) error {
	return s.DB.WithContext(ctx).Save(&model.StudyRecord{
		ID:      in.Id,
		Date:    in.Date,
		Content: in.Content,
	}).Error
}

func (s DBStore) DeleteStudyRecord(ctx context.Context, recordId uint) error {
	return s.DB.WithContext(ctx).Delete(&model.StudyRecord{}, recordId).Error
}

type UserPageQuery struct {
	Page    int
	Limit   int
	Keyword string
	Desc    bool
	RoleIds []int
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
	if len(query.RoleIds) != 0 {
		q = q.Joins("left join user_has_roles on users.id = user_has_roles.user_id").Where("user_has_roles.role_id in (?)", query.RoleIds)
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

func (s DBStore) GetCoachUserIds(ctx context.Context) ([]uint, error) {
	var userIds = []uint{}
	err := s.DB.WithContext(ctx).Model(&model.UserHasRole{}).Where("role_name = ?", "coach").Pluck("user_id", &userIds).Error
	return userIds, err
}

func (s DBStore) GetCoaches(ctx context.Context) ([]model.User, error) {
	coaches := []model.User{}
	userIds, err := s.GetCoachUserIds(ctx)
	if err != nil {
		return nil, fmt.Errorf("get coach user ids: %w", err)
	}
	err = s.DB.WithContext(ctx).Where("id in (?)", userIds).Find(&coaches).Error
	if err != nil {
		return nil, err
	}
	return coaches, nil
}

func (s DBStore) GetTeachingRecords(ctx context.Context, userId uint) ([]model.TeachingRecord, error) {
	records := []model.TeachingRecord{}
	err := s.DB.WithContext(ctx).Where("user_id = ?", userId).Order("id desc").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s DBStore) GetStudyRecords(ctx context.Context, userId uint) ([]model.StudyRecord, error) {
	records := []model.StudyRecord{}
	err := s.DB.WithContext(ctx).Where("user_id = ?", userId).Order("id desc").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}
