package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	util "github.com/zhan3333/goutil"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gongfu/internal/service/store"
	util2 "gongfu/pkg/util"
	"net/http"
	"strconv"
	"time"
)

// CreateCourse 创建课程
func (r UseCase) CreateCourse(c *app.Context) result.Result {
	req := struct {
		SchoolStartAt     int64  `form:"schoolStartAt"`     // 上课时间
		Address           string `form:"address"`           // 上课地点
		Name              string `form:"name"`              // 课程名称
		CoachId           uint   `form:"coachId"`           // 教练
		AssistantCoachIds []uint `form:"assistantCoachIds"` // 助理教练列表
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(err)
	}
	startAt := time.Unix(req.SchoolStartAt, 0)
	err := r.Store.CreateCourse(context.TODO(), store.CreateCourseInput{
		SchoolStartAt:     startAt,
		Address:           req.Address,
		Name:              req.Name,
		CoachId:           req.CoachId,
		AssistantCoachIds: req.AssistantCoachIds,
	})
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

// GetCoursePage 获取 course 分页
func (r UseCase) GetCoursePage(c *app.Context) result.Result {
	req := struct {
		Page     int     `json:"page" form:"page"`
		Limit    int     `json:"limit" form:"limit"`
		Keyword  *string `json:"keyword" form:"keyword"`
		Desc     bool    `json:"desc" form:"desc"`
		CourseId *uint   `json:"course_id" form:"courseId"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	page, err := r.Store.GetCoursePage(context.Background(), store.GetCoursePageInput{
		Page:    req.Page,
		Limit:   req.Limit,
		Keyword: req.Keyword,
		Desc:    req.Desc,
		CoachId: req.CourseId,
	})
	if err != nil {
		return result.Err(err)
	}
	var userIds = util.Reduce(page.Items, func(arr []uint, course *model.Course) []uint {
		return append(arr, course.GetRelatedUserIds()...)
	})
	usersQuery, err := r.NewUsersQuery(userIds...)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(gin.H{
		"items": util.Map(page.Items, func(course *model.Course) gin.H {
			return gin.H{
				"id":               course.ID,
				"createdAt":        course.CreatedAt.Unix(),
				"schoolStartAt":    course.SchoolStartAt.Time.Unix(),
				"address":          course.Address,
				"name":             course.Name,
				"coach":            usersQuery.GetCoach(course.CoachId),
				"assistantCoaches": usersQuery.GetCoaches(course.GetAssistantCoachIds()...),
				"checkInBy":        usersQuery.GetCoach(course.CheckInBy),
				"checkOutBy":       usersQuery.GetCoach(course.CheckOutBy),
				"checkInAt":        course.CheckInAt.Time.Unix(),
				"checkOutAt":       course.CheckOutAt.Time.Unix(),
				"images":           course.GetImages(),
				"summary":          course.Summary,
			}
		}),
		"page":  page.Page,
		"count": page.Count,
		"limit": page.Limit,
	})
}

func (r UseCase) GetCourse(c *app.Context) result.Result {
	s := c.Param("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		return result.Err(err)
	}
	course, err := r.Store.GetCourse(context.TODO(), uint(id))
	if err != nil {
		return result.Err(err)
	}
	if course == nil {
		c.String(http.StatusNotFound, "user not found")
		return result.Err(nil)
	}
	userQuery, err := r.NewUsersQuery(course.GetRelatedUserIds()...)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(gin.H{
		"id":               course.ID,
		"schoolStartAt":    util2.DBTimeToTimestamp(course.SchoolStartAt),
		"name":             course.Name,
		"coach":            userQuery.GetCoach(course.CoachId),
		"address":          course.Address,
		"summary":          course.Summary,
		"images":           course.GetImages(),
		"assistantCoaches": userQuery.GetCoaches(course.GetAssistantCoachIds()...),
		"checkInAt":        util2.DBTimeToTimestamp(course.CheckInAt),
		"checkOutAt":       util2.DBTimeToTimestamp(course.CheckOutAt),
		"checkInBy":        userQuery.GetCoach(course.CheckInBy),
		"checkOutBy":       userQuery.GetCoach(course.CheckOutBy),
	})
}
