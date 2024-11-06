package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

import "gorm.io/datatypes"

type Coach struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id" gorm:"index"`
	// 等级
	Level string `json:"level"`
	// 任教单位
	TeachingSpace string `json:"teaching_space"`
	// 任教年限
	TeachingAge string `json:"teaching_age"`
	// 任教经历
	TeachingExperiences datatypes.JSON `json:"teaching_experiences"`
}

func (c Coach) GetTeachingExperiences() []string {
	v := []string{}
	_ = json.Unmarshal(c.TeachingExperiences, &v)
	return v
}
