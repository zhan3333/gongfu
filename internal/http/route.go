package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount"
	"gongfu/internal/config"
	"gongfu/internal/http/controller"
	"gongfu/internal/http/middlewares"
	"net/http"
	"strings"
	"time"
)

type Route struct {
	Config          *config.Config
	OfficialAccount *officialaccount.OfficialAccount
	Controller      *controller.Controller
	Middleware      middlewares.Middlewares
}

func NewRoute(
	config *config.Config,
	officialAccount *officialaccount.OfficialAccount,
	controller *controller.Controller,
	middleware middlewares.Middlewares,
) *Route {
	return &Route{Config: config, OfficialAccount: officialAccount, Controller: controller, Middleware: middleware}
}

func (r Route) Route(app *gin.Engine) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"X-Request-Id"},
		MaxAge:           12 * time.Hour,
	}))
	app.GET("login", controller.Wrap(r.Controller.Login))
	api := app.Group("/api")
	{
		authedApi := api.Group("", r.Middleware.Auth())
		{
			authedApi.GET("me", controller.Wrap(r.Controller.Me))
			authedApi.POST("bind/phone", controller.Wrap(r.Controller.GetBindCode))
			authedApi.POST("bind/phone/valid", controller.Wrap(r.Controller.ValidBindCode))
			authedApi.GET("check-in/today", controller.Wrap(r.Controller.GetTodayCheckIn))
			authedApi.POST("check-in", controller.Wrap(r.Controller.PostCheckIn))
			// 获取上传文件的 token
			authedApi.GET("storage/upload-token", controller.Wrap(r.Controller.GetUploadToken))
		}
		api.GET("wechat/js-config", controller.Wrap(r.Controller.JSConfig))
		api.GET("check-in/top", controller.Wrap(r.Controller.GetCheckInTop))
		api.GET("check-in/top/count", controller.Wrap(r.Controller.GetCheckInCountTop))
		api.GET("check-in/top/continuous", controller.Wrap(r.Controller.GetCheckInContinuousTop))
		api.GET("check-in/histories", controller.Wrap(r.Controller.GetCheckInHistories))
		api.GET("check-in/:key", controller.Wrap(r.Controller.GetCheckIn))
	}

	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/me")
	})
	// 本地文件存储
	app.Static("/storage", "storage")
	// web 文件
	app.Static("/web", r.Config.WebPath)
	app.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/web") {
			c.File(r.Config.WebPath + "/index.html")
		} else {
			c.String(http.StatusNotFound, "not found route")
		}
	})
}
