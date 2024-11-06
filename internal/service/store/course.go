package store

import (
	"context"
	"fmt"
	"gongfu/internal/model"
	"gongfu/pkg/util"
	"gorm.io/datatypes"
)

var _ Course = (*DBStore)(nil)

type Course interface {
	// CreateCourse 创建课程
	CreateCourse(ctx context.Context, input CreateCourseInput) error
	// GetCoursePage 获取课程分页
	GetCoursePage(ctx context.Context, input GetCoursePageInput) (*GetCoursePageOutput, error)
	// GetCoursesByUser 获取所有课程
	GetCoursesByUser(ctx context.Context, userId uint) ([]*model.Course, error)
	// GetCourse 查询课程
	GetCourse(ctx context.Context, courseId uint) (*model.Course, error)
	// UpdateCourse 更新课程信息
	UpdateCourse(ctx context.Context, course *model.Course) error
	// DeleteCourse 删除课程
	DeleteCourse(ctx context.Context, id uint) error
}

func (s DBStore) GetCoursesByUser(ctx context.Context, userId uint) ([]*model.Course, error) {
	var courses = []*model.Course{}
	var err = s.DB.WithContext(ctx).Model(&model.Course{}).
		Where("coach_id = ?", userId).
		Or("manager_id = ?", userId).
		Or(datatypes.JSONArrayQuery("assistant_coach_ids").Contains(userId)).
		Order("start_date desc").
		Order("id desc").
		Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (s DBStore) CreateCourse(ctx context.Context, input CreateCourseInput) error {
	course := model.Course{
		StartTime:         input.StartTime,
		StartDate:         input.StartDate,
		SchoolId:          input.SchoolId,
		CoachId:           input.CoachId,
		AssistantCoachIds: util.JSON(input.AssistantCoachIds),
		Images:            datatypes.JSON{},
		Summary:           "",
		ManagerId:         input.ManagerId,
	}
	return s.DB.Model(&model.Course{}).Create(&course).Error
}

type CreateCourseInput struct {
	StartDate         string `json:"start_date"`          // 上课日期
	StartTime         string `json:"start_time"`          // 上课时间
	CoachId           *uint  `json:"coach_id"`            // 教练
	SchoolId          uint   `json:"school_id"`           // 学校
	AssistantCoachIds []uint `json:"assistant_coach_ids"` // 助理教练列表
	ManagerId         uint   `json:"manager_id"`          // 负责人
}

type GetCoursePageInput struct {
	Page    int
	Limit   int
	Keyword *string
	Desc    bool
	UserId  *uint
}

type GetCoursesInput struct {
	CoachId           *uint
	ManagerId         *uint
	AssistantCoachIds []uint
}

type GetCoursePageOutput struct {
	Items []*model.Course
	Count int64 // 总数
	Page  int
	Limit int
}

func (s DBStore) GetCoursePage(ctx context.Context, query GetCoursePageInput) (*GetCoursePageOutput, error) {
	q := s.DB.WithContext(ctx).Model(&model.Course{})
	if query.Desc {
		q = q.Order("id desc")
	}
	if query.UserId != nil {
		q = q.
			Where("manager_id = ?", *query.UserId).
			Where("coach_id = ?", *query.UserId).
			Or(datatypes.JSONArrayQuery("assistant_coach_ids").Contains(*query.UserId))
	}
	count := int64(0)
	if err := q.Count(&count).Error; err != nil {
		return nil, err
	}
	q = q.Limit(query.Limit).Offset(query.Page * query.Limit)
	courses := []*model.Course{}
	if err := q.Find(&courses).Error; err != nil {
		return nil, err
	}
	return &GetCoursePageOutput{
		Items: courses,
		Count: count,
		Page:  query.Page,
		Limit: query.Limit,
	}, nil
}

func (s DBStore) GetCourse(ctx context.Context, courseId uint) (*model.Course, error) {
	var course = model.Course{}
	err := s.DB.WithContext(ctx).Model(&model.Course{}).Where("id = ?", courseId).Find(&course).Error
	if err != nil {
		return nil, fmt.Errorf("query course: %w", err)
	}
	if course.ID == 0 {
		return nil, nil
	}
	return &course, nil
}

func (s DBStore) UpdateCourse(ctx context.Context, course *model.Course) error {
	if course == nil || course.ID == 0 {
		return fmt.Errorf("invalid course model")
	}
	if err := s.DB.WithContext(ctx).Save(course).Error; err != nil {
		return fmt.Errorf("update course: %w", err)
	}
	return nil
}

func (s DBStore) DeleteCourse(ctx context.Context, id uint) error {
	return s.DB.WithContext(ctx).Delete(&model.Course{}, id).Error
}
