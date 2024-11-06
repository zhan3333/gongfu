package model

import (
	"database/sql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Course struct {
	ID                uint                        `json:"id" gorm:"primarykey"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at"`
	DeletedAt         gorm.DeletedAt              `json:"deleted_at" gorm:"index"`
	StartDate         string                      `json:"start_date" gorm:"index"` // 上课日期
	StartTime         string                      `json:"start_time"`              // 上课时间
	SchoolId          uint                        `json:"school_id"`               // 学校 id
	ManagerId         uint                        `json:"manager_id"`              // 负责人
	CoachId           uint                        `json:"coach_id" gorm:"index"`   // 教练
	AssistantCoachIds datatypes.JSONSlice[uint]   `json:"assistant_coach_ids"`     // 助理教练列表
	CheckInBy         uint                        `json:"check_in_by"`             // 签到者
	CheckOutBy        uint                        `json:"check_out_by"`            // 签出者
	CheckInAt         sql.NullTime                `json:"check_in_at"`             // 签到时间
	CheckOutAt        sql.NullTime                `json:"check_out_at"`            // 签出时间
	Images            datatypes.JSONSlice[string] `json:"images"`                  // 图片列表
	Summary           string                      `json:"summary"`                 // 总结
	Content           string                      `json:"content"`                 // 内容
}

func (c Course) GetImages() []string {
	return c.Images
}

func (c Course) GetAssistantCoachIds() []uint {
	return c.AssistantCoachIds
}

func (c Course) GetRelatedUserIds() []uint {
	var userIds = []uint{c.ManagerId}
	if c.CoachId != 0 {
		userIds = append(userIds, c.CoachId)
	}
	if c.CheckInBy != 0 {
		userIds = append(userIds, c.CheckInBy)
	}
	if c.CheckOutBy != 0 {
		userIds = append(userIds, c.CheckOutBy)
	}
	if len(c.GetAssistantCoachIds()) != 0 {
		userIds = append(userIds, c.GetAssistantCoachIds()...)
	}
	return userIds
}
