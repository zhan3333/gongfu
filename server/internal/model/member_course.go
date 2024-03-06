package model

import "time"

type MemberCourse struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name" gorm:"column:name;type:varchar(255);not null;default:'';comment:课程名称"`
	// userId
	UserId uint `json:"userId" gorm:"column:user_id;type:bigint(32);comment:用户 id"`
	// 开始时间
	StartTime time.Time `json:"startTime" gorm:"column:start_time;type:timestamp;comment:开始时间"`
	// 结束时间
	EndTime time.Time `json:"endTime" gorm:"column:end_time;type:timestamp;comment:结束时间"`
	// 总课程数
	Total int `json:"total" gorm:"column:total;type:bigint(32);comment:总课程数"`
	// 剩余课程数
	Remain int `json:"remain" gorm:"column:remain;type:bigint(32);comment:剩余课程数"`
	// 备注
	Remark string `json:"remark" gorm:"column:remark;type:varchar(255);not null;default:'';comment:备注"`
	// 状态
	Status string `json:"status" gorm:"column:status;type:varchar(50);not null;default:'';comment:状态"`
}

func (m MemberCourse) TableName() string {
	return "member_courses"
}
