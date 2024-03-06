package store

import (
	"context"
	"errors"
	"gongfu/internal/model"
	"gorm.io/gorm"
	"time"
)
import _ "gorm.io/driver/mysql"

type UpdateMemberCourseInput struct {
	Name         string
	StartTime    time.Time
	EndTime      time.Time
	Total        int
	Remain       int
	Remark       string
	Status       string
	UpdateUserId uint32
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
	ChangeMemberCourseRemain(ctx context.Context, id uint, remain int) error
	GetMemberCourse(ctx context.Context, id uint) (*model.MemberCourse, error)
	CreateCheckInComment(ctx context.Context, checkInId uint32, userId uint32, content string) error
	GetCheckInComments(ctx context.Context, checkInId uint32) ([]*model.CheckInComment, error)
}

var _ Store = (*DBStore)(nil)

type DBStore struct {
	DB *gorm.DB
}

func (s DBStore) GetCheckInComments(ctx context.Context, checkInId uint32) ([]*model.CheckInComment, error) {
	var comments []*model.CheckInComment
	err := s.DB.WithContext(ctx).Where("check_in_id = ?", checkInId).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil

}

func (s DBStore) CreateCheckInComment(ctx context.Context, checkInId uint32, userId uint32, content string) error {
	// create comment
	comment := model.CheckInComment{
		Content:      content,
		CheckInId:    checkInId,
		CreateUserId: userId,
	}
	return s.DB.WithContext(ctx).Create(&comment).Error
}

func (s DBStore) GetMemberCourse(ctx context.Context, id uint) (*model.MemberCourse, error) {
	var course model.MemberCourse
	err := s.DB.WithContext(ctx).Where("id = ?", id).First(&course).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (s DBStore) ChangeMemberCourseRemain(ctx context.Context, id uint, remain int) error {
	return s.DB.WithContext(ctx).Model(&model.MemberCourse{}).Where("id = ?", id).Update("remain", remain).Error
}

func (s DBStore) UpdateMemberCourse(ctx context.Context, id uint, in *UpdateMemberCourseInput) error {
	return s.DB.WithContext(ctx).Model(&model.MemberCourse{}).Where("id = ?", id).Updates(map[string]any{
		"name":         in.Name,
		"start_time":   in.StartTime,
		"end_time":     in.EndTime,
		"total":        in.Total,
		"remain":       in.Remain,
		"remark":       in.Remark,
		"status":       in.Status,
		"updateUserId": in.UpdateUserId,
	}).Error
}

func (s DBStore) DeleteMemberCourse(ctx context.Context, id uint) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.MemberCourse{}).Error
}

func (s DBStore) CreateMemberCourse(ctx context.Context, in *CreateMemberCourseInput) error {
	course := &model.MemberCourse{
		UserId:       uint(in.UserId),
		Name:         in.Name,
		StartTime:    in.StartTime,
		EndTime:      in.EndTime,
		Total:        in.Total,
		Remark:       in.Remark,
		Status:       in.Status,
		Remain:       in.Remain,
		CreateUserId: in.CreateUserId,
	}
	return s.DB.WithContext(ctx).Create(course).Error
}

type CreateMemberCourseInput struct {
	UserId       uint32
	Name         string
	StartTime    time.Time
	EndTime      time.Time
	Total        int
	Remark       string
	Status       string
	Remain       int
	CreateUserId uint32
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
