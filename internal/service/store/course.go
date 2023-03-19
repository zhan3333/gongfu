package store

import (
	"context"
	"database/sql"
	"fmt"
	"gongfu/internal/model"
	"gongfu/pkg/util"
	"gorm.io/datatypes"
	"strconv"
	"time"
)

var _ Course = (*DBStore)(nil)

type Course interface {
	// CreateCourse 创建课程
	CreateCourse(ctx context.Context, input CreateCourseInput) error
	// GetCoursePage 获取课程分页
	GetCoursePage(ctx context.Context, input GetCoursePageInput) (*GetCoursePageOutput, error)
	// GetCourse 查询课程
	GetCourse(ctx context.Context, courseId uint) (*model.Course, error)
	// UpdateCourse 更新课程信息
	UpdateCourse(ctx context.Context, course *model.Course) error
}

func (s DBStore) CreateCourse(ctx context.Context, input CreateCourseInput) error {
	course := model.Course{
		SchoolStartAt:     sql.NullTime{input.SchoolStartAt, true},
		Address:           input.Address,
		Name:              input.Name,
		CoachId:           &input.CoachId,
		AssistantCoachIds: util.JSON(input.AssistantCoachIds),
		Images:            datatypes.JSON{},
		Summary:           "",
	}
	return s.DB.Model(&model.Course{}).Create(&course).Error
}

type CreateCourseInput struct {
	SchoolStartAt     time.Time `json:"school_start_at"`     // 上课时间
	Address           string    `json:"address"`             // 上课地点
	Name              string    `json:"name"`                // 课程名称
	CoachId           uint      `json:"coach_id"`            // 教练
	AssistantCoachIds []uint    `json:"assistant_coach_ids"` // 助理教练列表
}

type GetCoursePageInput struct {
	Page    int
	Limit   int
	Keyword *string
	Desc    bool
	CoachId *uint
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
	if query.Keyword != nil {
		q = q.Where("name like ?", "%"+*query.Keyword+"%").
			Or("address like ?", "%"+*query.Keyword+"%").
			Or("summary like ?", "%"+*query.Keyword+"%")
	}
	if query.CoachId != nil {
		q = q.Where("coach_id = ?", *query.CoachId).
			Or(datatypes.JSONQuery("assistant_coach_ids").HasKey(strconv.Itoa(int(*query.CoachId))))
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
