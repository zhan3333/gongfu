package controller

import (
	"errors"
	"gongfu/internal/app"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"strconv"
)

// ChangeMemberCourseRemain change member course remain
func (r UseCase) ChangeMemberCourseRemain(c *app.Context) result.Result {
	if !c.User.HasAnyRole([]string{model.ROLE_ADMIN, model.ROLE_COACH}) {
		return result.Err(errors.New("permission denied"))
	}
	memberCourseId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return result.Err(err)
	}
	if memberCourseId <= 0 {
		return result.Err(errors.New("invalid member course id"))
	}
	// check remain value
	mc, err := r.Store.GetMemberCourse(c.Request.Context(), uint(memberCourseId))
	if err != nil {
		return result.Err(err)
	}
	if mc == nil {
		return result.Err(errors.New("member course not found"))
	}

	req := struct {
		Remain int `json:"remain" binding:"required,omitempty"`
	}{}
	if err := c.Bind(&req); err != nil {
		return result.Err(nil)
	}
	if req.Remain < 0 {
		return result.Err(errors.New("invalid remain"))
	}
	if mc.Total < req.Remain {
		return result.Err(errors.New("remain should not be greater than total"))
	}
	err = r.Store.ChangeMemberCourseRemain(c.Request.Context(), uint(memberCourseId), req.Remain)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(nil)
}
