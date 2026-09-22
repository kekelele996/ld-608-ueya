package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/controllers"
	"groundTurn/src/middlewares"
	"groundTurn/src/services"
)

// Start builds the full router: auth + RBAC + audit + rate limiting, with
// every entity in its own routes section.
func Start(addr string, cfg *config.Config, svc *services.Services) error {
	if cfg.SQLitePath == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(middlewares.ErrorHandlerMiddleware())
	r.Use(middlewares.RateLimitMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ground-turn"})
	})

	auth := controllers.NewAuthController(svc)
	turns := controllers.NewTurnaroundController(svc)
	tasks := controllers.NewTaskController(svc)
	resources := controllers.NewResourceController(svc)
	bookings := controllers.NewBookingController(svc)
	delays := controllers.NewDelayController(svc)
	reports := controllers.NewReportController(svc)

	api := r.Group("/api")

	// Public login; everything else requires a valid JWT.
	api.POST("/auth/login", auth.Login)

	authed := api.Group("")
	authed.Use(middlewares.AuthMiddleware(cfg))
	authed.Use(middlewares.AuditLogMiddleware(svc.DB))
	authed.GET("/auth/me", auth.Me)

	// 运行看板 / 操作日志（所有角色只读）
	authed.GET("/reports/dashboard", reports.Dashboard)
	authed.GET("/reports/audit-logs", reports.AuditLogs)

	// 航班过站
	authed.GET("/turnarounds", turns.List)
	authed.GET("/turnarounds/:id", turns.Detail)
	authed.POST("/turnarounds",
		middlewares.RBACMiddleware(constants.ActionTurnaroundCreate), turns.Create)
	authed.POST("/turnarounds/:id/arrive",
		middlewares.RBACMiddleware(constants.ActionTurnaroundArrive), turns.Arrive)
	authed.POST("/turnarounds/:id/generate-plan",
		middlewares.RBACMiddleware(constants.ActionPlanGenerate), turns.GeneratePlan)
	authed.POST("/turnarounds/:id/release",
		middlewares.RBACMiddleware(constants.ActionTurnaroundRelease), turns.Release)

	// 地勤任务
	authed.GET("/tasks", tasks.List)
	authed.POST("/tasks/:id/sign",
		middlewares.RBACMiddleware(constants.ActionTaskSign), tasks.Sign)
	authed.POST("/tasks/:id/finish",
		middlewares.RBACMiddleware(constants.ActionTaskFinish), tasks.Finish)
	authed.POST("/tasks/:id/block",
		middlewares.RBACMiddleware(constants.ActionTaskBlock), tasks.Block)

	// 保障资源
	authed.GET("/resources", resources.List)
	authed.GET("/resource-bookings", resources.Bookings)
	authed.GET("/resource-bookings/pending", resources.PendingBookings)
	authed.POST("/resources/:id/status",
		middlewares.RBACMiddleware(constants.ActionResourceMaintenance), resources.SetMaintenance)

	// 资源预约调整/释放
	authed.POST("/resource-bookings/:id/adjust",
		middlewares.RBACMiddleware(constants.ActionBookingAdjust), bookings.Adjust)
	authed.POST("/resource-bookings/:id/release",
		middlewares.RBACMiddleware(constants.ActionBookingRelease), bookings.Release)

	// 延误登记 -> 触发整批重排
	authed.GET("/delay-events", delays.List)
	authed.GET("/delay-events/impacts", delays.Impacts)
	authed.POST("/turnarounds/:id/delays",
		middlewares.RBACMiddleware(constants.ActionDelayRegister), delays.Register)
	authed.POST("/delay-events/:id/resolve",
		middlewares.RBACMiddleware(constants.ActionDelayResolve), delays.Resolve)

	return r.Run(addr)
}
