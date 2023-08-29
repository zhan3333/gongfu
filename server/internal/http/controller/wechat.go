package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/result"
)

func (r UseCase) JSConfig(c *app.Context) result.Result {
	req := struct {
		Uri string `form:"uri" binding:"required"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		return result.Err(err)
	}
	config, err := r.OfficialAccount.GetJs().GetConfig(req.Uri)
	if err != nil {
		return result.Err(fmt.Errorf("get wechat js config: %w", err))
	}
	return result.Ok(gin.H{
		"appID":     config.AppID,
		"timestamp": config.Timestamp,
		"nonceStr":  config.NonceStr,
		"signature": config.Signature,
	})
}
