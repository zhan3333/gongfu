package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	util "github.com/zhan3333/goutil"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	util2 "gongfu/pkg/util"
	"net/http"
	"strconv"
)

// GetCourses 获取所有课程信息
func (r UseCase) GetCourses(c *app.Context) result.Result {
	courses, err := r.Store.GetCoursesByUser(context.Background(), c.UserID)
	if err != nil {
		return result.Err(err)
	}
	var userIds = util.Reduce(courses, func(arr []uint, course *model.Course) []uint {
		return append(arr, course.GetRelatedUserIds()...)
	})
	usersQuery, err := r.NewUsersQuery(userIds...)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(util.Map(courses, func(course *model.Course) gin.H {
		// todo 优化查询
		school, _ := r.Store.GetSchool(context.TODO(), course.SchoolId)
		return gin.H{
			"id":               course.ID,
			"createdAt":        course.CreatedAt.Unix(),
			"startDate":        course.StartDate,
			"startTime":        course.StartTime,
			"schoolId":         course.SchoolId,
			"school":           school,
			"manager":          usersQuery.GetCoach(course.ManagerId),
			"coach":            usersQuery.GetCoach(course.CoachId),
			"assistantCoaches": usersQuery.GetCoaches(course.GetAssistantCoachIds()...),
			"checkInBy":        usersQuery.GetCoach(course.CheckInBy),
			"checkOutBy":       usersQuery.GetCoach(course.CheckOutBy),
			"checkInAt":        course.CheckInAt.Time.Unix(),
			"checkOutAt":       course.CheckOutAt.Time.Unix(),
			"images":           course.GetImages(),
			"summary":          course.Summary,
		}
	}))
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
	school, err := r.Store.GetSchool(context.TODO(), course.SchoolId)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(gin.H{
		"id":               course.ID,
		"startDate":        course.StartDate,
		"startTime":        course.StartTime,
		"coach":            userQuery.GetCoach(course.CoachId),
		"schoolId":         course.SchoolId,
		"school":           school,
		"manager":          userQuery.GetCoach(course.ManagerId),
		"content":          course.Content,
		"summary":          course.Summary,
		"images":           course.GetImages(),
		"assistantCoaches": userQuery.GetCoaches(course.GetAssistantCoachIds()...),
		"checkInAt":        util2.DBTimeToTimestamp(course.CheckInAt),
		"checkOutAt":       util2.DBTimeToTimestamp(course.CheckOutAt),
		"checkInBy":        userQuery.GetCoach(course.CheckInBy),
		"checkOutBy":       userQuery.GetCoach(course.CheckOutBy),
	})
}

func (r UseCase) UpdateCourse(c *app.Context) result.Result {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return result.Err(err)
	}
	course, err := r.Store.GetCourse(context.TODO(), uint(id))
	if err != nil {
		return result.Err(err)
	}
	if course == nil {
		c.String(http.StatusNotFound, "course not found")
		return result.Err(nil)
	}
	req := struct {
		Content           string
		Summary           string
		Images            []string `form:"images"`            // 图片 key 列表
		ManagerId         uint     `form:"managerId"`         // 学校
		CoachId           uint     `form:"coachId"`           // 教练
		AssistantCoachIds []uint   `form:"assistantCoachIds"` // 助理教练列表
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}

	course.Images = req.Images
	course.Content = req.Content
	course.Summary = req.Summary
	course.ManagerId = req.ManagerId
	course.CoachId = req.CoachId
	course.AssistantCoachIds = req.AssistantCoachIds

	if err := r.Store.UpdateCourse(context.TODO(), course); err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}
