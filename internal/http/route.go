package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount"
	"gongfu/internal/config"
	"gongfu/internal/http/action"
	"gongfu/internal/http/admin"
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
	AdminUseCase    admin.UseCase
	Middleware      middlewares.Middlewares
}

func NewRoute(
	config *config.Config,
	officialAccount *officialaccount.OfficialAccount,
	controller *controller.Controller,
	adminUseCase admin.UseCase,
	middleware middlewares.Middlewares,
) *Route {
	return &Route{
		Config:          config,
		OfficialAccount: officialAccount,
		Controller:      controller,
		AdminUseCase:    adminUseCase,
		Middleware:      middleware,
	}
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
	app.GET("wechat-login", action.Wrap(r.Controller.WeChatLogin))
	api := app.Group("/api")
	{
		api.POST("/auth/login", action.Wrap(r.Controller.Login))
		authedApi := api.Group("", r.Middleware.Auth())
		{
			authedApi.POST("me", action.Wrap(r.Controller.EditMe))
			authedApi.GET("me", action.Wrap(r.Controller.Me))
			authedApi.POST("bind/phone", action.Wrap(r.Controller.GetBindCode))
			authedApi.POST("bind/phone/valid", action.Wrap(r.Controller.ValidBindCode))
			authedApi.GET("check-in/today", action.Wrap(r.Controller.GetTodayCheckIn))
			authedApi.POST("check-in", action.Wrap(r.Controller.PostCheckIn))
			// 获取上传文件的 token
			authedApi.GET("storage/upload-token", action.Wrap(r.Controller.GetUploadToken))
			authedApi.GET("coach", action.Wrap(r.Controller.GetCoach))

			adminApi := authedApi.Group("admin", r.Middleware.Role(model.ROLE_ADMIN))
			{
				adminApi.GET("users", action.Wrap(r.AdminUseCase.AdminGetUsers))
				adminApi.GET("user/:id", action.Wrap(r.AdminUseCase.AdminGetUser))
				adminApi.PUT("user/:id", action.Wrap(r.AdminUseCase.AdminUpdateUser))
				adminApi.GET("coach/:id", action.Wrap(r.AdminUseCase.AdminGetCoach))
				adminApi.GET("coaches", action.Wrap(r.AdminUseCase.GetCoaches))
				adminApi.POST("course", action.Wrap(r.AdminUseCase.CreateCourse))
				adminApi.GET("courses", action.Wrap(r.AdminUseCase.GetCoursePage))
				adminApi.GET("courses/:id", action.Wrap(r.AdminUseCase.GetCourse))
				adminApi.GET("role-names", action.Wrap(r.AdminUseCase.AdminGetRoleNames))
			}
		}
		api.GET("wechat/js-config", action.Wrap(r.Controller.JSConfig))
		api.GET("check-in/top", action.Wrap(r.Controller.GetCheckInTop))
		api.GET("check-in/top/count", action.Wrap(r.Controller.GetCheckInCountTop))
		api.GET("check-in/top/continuous", action.Wrap(r.Controller.GetCheckInContinuousTop))
		api.GET("check-in/histories", action.Wrap(r.Controller.GetCheckInHistories))
		api.GET("check-in/:key", action.Wrap(r.Controller.GetCheckIn))

		// 用户详情页
		api.GET("profile/:uuid", action.Wrap(r.Controller.Profile))
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
