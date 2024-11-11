package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gongfu/internal/config"
	"gongfu/internal/model"
	"gongfu/pkg/util"
	"gorm.io/driver/mysql"
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

type CreateEnrollInput struct {
	UserId      uint
	ActivityId  string
	Status      string
	Amount      int64
	OutTradeNo  string
	Description string
	Attach      string
	Username    string
	Phone       string
	Sex         string
}

type SchoolMap map[uint]*model.School

func (receiver SchoolMap) Name(schoolId uint) string {
	s, ok := receiver[schoolId]
	if !ok {
		return "unknown"
	}
	return s.Name
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

	CreateEnroll(ctx context.Context, in *CreateEnrollInput) error
	UpdateEnroll(ctx context.Context, enroll *model.Enroll) error
	GetEnroll(ctx context.Context, outTradeNo string) (*model.Enroll, error)
	GetEnrollByActivityId(ctx context.Context, activityId string) (*model.Enroll, error)

	SchoolMap(ctx context.Context, schoolIds ...uint) (SchoolMap, error)
}

var _ Store = (*DBStore)(nil)

type DBStore struct {
	DB *gorm.DB
}

func (s DBStore) SchoolMap(ctx context.Context, schoolIds ...uint) (SchoolMap, error) {
	var schools []*model.School
	err := s.DB.WithContext(ctx).Where("id in ?", schoolIds).Find(&schools).Error
	if err != nil {
		return nil, err
	}
	m := make(SchoolMap)
	for _, s := range schools {
		m[s.ID] = s
	}
	return m, nil
}

func (s DBStore) UpdateEnroll(ctx context.Context, enroll *model.Enroll) error {
	return s.DB.WithContext(ctx).Save(enroll).Error
}

func (s DBStore) GetEnroll(ctx context.Context, outTradeNo string) (*model.Enroll, error) {
	var enroll model.Enroll
	err := s.DB.WithContext(ctx).Where("out_trade_no = ?", outTradeNo).First(&enroll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &enroll, nil
}

func (s DBStore) GetEnrollByActivityId(ctx context.Context, activityId string) (*model.Enroll, error) {
	var enroll model.Enroll
	err := s.DB.WithContext(ctx).Where("activity_id = ?", activityId).Order("created_at desc").First(&enroll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &enroll, nil
}

func (s DBStore) CreateEnroll(ctx context.Context, in *CreateEnrollInput) error {
	return s.DB.WithContext(ctx).Create(&model.Enroll{
		ActivityId:  in.ActivityId,
		Status:      in.Status,
		Amount:      in.Amount,
		UserId:      in.UserId,
		OutTradeNo:  in.OutTradeNo,
		Description: in.Description,
		Attach:      in.Attach,
		Username:    in.Username,
		Phone:       in.Phone,
		Sex:         in.Sex,
	}).Error
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

func NewStore(conf *config.Config) (Store, error) {
	db, err := gorm.Open(mysql.Open(conf.DB.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.AutoMigrate(
		&model.User{},
		&model.CheckIn{},
		&model.CheckInDay{},
		&model.CheckInCount{},
		&model.Coach{},
		&model.Migrate{},
		&model.Role{},
		&model.UserHasRole{},
		&model.School{},
		&model.Course{},
		&model.TeachingRecord{},
		&model.StudyRecord{},
		&model.MemberCourse{},
		&model.CheckInComment{},
		&model.Enroll{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	migrates := map[string]func(db2 *gorm.DB) error{
		"init_users_uuid": initUsersUUID,
		"init_roles":      initRoles,
		"init_schools":    initSchools,
	}
	for action, migrate := range migrates {
		var version = model.Migrate{}
		err = db.Model(&model.Migrate{}).Where("action = ?", action).First(&version).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("find max version: %w", err)
			}
			// 未找到记录
			if err := migrate(db); err != nil {
				return nil, fmt.Errorf("migrate action=%s: %w", action, err)
			}
			version.Action = action
			if err = db.Create(&version).Error; err != nil {
				return nil, fmt.Errorf("create migrate: %w", err)
			}
			logrus.Infof("migrate action=%s migrated success", action)
		} else {
			logrus.Infof("migrate action=%s already migrated", action)
		}
	}

	return &DBStore{DB: db}, nil
}

func initSchools(db *gorm.DB) error {
	schools := []*model.School{
		{
			Name:    "丁字桥小学",
			Address: "",
		},
		{
			Name:    "水果湖一小",
			Address: "",
		},
		{
			Name:    "洪山区第一小学 书城路分校",
			Address: "",
		},
		{
			Name:    "陆家街中学",
			Address: "",
		},
		{
			Name:    "粮道街中学 三角路校区",
			Address: "",
		},
		{
			Name:    "粮道街中学 积玉桥",
			Address: "",
		},
		{
			Name:    "粮道街中学 本部",
			Address: "",
		},
	}
	if err := db.Create(&schools).Error; err != nil {
		return fmt.Errorf("create schools failed: %w", err)
	}
	return nil
}

func initRoles(db *gorm.DB) error {
	roles := []*model.Role{
		{
			ID:   1,
			Name: "admin",
		},
		{
			ID:   2,
			Name: "coach",
		},
		{
			ID:   3,
			Name: "user",
		},
	}
	if err := db.Create(&roles).Error; err != nil {
		return fmt.Errorf("create roles: %w", err)
	}

	// 为所有用户赋予普通用户角色
	users := []*model.User{}
	err := db.FindInBatches(&users, 2000, func(tx *gorm.DB, batch int) error {
		var userRoles []*model.UserHasRole
		for _, user := range users {
			userRoles = append(userRoles, &model.UserHasRole{
				UserID:   user.ID,
				RoleID:   3,
				RoleName: "user",
			})
		}
		if err := db.Create(userRoles).Error; err != nil {
			return fmt.Errorf("create user roles: %w", err)
		}
		return nil
	}).Error
	if err != nil {
		return fmt.Errorf("set users roles: %w", err)
	}
	// 删除 role 列
	if err := db.Migrator().DropColumn(&model.User{}, "role"); err != nil {
		return fmt.Errorf("drop user role column: %w", err)
	}
	return nil
}

func initUsersUUID(db *gorm.DB) error {
	// 为所有用户加上 uuid
	users := []*model.User{}
	err := db.Where("uuid is null").FindInBatches(&users, 2000, func(tx *gorm.DB, batch int) error {
		for _, user := range users {
			user.UUID = util.UUID()
			if err := db.Save(user).Error; err != nil {
				return fmt.Errorf("update user: %w", err)
			}
		}
		return nil
	}).Error
	if err != nil {
		return fmt.Errorf("set users uuid: %w", err)
	}
	return nil
}
