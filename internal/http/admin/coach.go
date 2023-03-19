package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gorm.io/datatypes"
	"strconv"
)

func (r UseCase) AdminGetCoach(c *app.Context) result.Result {
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

func (r UseCase) GetCoaches(c *app.Context) result.Result {
	coaches, err := r.Store.GetCoaches(context.Background())
	if err != nil {
		return result.Err(err)
	}
	var ret = []gin.H{}
	for _, coach := range coaches {
		ret = append(ret, gin.H{
			"id":   coach.ID,
			"name": coach.Nickname,
		})
	}
	return result.Ok(ret)
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
