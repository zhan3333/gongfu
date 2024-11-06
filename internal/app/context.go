package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gongfu/internal/model"
)

type Context struct {
	*gin.Context
	Return
	User      model.User
	UserID    uint
	Logger    *logrus.Entry
	RequestID string
}

type Return struct {
	Context *gin.Context
}

func (r Return) Message(code int, msg string) {
	r.Context.JSON(code, gin.H{"msg": msg})
}

func (c Context) Ctx() context.Context {
	return c.Context.Request.Context()
}
