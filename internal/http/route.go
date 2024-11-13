package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gongfu/internal/client"
	"gongfu/internal/config"
	"gongfu/internal/http/action"
	"gongfu/internal/http/admin"
	"gongfu/internal/http/controller"
	"gongfu/internal/http/middlewares"
	"gongfu/internal/model"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Route struct {
	Config          *config.Config
	OfficialAccount *client.OfficialAccount
	UserUseCase     *controller.UseCase
	AdminUseCase    admin.UseCase
	Middleware      middlewares.Middlewares
}

func NewRoute(
	config *config.Config,
	officialAccount *client.OfficialAccount,
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
		api.GET("storage/visit/*key", action.Wrap(r.UserUseCase.VisitFile))
		authedApi := api.Group("", r.Middleware.Auth())
		{
			authedApi.POST("me", action.Wrap(r.UserUseCase.EditMe))
			authedApi.GET("me", action.Wrap(r.UserUseCase.Me))
			authedApi.POST("bind/phone", action.Wrap(r.UserUseCase.GetBindCode))
			authedApi.POST("bind/phone/valid", action.Wrap(r.UserUseCase.ValidBindCode))
			checkIn := authedApi.Group("check-in")
			{
				// 获取今日打卡信息
				checkIn.GET("today", action.Wrap(r.UserUseCase.GetTodayCheckIn))
				// 创建打卡
				checkIn.POST("", action.Wrap(r.UserUseCase.PostCheckIn))
				// 评论打卡
				checkIn.POST("comment", action.Wrap(r.UserUseCase.CreateCheckInComment))
				// 获取打卡评论列表
				checkIn.GET("comments", action.Wrap(r.UserUseCase.GetCheckInComments))
			}

			// 获取上传文件的 token
			authedApi.GET("storage/upload-token", action.Wrap(r.UserUseCase.GetUploadToken))
			authedApi.GET("coach", action.Wrap(r.UserUseCase.GetCoach))

			// courses
			authedApi.GET("courses", action.Wrap(r.UserUseCase.GetCourses))
			authedApi.GET("courses/:id", action.Wrap(r.UserUseCase.GetCourse))
			authedApi.PUT("courses/:id", action.Wrap(r.UserUseCase.UpdateCourse))

			// update member courses remain
			authedApi.PUT("member-courses/:id/remain", action.Wrap(r.UserUseCase.ChangeMemberCourseRemain))

			adminApi := authedApi.Group("admin", r.Middleware.Role(model.ROLE_ADMIN))
			{
				adminApi.GET("users", action.Wrap(r.AdminUseCase.AdminGetUsers))
				adminApi.GET("user/:id", action.Wrap(r.AdminUseCase.AdminGetUser))
				adminApi.PUT("user/:id", action.Wrap(r.AdminUseCase.AdminUpdateUser))
				adminApi.GET("coach/:id", action.Wrap(r.AdminUseCase.AdminGetCoach))
				adminApi.GET("coaches", action.Wrap(r.AdminUseCase.GetCoaches))
				adminApi.POST("course", action.Wrap(r.AdminUseCase.CreateCourse))
				adminApi.GET("course/list", action.Wrap(r.AdminUseCase.GetCourseList))
				adminApi.GET("schools", action.Wrap(r.AdminUseCase.GetSchools))
				adminApi.GET("courses", action.Wrap(r.AdminUseCase.GetCoursePage))
				adminApi.GET("courses/:id", action.Wrap(r.AdminUseCase.GetCourse))
				adminApi.PUT("courses/:id", action.Wrap(r.AdminUseCase.UpdateCourse))
				adminApi.DELETE("courses/:id", action.Wrap(r.AdminUseCase.DeleteCourse))
				adminApi.GET("role-names", action.Wrap(r.AdminUseCase.AdminGetRoleNames))
				adminApi.POST("teaching-record", action.Wrap(r.AdminUseCase.AdminEditTeachingRecord))
				adminApi.DELETE("teaching-record/:id", action.Wrap(r.AdminUseCase.AdminDeleteTeachingRecord))
				adminApi.POST("study-record", action.Wrap(r.AdminUseCase.AdminEditStudyRecord))
				adminApi.DELETE("study-record/:id", action.Wrap(r.AdminUseCase.AdminDeleteStudyRecord))

				memberCourse := adminApi.Group("member-courses")
				{
					memberCourse.GET("", action.Wrap(r.AdminUseCase.GetMemberCourses))
					memberCourse.POST("", action.Wrap(r.AdminUseCase.CreateMemberCourse))
					memberCourse.PUT("/:id", action.Wrap(r.AdminUseCase.UpdateMemberCourse))
					memberCourse.DELETE("/:id", action.Wrap(r.AdminUseCase.DeleteMemberCourse))
				}
			}
		}
		api.GET("wechat/js-config", action.Wrap(r.UserUseCase.JSConfig))
		api.GET("wechat/pay", r.Middleware.Auth(), action.Wrap(r.UserUseCase.Pay))
		api.POST("wechat/pay-notify", action.Wrap(r.UserUseCase.PayNotify))
		api.GET("wechat/enroll", r.Middleware.Auth(), action.Wrap(r.UserUseCase.GetEnroll))

		checkIn := api.Group("check-in")
		{
			// 打卡排行榜(时间)
			checkIn.GET("top", action.Wrap(r.UserUseCase.GetCheckInTop))
			// 打卡次数排行榜
			checkIn.GET("top/count", action.Wrap(r.UserUseCase.GetCheckInCountTop))
			// 连续打卡排行榜
			checkIn.GET("top/continuous", action.Wrap(r.UserUseCase.GetCheckInContinuousTop))
			// 获取个人的打卡历史
			checkIn.GET("histories", action.Wrap(r.UserUseCase.GetCheckInHistories))
			// 获取打卡详情
			checkIn.GET(":key", action.Wrap(r.UserUseCase.GetCheckIn))
		}

		// 用户详情页
		api.GET("profile/:uuid", action.Wrap(r.UserUseCase.Profile))
	}

	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/me")
	})

	if r.Config.IsProd() {
		// prod 代理到静态文件
		app.Static("/web", r.Config.WebPath)
		app.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.Path, "/web") {
				c.File(r.Config.WebPath + "/index.html")
			} else {
				c.String(http.StatusNotFound, "not found route")
			}
		})
	} else {
		// local 代理到 :4200
		app.GET("/web/*webPath", func(c *gin.Context) {
			uri := "http://127.0.0.1:4200"
			remote, err := url.Parse(uri) // backend server
			if err != nil {
				log.Fatal(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}
}
