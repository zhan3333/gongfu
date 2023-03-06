package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount"
	"gongfu/internal/config"
	"gongfu/internal/http/controller"
	"gongfu/internal/http/middlewares"
	"gongfu/internal/model"
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
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
	}))
	app.GET("wechat-login", controller.Wrap(r.Controller.WeChatLogin))
	api := app.Group("/api")
	{
		api.POST("/auth/login", controller.Wrap(r.Controller.Login))
		authedApi := api.Group("", r.Middleware.Auth())
		{
			authedApi.POST("me", controller.Wrap(r.Controller.EditMe))
			authedApi.GET("me", controller.Wrap(r.Controller.Me))
			authedApi.POST("bind/phone", controller.Wrap(r.Controller.GetBindCode))
			authedApi.POST("bind/phone/valid", controller.Wrap(r.Controller.ValidBindCode))
			authedApi.GET("check-in/today", controller.Wrap(r.Controller.GetTodayCheckIn))
			authedApi.POST("check-in", controller.Wrap(r.Controller.PostCheckIn))
			// 获取上传文件的 token
			authedApi.GET("storage/upload-token", controller.Wrap(r.Controller.GetUploadToken))
			authedApi.GET("coach", controller.Wrap(r.Controller.GetCoach))

			adminApi := authedApi.Group("admin", r.Middleware.Role(model.ROLE_ADMIN))
			{
				adminApi.GET("users", controller.Wrap(r.Controller.AdminGetUsers))
				adminApi.GET("user/:id", controller.Wrap(r.Controller.AdminGetUser))
				adminApi.PUT("user/:id", controller.Wrap(r.Controller.AdminUpdateUser))
				adminApi.GET("coach/:id", controller.Wrap(r.Controller.AdminGetCoach))
				adminApi.GET("role-names", controller.Wrap(r.Controller.AdminGetRoleNames))
			}
		}
		api.GET("wechat/js-config", controller.Wrap(r.Controller.JSConfig))
		api.GET("check-in/top", controller.Wrap(r.Controller.GetCheckInTop))
		api.GET("check-in/top/count", controller.Wrap(r.Controller.GetCheckInCountTop))
		api.GET("check-in/top/continuous", controller.Wrap(r.Controller.GetCheckInContinuousTop))
		api.GET("check-in/histories", controller.Wrap(r.Controller.GetCheckInHistories))
		api.GET("check-in/:key", controller.Wrap(r.Controller.GetCheckIn))

		// 用户详情页
		api.GET("profile/:uuid", controller.Wrap(r.Controller.Profile))
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
