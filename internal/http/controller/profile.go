package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gorm.io/datatypes"
	"net/http"
)

type ProfileResponse struct {
	ID              uint                   `json:"id"`
	Nickname        string                 `json:"nickname"`
	HeadImgURL      string                 `json:"headimgurl"`
	RoleNames       []string               `json:"roleNames"`
	UUID            string                 `json:"uuid"`
	Coach           ProfileCoach           `json:"coach"`
	TeachingRecords []model.TeachingRecord `json:"teachingRecords"` // 授课记录
	StudyRecords    []model.StudyRecord    `json:"studyRecords"`    // 学习记录
	MemberCourses   []*model.MemberCourse  `json:"memberCourses"`   // 会员课程
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

func (r UseCase) Profile(c *app.Context) result.Result {
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
	teachingRecords, err := r.Store.GetTeachingRecords(c.Context.Request.Context(), user.ID)
	if err != nil {
		return result.Err(fmt.Errorf("get user teaching records failed: %w", err))
	}
	studyRecords, err := r.Store.GetStudyRecords(c.Context.Request.Context(), user.ID)
	if err != nil {
		return result.Err(fmt.Errorf("get user study records failed: %w", err))
	}
	memberCourses, err := r.Store.GetMemberCourses(c.Context.Request.Context(), uint32(user.ID))
	if err != nil {
		return result.Err(fmt.Errorf("get user member course failed: %w", err))
	}
	return result.Ok(ProfileResponse{
		ID:         user.ID,
		Nickname:   user.Nickname,
		HeadImgURL: r.Storage.GetHeadImageVisitURL(user.HeadImgURL),
		RoleNames:  user.GetRoleNames(),
		UUID:       user.UUID,
		Coach: ProfileCoach{
			Level:               coach.Level,
			TeachingSpace:       coach.TeachingSpace,
			TeachingAge:         coach.TeachingAge,
			TeachingExperiences: coach.GetTeachingExperiences(),
		},
		TeachingRecords: teachingRecords,
		StudyRecords:    studyRecords,
		MemberCourses:   memberCourses,
	})
}
