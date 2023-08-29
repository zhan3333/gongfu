package controller

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gongfu/internal/app"
	"gongfu/internal/result"
	"net/http"
)

// GetBindCode 获取绑定验证码
func (r UseCase) GetBindCode(c *app.Context) result.Result {
	var req = struct {
		Phone string `binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.Message(http.StatusBadRequest, fmt.Sprintf("bad request params: %s", err))
		return result.Err(nil)
	}
	if c.User.Phone != nil {
		c.Message(http.StatusBadRequest, "you already bind phone")
		return result.Err(nil)
	}
	if code, err := r.AuthCode.Send(req.Phone); err != nil {
		return result.Err(err)
	} else {
		logrus.WithField("valid_code", code).WithField("phone", req.Phone).Info("get_bind_code")
	}
	return result.Ok(nil)
}

// ValidBindCode 验证绑定验证码
func (r UseCase) ValidBindCode(c *app.Context) result.Result {
	var req = struct {
		Phone string `binding:"required"`
		Code  string `binding:"required"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.Message(http.StatusBadRequest, fmt.Sprintf("bad request params: %s", err))
		return result.Err(nil)
	}
	if c.User.Phone != nil {
		c.Message(http.StatusBadRequest, "already bind phone")
		return result.Err(nil)
	}
	if user, err := r.Store.GetUserByOpenID(context.TODO(), req.Phone); err != nil {
		return result.Err(fmt.Errorf("phone query user: %w", err))
	} else if user != nil {
		c.Message(http.StatusBadRequest, "phone already used")
		return result.Err(nil)
	}
	if valid, err := r.AuthCode.Valid(req.Phone, req.Code); err != nil {
		return result.Err(err)
	} else {
		if !valid {
			c.Message(http.StatusBadRequest, "valid code invalid")
			return result.Err(nil)
		} else {
			c.User.Phone = &req.Phone
			if err := r.Store.UpdateUser(context.TODO(), &c.User); err != nil {
				return result.Err(err)
			}
			return result.Ok(nil)
		}
	}
}
