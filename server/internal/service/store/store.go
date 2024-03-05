package store

import (
	"context"
	"gongfu/internal/model"
	"gorm.io/gorm"
	"time"
)
import _ "gorm.io/driver/mysql"

type UpdateMemberCourseInput struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Total     int
	Remain    int
	Remark    string
	Status    string
}

type Store interface {
	// CheckIn 打卡相关查询
	CheckIn
	// CheckInCount 打卡统计相关查询
	CheckInCount
	// User 用户相关查询
	User
	// Role 角色相关查询
	Role
	// Coach 教练信息相关查询
	Coach
	// Course 课程
	Course
	// School 学校
	School
	GetMemberCourses(ctx context.Context, userId uint32) ([]*model.MemberCourse, error)
	CreateMemberCourse(ctx context.Context, in *CreateMemberCourseInput) error
	DeleteMemberCourse(ctx context.Context, userId uint) error
	UpdateMemberCourse(ctx context.Context, id uint, in *UpdateMemberCourseInput) error
}

var _ Store = (*DBStore)(nil)

type DBStore struct {
	DB *gorm.DB
}

func (s DBStore) UpdateMemberCourse(ctx context.Context, id uint, in *UpdateMemberCourseInput) error {
	return s.DB.WithContext(ctx).Model(&model.MemberCourse{}).Where("id = ?", id).Updates(map[string]any{
		"name":       in.Name,
		"start_time": in.StartTime,
		"end_time":   in.EndTime,
		"total":      in.Total,
		"remain":     in.Remain,
		"remark":     in.Remark,
		"status":     in.Status,
	}).Error
}

func (s DBStore) DeleteMemberCourse(ctx context.Context, id uint) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.MemberCourse{}).Error
}

func (s DBStore) CreateMemberCourse(ctx context.Context, in *CreateMemberCourseInput) error {
	course := &model.MemberCourse{
		UserId:    uint(in.UserId),
		Name:      in.Name,
		StartTime: in.StartTime,
		EndTime:   in.EndTime,
		Total:     in.Total,
		Remark:    in.Remark,
		Status:    in.Status,
		Remain:    in.Remain,
	}
	return s.DB.WithContext(ctx).Create(course).Error
}

type CreateMemberCourseInput struct {
	UserId    uint32
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Total     int
	Remark    string
	Status    string
	Remain    int
}

func (s DBStore) GetMemberCourses(ctx context.Context, userId uint32) ([]*model.MemberCourse, error) {
	var courses []*model.MemberCourse
	err := s.DB.WithContext(ctx).Where("user_id = ?", userId).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func NewDBStore(DB *gorm.DB) Store {
	return &DBStore{DB: DB}
}
