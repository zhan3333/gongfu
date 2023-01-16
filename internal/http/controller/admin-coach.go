package controller

import (
	"context"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gorm.io/datatypes"
	"strconv"
)

func (r Controller) AdminGetCoach(c *app.Context) result.Result {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return result.Err(err)
	}
	coach, err := r.Store.GetCoach(context.Background(), uint(userID))
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
