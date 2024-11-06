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
)

// GetSchools 获取学校列表
func (r UseCase) GetSchools(c *app.Context) result.Result {
	schools, err := r.Store.GetSchools(context.TODO())
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(schools)
}

// CreateCourse 创建课程
func (r UseCase) CreateCourse(c *app.Context) result.Result {
	req := struct {
		StartDate         string `form:"startDate"`         // 上课日期
		StartTime         string `form:"startTime"`         // 上课时间
		SchoolId          uint   `form:"schoolId"`          // 学校
		ManagerId         uint   `form:"managerId"`         // 学校
		CoachId           uint   `form:"coachId"`           // 教练
		AssistantCoachIds []uint `form:"assistantCoachIds"` // 助理教练列表
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(err)
	}
	var err = r.Store.CreateCourse(context.TODO(), store.CreateCourseInput{
		StartDate:         req.StartDate,
		StartTime:         req.StartTime,
		SchoolId:          req.SchoolId,
		ManagerId:         req.ManagerId,
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
		Page    int     `json:"page" form:"page"`
		Limit   int     `json:"limit" form:"limit"`
		Keyword *string `json:"keyword" form:"keyword"`
		Desc    bool    `json:"desc" form:"desc"`
		UserId  *uint   `json:"user_id" form:"userId"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	page, err := r.Store.GetCoursePage(context.Background(), store.GetCoursePageInput{
		Page:    req.Page,
		Limit:   req.Limit,
		Keyword: req.Keyword,
		Desc:    req.Desc,
		UserId:  req.UserId,
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
	var schoolIds = util.Reduce(page.Items, func(arr []uint, course *model.Course) []uint {
		return append(arr, course.SchoolId)
	})
	schoolsQuery, err := r.Store.SchoolMap(c.Ctx(), schoolIds...)
	if err != nil {
		return result.Err(err)
	}

	return result.Ok(gin.H{
		"items": util.Map(page.Items, func(course *model.Course) gin.H {
			return gin.H{
				"id":               course.ID,
				"createdAt":        course.CreatedAt.Unix(),
				"startDate":        course.StartDate,
				"startTime":        course.StartTime,
				"schoolId":         course.SchoolId,
				"school":           schoolsQuery.Name(course.SchoolId),
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
	school, err := r.Store.GetSchool(context.TODO(), course.SchoolId)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(gin.H{
		"id":                course.ID,
		"startDate":         course.StartDate,
		"startTime":         course.StartTime,
		"coach":             userQuery.GetCoach(course.CoachId),
		"schoolId":          course.SchoolId,
		"school":            school,
		"manager":           userQuery.GetCoach(course.ManagerId),
		"summary":           course.Summary,
		"content":           course.Content,
		"images":            course.GetImages(),
		"assistantCoaches":  userQuery.GetCoaches(course.GetAssistantCoachIds()...),
		"assistantCoachIds": course.GetAssistantCoachIds(),
		"checkInAt":         util2.DBTimeToTimestamp(course.CheckInAt),
		"checkOutAt":        util2.DBTimeToTimestamp(course.CheckOutAt),
		"checkInBy":         userQuery.GetCoach(course.CheckInBy),
		"checkOutBy":        userQuery.GetCoach(course.CheckOutBy),
	})
}

// DeleteCourse 删除课程
func (r UseCase) DeleteCourse(c *app.Context) result.Result {
	s := c.Param("id")
	id, err := strconv.Atoi(s)
	if err != nil {
		return result.Err(err)
	}
	err = r.Store.DeleteCourse(context.TODO(), uint(id))
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

func (r UseCase) UpdateCourse(c *app.Context) result.Result {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return result.Err(err)
	}
	req := struct {
		StartDate         string   `form:"startDate"`         // 上课日期
		StartTime         string   `form:"startTime"`         // 上课时间
		SchoolId          uint     `form:"schoolId"`          // 学校
		ManagerId         uint     `form:"managerId"`         // 学校
		CoachId           *uint    `form:"coachId"`           // 教练
		AssistantCoachIds []uint   `form:"assistantCoachIds"` // 助理教练列表
		Images            []string `form:"images"`
		Content           string   `form:"content"`
		Summary           string   `form:"summary"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	course, err := r.Store.GetCourse(context.TODO(), uint(id))
	if err != nil {
		return result.Err(err)
	}
	if course == nil {
		c.String(http.StatusNotFound, "course not found")
		return result.Err(nil)
	}
	course.StartDate = req.StartDate
	course.StartTime = req.StartTime
	course.SchoolId = req.SchoolId
	course.ManagerId = req.ManagerId
	if req.CoachId != nil {
		course.CoachId = *req.CoachId
	}
	course.Images = req.Images
	course.Summary = req.Summary
	course.Content = req.Content
	course.AssistantCoachIds = req.AssistantCoachIds

	if err := r.Store.UpdateCourse(c.Ctx(), course); err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

func (r UseCase) GetCourseList(c *app.Context) result.Result {
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		return result.Err(err)
	}

	courses, err := r.Store.GetCoursesByUser(context.Background(), uint(userId))
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
	var schoolIds = util.Reduce(courses, func(arr []uint, course *model.Course) []uint {
		return append(arr, course.SchoolId)
	})
	schoolMap, err := r.Store.SchoolMap(c.Ctx(), schoolIds...)
	if err != nil {
		return result.Err(err)
	}

	return result.Ok(util.Map(courses, func(course *model.Course) gin.H {
		return gin.H{
			"id":             course.ID,
			"schoolName":     schoolMap.Name(course.ID),
			"startDate":      course.StartDate,
			"startTime":      course.StartTime,
			"managerName":    usersQuery.Name(course.ManagerId),
			"coachName":      usersQuery.Name(course.CoachId),
			"assistantNames": usersQuery.Names(course.GetAssistantCoachIds()...),
			"checkInBy":      usersQuery.Name(course.CheckInBy),
			"checkOutBy":     usersQuery.Name(course.CheckOutBy),
			"checkInAt":      util2.NullTime(course.CheckInAt),
			"checkOutAt":     util2.NullTime(course.CheckOutAt),
			"summary":        course.Summary,
			"content":        course.Content,
		}
	}))
}
