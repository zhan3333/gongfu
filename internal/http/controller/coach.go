package controller

import (
	"context"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gorm.io/datatypes"
)

func (r UseCase) GetCoach(c *app.Context) result.Result {
	coach, err := r.Store.GetCoach(context.Background(), c.UserID)
	if err != nil {
		return result.Err(err)
	}
	if coach == nil {
		coach = &model.Coach{
			Level:               "",
			TeachingSpace:       "",
			TeachingAge:         "",
			TeachingExperiences: datatypes.JSON{},
		}
	}
	if coach.TeachingExperiences == nil {
		coach.TeachingExperiences = datatypes.JSON{}
	}
	return result.Ok(CoachResponse{
		ID:                  coach.ID,
		UserID:              coach.UserID,
		Level:               coach.Level,
		TeachingSpace:       coach.TeachingSpace,
		TeachingAge:         coach.TeachingAge,
		TeachingExperiences: coach.GetTeachingExperiences(),
	})
}

type CoachResponse struct {
	ID     uint
	UserID uint
	// 等级
	Level string `json:"level"`
	// 任教单位
	TeachingSpace string `json:"teachingSpace"`
	// 任教年限
	TeachingAge string `json:"teachingAge"`
	// 任教经历
	TeachingExperiences []string `json:"teachingExperiences"`
}
