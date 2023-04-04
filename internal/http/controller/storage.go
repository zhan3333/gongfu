package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"gongfu/internal/app"
	"gongfu/internal/result"
	"gongfu/pkg/util"
	"path"
	"path/filepath"
	"strings"
)

// GetUploadToken 获取上传文件的 token
func (r UseCase) GetUploadToken(c *app.Context) result.Result {
	req := struct {
		FileName string `form:"fileName" binding:"required"`
		Dir      string `form:"dir"`
	}{}
	if err := c.ShouldBindQuery(&req); err != nil {
		return result.Err(err)
	}
	key := util.UUID() + filepath.Ext(req.FileName)
	if req.Dir != "" {
		key = path.Join(req.Dir, key)
	}
	uploadUrl, err := r.Storage.GetPresignedURL(context.TODO(), key)
	if err != nil {
		return result.Err(err)
	}
	return result.Ok(gin.H{
		"uploadUrl": uploadUrl,
		"key":       key,
	})
}

func (r UseCase) VisitFile(c *app.Context) result.Result {
	key := strings.TrimPrefix(c.Param("key"), "/")
	c.Redirect(302, r.Storage.GetFileVisitURL(key))
	return result.Ok(nil)
}

func (r UseCase) GetFileUrl(c *app.Context) result.Result {
	key := strings.TrimPrefix(c.Param("key"), "/")
	return result.Ok(gin.H{"url": r.Storage.GetFileVisitURL(key), "key": key})
}
