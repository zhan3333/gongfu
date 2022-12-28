package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gorm.io/datatypes"
	"net/http"
)

type ProfileResponse struct {
	ID         uint         `json:"id"`
	Nickname   string       `json:"nickname"`
	HeadImgURL string       `json:"headimgurl"`
	Role       string       `json:"role"`
	UUID       string       `json:"uuid"`
	Coach      ProfileCoach `json:"coach"`
}

type ProfileCoach struct {
	// 等级
	Level string `json:"level" json:"level,omitempty"`
	// 任教单位
	TeachingSpace string `json:"teachingSpace" json:"teaching_space,omitempty"`
	// 任教年限
	TeachingAge string `json:"teachingAge" json:"teaching_age,omitempty"`
	// 任教经历
	TeachingExperiences []string `json:"teachingExperiences" json:"teaching_experiences,omitempty"`
}

func (r Controller) Profile(c *app.Context) result.Result {
	uuid := c.Param("uuid")
	user, err := r.Store.GetUserByUUID(context.Background(), uuid)
	if err != nil {
		return result.Err(err)
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "user not found"})
		return result.Err(nil)
	}
	coach, err := r.Store.GetCoach(context.Background(), user.ID)
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
	return result.Ok(ProfileResponse{
		ID:         user.ID,
		Nickname:   user.Nickname,
		HeadImgURL: r.getUserHeadImgUrl(&c.User),
		Role:       user.Role,
		UUID:       user.UUID,
		Coach: ProfileCoach{
			Level:               coach.Level,
			TeachingSpace:       coach.TeachingSpace,
			TeachingAge:         coach.TeachingAge,
			TeachingExperiences: coach.GetTeachingExperiences(),
		},
	})
}
