package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/sirupsen/logrus"
	"gongfu/internal/app"
	"gongfu/internal/config"
	"gongfu/internal/model"
	"gongfu/internal/result"
	"gongfu/internal/service"
	"gongfu/internal/service/store"
	"net/http"
	"strings"
)

type Controller struct {
	Config          *config.Config
	OfficialAccount *officialaccount.OfficialAccount
	AuthCode        service.AuthCode
	Token           service.Token
	Store           store.Store
	Storage         service.Storage
}

func NewController(
	config *config.Config,
	officialAccount *officialaccount.OfficialAccount,
	authCode service.AuthCode,
	token service.Token,
	store store.Store,
	storage service.Storage,
) *Controller {
	return &Controller{Config: config, OfficialAccount: officialAccount, AuthCode: authCode, Token: token, Store: store, Storage: storage}
}

type Action func(c *app.Context) result.Result

func Wrap(f Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		c2 := &app.Context{
			Context:   c,
			Return:    app.Return{Context: c},
			User:      model.User{},
			UserID:    c.GetUint("user_id"),
			RequestID: c.Request.Header.Get("x-request-id"),
		}
		if c2.RequestID == "" {
			c2.RequestID = strings.ReplaceAll(uuid.New().String(), "-", "")
		}
		c.Header("x-request-id", c2.RequestID)
		defer func() {
			log := c2.Logger.WithField("code", c.Writer.Status())
			if c.Writer.Status() > 499 {
				log.Error("access_log")
			} else {
				log.Info("access_log")
			}
		}()
		if c2.UserID != 0 {
			if v, exist := c.Get("user"); exist {
				if user, ok := v.(model.User); ok {
					c2.User = user
				}
			}
		}
		c2.Logger = logrus.WithContext(c2).WithField("x-request-id", c2.RequestID).WithField("user_id", c2.UserID)
		r := f(c2)
		if c.Writer.Written() { // already write response
			return
		}
		if err := r.Err(); err != nil {
			c2.Logger = c2.Logger.WithField("error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}
		if r.Ok() == nil {
			c.JSON(http.StatusOK, gin.H{"msg": "ok"})
		} else {
			c.JSON(http.StatusOK, r.Ok())
		}
		return
	}
}
