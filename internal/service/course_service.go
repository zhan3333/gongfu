package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gongfu/internal/client"
	"gongfu/internal/model"
	"gongfu/internal/service/store"
	"gongfu/pkg/util"
	"gorm.io/gorm"
)

type CourseService struct {
	store           store.Store
	db              *gorm.DB
	officialAccount *client.OfficialAccount
}

func NewCourseService(store store.Store, db *gorm.DB, account *client.OfficialAccount) *CourseService {
	return &CourseService{store, db, account}
}

type CreateCourseInput struct {
	StartDate         string `json:"start_date"`          // 上课日期
	StartTime         string `json:"start_time"`          // 上课时间
	CoachId           uint   `json:"coach_id"`            // 教练
	SchoolId          uint   `json:"school_id"`           // 学校
	AssistantCoachIds []uint `json:"assistant_coach_ids"` // 助理教练列表
	ManagerId         uint   `json:"manager_id"`          // 负责人
}

func (s CourseService) CreateCourse(ctx context.Context, in *CreateCourseInput) error {
	course := model.Course{
		StartTime:         in.StartTime,
		StartDate:         in.StartDate,
		SchoolId:          in.SchoolId,
		CoachId:           in.CoachId,
		AssistantCoachIds: in.AssistantCoachIds,
		Images:            []string{},
		Summary:           "",
		ManagerId:         in.ManagerId,
	}
	err := s.db.Model(&model.Course{}).WithContext(ctx).Create(&course).Error
	if err != nil {
		return fmt.Errorf("create course: %w", err)
	}
	go func() {
		ctx2 := context.Background()
		var userIds []uint
		userIds = append(userIds, in.CoachId)
		userIds = append(userIds, in.ManagerId)
		userIds = append(userIds, in.AssistantCoachIds...)
		userMap, err := s.store.GetUsersMap(ctx2, userIds)
		if err != nil {
			logrus.WithError(err).WithContext(ctx2).Error("get users map failed")
			return
		}
		schoolName, err := s.store.GetSchoolName(ctx2, in.SchoolId)
		if err != nil {
			logrus.WithError(err).WithContext(ctx2).Error("get school name failed")
			return
		}

		var toolTip string
		toolTip += fmt.Sprintf("教练:%s", userMap.Name(in.CoachId))
		toolTip += fmt.Sprintf("\n负责人:%s", userMap.Name(in.ManagerId))
		for _, id := range in.AssistantCoachIds {
			toolTip += fmt.Sprintf("\n助理教练:%s", userMap.Name(id))
		}

		sin := &client.SubscribeSendInput{
			CourseName:    "您有一节待上的课程",
			CourseTime:    fmt.Sprintf("%s %s", in.StartDate, in.StartTime),
			CourseAddress: schoolName,
			TeacherName:   userMap.Name(in.CoachId),
			ToolTip:       toolTip,
		}

		var sentUserIds = map[uint]bool{}

		if coachOpenId := userMap.OpenId(in.CoachId); !util.Empty(coachOpenId) {
			sin.OpenId = *coachOpenId
			sentUserIds[in.CoachId] = true
			if err := s.officialAccount.SubscribeSend(sin); err != nil {
				logrus.WithError(err).WithContext(ctx2).Error("subscribe send to coach failed")
			}
		}
		if managerOpenId := userMap.OpenId(in.ManagerId); !util.Empty(managerOpenId) && !sentUserIds[in.CoachId] {
			sin.OpenId = *managerOpenId
			sentUserIds[in.ManagerId] = true
			if err := s.officialAccount.SubscribeSend(sin); err != nil {
				logrus.WithError(err).WithContext(ctx2).Error("subscribe send to manager failed")
			}
		}
		for _, id := range in.AssistantCoachIds {
			if assistantOpenId := userMap.OpenId(id); !util.Empty(assistantOpenId) && !sentUserIds[id] {
				sin.OpenId = *assistantOpenId
				sentUserIds[id] = true
				if err := s.officialAccount.SubscribeSend(sin); err != nil {
					logrus.WithError(err).WithContext(ctx2).Error("subscribe send to assistant failed")
				}
			}
		}
		logrus.WithContext(ctx2).Info("subscribe send end")
	}()
	return nil
}
