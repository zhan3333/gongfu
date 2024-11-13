package client

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gongfu/internal/config"
	"gongfu/internal/model"
	"gongfu/pkg/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(conf *config.Config) (*gorm.DB, error) {
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

	return db, nil
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
