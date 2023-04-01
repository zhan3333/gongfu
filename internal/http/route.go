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
	UserUseCase     *controller.UseCase
	AdminUseCase    admin.UseCase
	Middleware      middlewares.Middlewares
}

func NewRoute(
	config *config.Config,
	officialAccount *officialaccount.OfficialAccount,
	controller *controller.UseCase,
	adminUseCase admin.UseCase,
	middleware middlewares.Middlewares,
) *Route {
	return &Route{
		Config:          config,
		OfficialAccount: officialAccount,
		UserUseCase:     controller,
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
	app.GET("wechat-login", action.Wrap(r.UserUseCase.WeChatLogin))
	api := app.Group("/api")
	{
		api.POST("/auth/login", action.Wrap(r.UserUseCase.Login))
		authedApi := api.Group("", r.Middleware.Auth())
		{
			authedApi.POST("me", action.Wrap(r.UserUseCase.EditMe))
			authedApi.GET("me", action.Wrap(r.UserUseCase.Me))
			authedApi.POST("bind/phone", action.Wrap(r.UserUseCase.GetBindCode))
			authedApi.POST("bind/phone/valid", action.Wrap(r.UserUseCase.ValidBindCode))
			authedApi.GET("check-in/today", action.Wrap(r.UserUseCase.GetTodayCheckIn))
			authedApi.POST("check-in", action.Wrap(r.UserUseCase.PostCheckIn))
			// 获取上传文件的 token
			authedApi.GET("storage/upload-token", action.Wrap(r.UserUseCase.GetUploadToken))
			authedApi.GET("coach", action.Wrap(r.UserUseCase.GetCoach))

			// courses
			authedApi.GET("courses", action.Wrap(r.UserUseCase.GetCourses))
			authedApi.GET("courses/:id", action.Wrap(r.UserUseCase.GetCourse))
			authedApi.PUT("courses/:id", action.Wrap(r.UserUseCase.UpdateCourse))

			adminApi := authedApi.Group("admin", r.Middleware.Role(model.ROLE_ADMIN))
			{
				adminApi.GET("users", action.Wrap(r.AdminUseCase.AdminGetUsers))
				adminApi.GET("user/:id", action.Wrap(r.AdminUseCase.AdminGetUser))
				adminApi.PUT("user/:id", action.Wrap(r.AdminUseCase.AdminUpdateUser))
				adminApi.GET("coach/:id", action.Wrap(r.AdminUseCase.AdminGetCoach))
				adminApi.GET("coaches", action.Wrap(r.AdminUseCase.GetCoaches))
				adminApi.POST("course", action.Wrap(r.AdminUseCase.CreateCourse))
				adminApi.GET("schools", action.Wrap(r.AdminUseCase.GetSchools))
				adminApi.GET("courses", action.Wrap(r.AdminUseCase.GetCoursePage))
				adminApi.GET("courses/:id", action.Wrap(r.AdminUseCase.GetCourse))
				adminApi.PUT("courses/:id", action.Wrap(r.AdminUseCase.UpdateCourse))
				adminApi.DELETE("courses/:id", action.Wrap(r.AdminUseCase.DeleteCourse))
				adminApi.GET("role-names", action.Wrap(r.AdminUseCase.AdminGetRoleNames))
			}
		}
		api.GET("wechat/js-config", action.Wrap(r.UserUseCase.JSConfig))
		api.GET("check-in/top", action.Wrap(r.UserUseCase.GetCheckInTop))
		api.GET("check-in/top/count", action.Wrap(r.UserUseCase.GetCheckInCountTop))
		api.GET("check-in/top/continuous", action.Wrap(r.UserUseCase.GetCheckInContinuousTop))
		api.GET("check-in/histories", action.Wrap(r.UserUseCase.GetCheckInHistories))
		api.GET("check-in/:key", action.Wrap(r.UserUseCase.GetCheckIn))

		// 用户详情页
		api.GET("profile/:uuid", action.Wrap(r.UserUseCase.Profile))
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
