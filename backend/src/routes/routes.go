package routes

import (
	"net/http"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/controllers"
	"groundTurn/src/middlewares"
	"groundTurn/src/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Start builds the Gin engine with every entity route behind JWT + RBAC.
func Start(addr string, db *gorm.DB, cfg config.Config) error {
	r := NewRouter(db, cfg)
	return r.Run(addr)
}

// NewRouter wires the full HTTP stack; extracted so integration tests can
// drive the real router over httptest.
func NewRouter(db *gorm.DB, cfg config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.RateLimit(cfg.RateLimitQPS))

	// Services (constructor wiring layer).
	planSvc := services.NewPlanService(db)
	querySvc := services.NewQueryService(db)
	authSvc := services.NewAuthService(db, cfg.JWTSecret)
	taskSvc := services.NewTaskService(db)
	bookingSvc := services.NewBookingService(db)
	resourceSvc := services.NewResourceService(db)
	delaySvc := services.NewDelayService(db, planSvc)

	// Controllers.
	authCtrl := controllers.NewAuthController(authSvc)
	turnCtrl := controllers.NewTurnaroundController(planSvc, querySvc)
	taskCtrl := controllers.NewTaskController(taskSvc, querySvc)
	resourceCtrl := controllers.NewResourceController(resourceSvc, querySvc)
	bookingCtrl := controllers.NewBookingController(bookingSvc, querySvc)
	delayCtrl := controllers.NewDelayController(delaySvc, querySvc)
	dashCtrl := controllers.NewDashboardController(querySvc)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ground-turn"})
	})

	api := r.Group("/api")
	// Public login; all other routes require JWT.
	api.POST("/auth/login", authCtrl.Login)

	authed := api.Group("")
	authed.Use(middlewares.AuthRequired(cfg.JWTSecret))
	authed.Use(middlewares.AuditContext())
	{
		authed.GET("/auth/me", authCtrl.Me)

		// Dashboard: every authenticated role may view.
		authed.GET("/dashboard/stats", dashCtrl.Stats)
		authed.GET("/dashboard/pending-bookings", dashCtrl.PendingBookings)
		authed.GET("/audit-logs", middlewares.RequireRole(
			constants.RoleDispatcher, constants.RoleSupervisor), dashCtrl.AuditLogs)

		// Turnarounds.
		authed.GET("/turnarounds", turnCtrl.List)
		authed.GET("/turnarounds/:id", turnCtrl.Detail)
		authed.POST("/turnarounds",
			middlewares.RequireRole(constants.RoleDispatcher, constants.RoleSupervisor),
			turnCtrl.Register)
		authed.POST("/turnarounds/:id/plan",
			middlewares.RequireRole(constants.RoleDispatcher, constants.RoleSupervisor),
			turnCtrl.GeneratePlan)

		// Ground tasks: teams sign/complete; dispatcher/supervisor can block.
		authed.GET("/ground-tasks", taskCtrl.List)
		authed.POST("/ground-tasks/:id/accept",
			middlewares.RequireRole(constants.RoleDispatcher, constants.RoleTeam, constants.RoleSupervisor),
			taskCtrl.Accept)
		authed.POST("/ground-tasks/:id/complete",
			middlewares.RequireRole(constants.RoleDispatcher, constants.RoleTeam, constants.RoleSupervisor),
			taskCtrl.Complete)
		authed.POST("/ground-tasks/:id/block",
			middlewares.RequireRole(constants.RoleDispatcher, constants.RoleTeam, constants.RoleSupervisor),
			taskCtrl.Block)

		// Resources.
		authed.GET("/ground-resources", resourceCtrl.List)
		authed.GET("/resource-calendar", resourceCtrl.Calendar)
		authed.PATCH("/ground-resources/:id/status",
			middlewares.RequireRole(constants.RoleResourceManager, constants.RoleSupervisor),
			resourceCtrl.UpdateStatus)

		// Bookings: resource managers + supervisors adjust pending windows.
		authed.GET("/resource-bookings", bookingCtrl.List)
		authed.POST("/resource-bookings/:id/adjust",
			middlewares.RequireRole(constants.RoleResourceManager, constants.RoleDispatcher, constants.RoleSupervisor),
			bookingCtrl.Adjust)
		authed.POST("/resource-bookings/:id/release",
			middlewares.RequireRole(constants.RoleResourceManager, constants.RoleDispatcher, constants.RoleSupervisor),
			bookingCtrl.Release)

		// Delays.
		authed.GET("/delay-events", delayCtrl.List)
		authed.POST("/turnarounds/:id/delays",
			middlewares.RequireRole(constants.RoleDispatcher, constants.RoleSupervisor),
			delayCtrl.Register)
		authed.POST("/delay-events/:id/resolve",
			middlewares.RequireRole(constants.RoleDispatcher, constants.RoleSupervisor),
			delayCtrl.Resolve)
	}

	return r
}
