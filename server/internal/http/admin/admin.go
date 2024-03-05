package admin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/http/common"
	"gongfu/internal/result"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
	"strconv"
	"time"
)

type UseCase struct {
	Store   store.Store
	Storage service.Storage
}

func NewUseCase(store store.Store, storage service.Storage) UseCase {
	return UseCase{
		Store:   store,
		Storage: storage,
	}
}

func (r UseCase) NewUsersQuery(userIds ...uint) (*common.UsersQuery, error) {
	usersMap, err := r.Store.GetUsersMap(context.TODO(), userIds)
	if err != nil {
		return nil, fmt.Errorf("get users map: %w", err)
	}
	return &common.UsersQuery{UsersMap: usersMap, Storage: r.Storage}, nil
}

func (r UseCase) GetMemberCourses(c *app.Context) result.Result {
	userID, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		return result.Err(err)
	}
	if userID == 0 {
		return result.Err(fmt.Errorf("invalid user id"))
	}

	courses, err := r.Store.GetMemberCourses(c.Request.Context(), uint32(userID))
	if err != nil {
		return result.Err(err)
	}
	var res []gin.H
	for _, v := range courses {
		res = append(res, gin.H{
			"id":        v.ID,
			"name":      v.Name,
			"startTime": v.StartTime,
			"endTime":   v.EndTime,
			"total":     v.Total,
			"remain":    v.Remain,
			"remark":    v.Remark,
			"status":    v.Status,
		})
	}
	return result.Ok(res)
}

type CreateMemberCourseRequest struct {
	UserID    uint32    `form:"useId" json:"userId" binding:"required"`
	Name      string    `form:"name" json:"name" binding:"required"`
	StartTime time.Time `form:"startTime" json:"startTime" binding:"required"`
	EndTime   time.Time `form:"endTime" json:"endTime" binding:"required"`
	Total     int       `form:"total" json:"total" binding:"required"`
	Remark    string    `form:"remark" json:"remark"`
}

func (r UseCase) CreateMemberCourse(c *app.Context) result.Result {
	var req CreateMemberCourseRequest
	if err := c.Bind(&req); err != nil {
		return result.Err(err)
	}
	if req.Total <= 0 {
		return result.Err(fmt.Errorf("invalid total"))
	}
	err := r.Store.CreateMemberCourse(c.Request.Context(), &store.CreateMemberCourseInput{
		UserId:    req.UserID,
		Name:      req.Name,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Total:     req.Total,
		Remain:    req.Total,
		Remark:    req.Remark,
		Status:    "normal",
	})
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

func (r UseCase) DeleteMemberCourse(c *app.Context) result.Result {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return result.Err(err)
	}
	err = r.Store.DeleteMemberCourse(c.Request.Context(), uint(id))
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}

type UpdateMemberCourseRequest struct {
	Name      string    `form:"name" json:"name" binding:"required"`
	StartTime time.Time `form:"startTime" json:"startTime" binding:"required"`
	EndTime   time.Time `form:"endTime" json:"endTime" binding:"required"`
	Total     int       `form:"total" json:"total" binding:"required"`
	Remain    int       `form:"remain" json:"remain" binding:"required"`
	Remark    string    `form:"remark" json:"remark"`
	Status    string    `form:"status" json:"status" binding:"required"`
}

func (r UseCase) UpdateMemberCourse(c *app.Context) result.Result {
	var req UpdateMemberCourseRequest
	if err := c.Bind(&req); err != nil {
		return result.Err(err)
	}
	if req.Total <= 0 {
		return result.Err(fmt.Errorf("invalid total"))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return result.Err(err)
	}
	err = r.Store.UpdateMemberCourse(c.Request.Context(), uint(id), &store.UpdateMemberCourseInput{
		Name:      req.Name,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Total:     req.Total,
		Remain:    req.Remain,
		Remark:    req.Remark,
		Status:    req.Status,
	})
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}
