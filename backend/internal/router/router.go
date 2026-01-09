package router

import (
	"app/internal/app/user/handler"
	"app/internal/app/user/repository"
	"app/internal/app/user/service"
	"app/internal/db"

	authH "app/internal/app/auth/handler"
	authS "app/internal/app/auth/service"
	authM "app/internal/middleware"

	workSessionH "app/internal/app/work-session/handler"
	workSessionR "app/internal/app/work-session/repository"
	workSessionS "app/internal/app/work-session/service"

	BreakH "app/internal/app/break/handler"
	BreakR "app/internal/app/break/repository"
	BreakS "app/internal/app/break/service"

	TeamH "app/internal/app/team/handler"
	TeamR "app/internal/app/team/repository"
	TeamS "app/internal/app/team/service"

	WeeklyRatesH "app/internal/app/weekly-rate/handler"
	WeeklyRatesR "app/internal/app/weekly-rate/repository"
	WeeklyRatesS "app/internal/app/weekly-rate/service"

	KPIH "app/internal/app/kpi/handler"
	KPIR "app/internal/app/kpi/repository"
	KPIService "app/internal/app/kpi/service"

	"app/internal/app/mailer"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	database := db.ConnectPostgres()
	db.ConnectRedis()

	// 1) Repos
	userRepo := repository.NewUserRepository(database)
	workSessionRepo := workSessionR.NewWorkSessionRepository(database)
	breakRepo := BreakR.NewBreakRepository(database)
	kpiRepo := KPIR.NewKPIRepository(database)
	teamRepo := TeamR.NewTeamRepository(database)
	weeklyRateRepo := WeeklyRatesR.NewWeeklyRateRepository(database)

	// 2) Services
	userService := service.NewUserService(userRepo, mailer.Service)
	weeklyRateService := WeeklyRatesS.NewWeeklyRateService(weeklyRateRepo, userService)

	userService.SetWeeklyRateService(weeklyRateService)

	workSessionService := workSessionS.NewWorkSessionService(workSessionRepo, userService, breakRepo)
	breakService := BreakS.NewBreakService(breakRepo, workSessionRepo)
	teamService := TeamS.NewTeamService(teamRepo, userService)
	kpiService := KPIService.NewKPIService(breakService, teamService, userService, weeklyRateService, kpiRepo)
	authService := authS.NewAuthService(userService)

	// 3) Handlers
	userHandler := handler.NewUserHandler(userService)
	weeklyRateHandler := WeeklyRatesH.NewWeeklyRateHandler(weeklyRateService)
	workSessionHandler := workSessionH.NewWorkSessionHandler(workSessionService)
	breakHandler := BreakH.NewBreakHandler(breakService)
	teamHandler := TeamH.NewTeamHandler(teamService, userService)
	kpiHandler := KPIH.NewKPIHandler(kpiService)
	authHandler := authH.NewAuthHandler(authService)
	authMiddleware := &authM.AuthHandler{Service: authService}

	/**
	* Public Routes
	 */
	r.POST("/api/auth/login", authHandler.LoginHandler)
	r.POST("/api/users/reset-password", userHandler.ResetPassword)
	r.POST("/api/users/update-password", userHandler.UpdateCurrentUserPassword)

	/**
	 * Protected Routes
	 */
	protected := r.Group("/api")
	protected.Use(authMiddleware.AuthenticationMiddleware)
	{
		protected.POST("/auth/logout", authHandler.LogoutHandler)

		protected.GET("/auth/me", authHandler.MeHandler)

		/**
		 * User Management Routes
		 */
		protected.POST("/users/register", authMiddleware.RequireRoles("admin"), userHandler.RegisterUser)
		protected.POST("/users/weekly-rates/create", authMiddleware.RequireRoles("admin"), weeklyRateHandler.Create)
		protected.POST("/users/weekly-rates/:weekly_rate_uuid/assign-to-user/:user_uuid", authMiddleware.RequireRoles("admin"), weeklyRateHandler.AssignToUser)

		protected.PUT("/users/update-status", authMiddleware.RequireRoles("admin"), userHandler.UpdateUserStatus)
		protected.PUT("/users", authMiddleware.RequireRoles("admin"), userHandler.UpdateUser)
		protected.PUT("/users/weekly-rates/:uuid/update", authMiddleware.RequireRoles("admin"), weeklyRateHandler.Update)
		protected.PUT("/users/current-user-dashboard-layout/edit", authMiddleware.RequireRoles("all"), userHandler.UpdateCurrentUserDashboardLayout)

		protected.GET("/users", authMiddleware.RequireRoles("admin"), userHandler.GetUsers)
		protected.GET("/users/:uuid", authMiddleware.RequireRoles("all"), userHandler.GetUserByUUID)
		protected.GET("/users/weekly-rates", authMiddleware.RequireRoles("all"), weeklyRateHandler.GetAll)
		protected.GET("/users/current-user-dashboard-layout", authMiddleware.RequireRoles("all"), userHandler.GetCurrentUserDashboardLayout)

		protected.DELETE("/users/weekly-rates/:uuid/delete", authMiddleware.RequireRoles("admin"), weeklyRateHandler.Delete)
		protected.DELETE("/users/delete/:uuid", authMiddleware.RequireRoles("admin"), userHandler.DeleteUser)
		protected.DELETE("/users/current-user-dashboard-layout/delete", authMiddleware.RequireRoles("all"), userHandler.DeleteCurrentUserDashboardLayout)

		/**
		 * Work Sessions & Breaks Routes
		 */
		protected.POST("/work-session/update-clocking", authMiddleware.RequireRoles("all"), workSessionHandler.UpdateWorkSessionClocking)
		protected.POST("/work-session/update-breaking", authMiddleware.RequireRoles("all"), breakHandler.UpdateBreak)

		protected.GET("/work-session/history", authMiddleware.RequireRoles("all"), workSessionHandler.GetWorkSessionHistory)

		protected.GET("/work-session/status", authMiddleware.RequireRoles("all"), workSessionHandler.GetWorkSessionStatus)

		/**
		 * Teams Routes
		 */
		protected.GET("/teams", authMiddleware.RequireRoles("admin"), teamHandler.GetTeams)
		protected.GET("/teams/:uuid", authMiddleware.RequireRoles("all"), teamHandler.GetTeamByUUID)

		protected.DELETE("/teams/:uuid", authMiddleware.RequireRoles("admin"), teamHandler.DeleteTeamByUUID)
		protected.DELETE("/teams/users/:team_uuid/:user_uuid", authMiddleware.RequireRoles("admin"), teamHandler.RemoveUserFromTeam)

		protected.POST("/teams", authMiddleware.RequireRoles("admin"), teamHandler.CreateTeam)
		protected.POST("/teams/add-users", authMiddleware.RequireRoles("admin"), teamHandler.AddUsersToTeam)

		protected.PUT("/teams/edit/:uuid", authMiddleware.RequireRoles("admin"), teamHandler.UpdateTeamByUUID)
		protected.PUT("/teams/:team_uuid/users/:user_uuid/edit-manager-status/:is_manager", authMiddleware.RequireRoles("admin"), teamHandler.UpdateTeamUserManagerStatus)

		/**
		 * KPI Routes
		 */
		// TODO: route to get all the weekly dates from a user where he worked for the frontend filter/dropdown

		protected.GET("/kpi/work-session-user-weekly-total/:user_uuid/:start_date/:end_date", authMiddleware.RequireRoles("all"), kpiHandler.GetWorkSessionUserWeeklyTotal)
		protected.GET("/kpi/work-session-team-weekly-total/:team_uuid/:start_date/:end_date", authMiddleware.RequireRoles("manager"), kpiHandler.GetWorkSessionTeamWeeklyTotal)
		protected.GET("/kpi/presence-rate/:user_uuid/:start_date/:end_date", authMiddleware.RequireRoles("manager, admin"), kpiHandler.GetPresenceRate)
		protected.GET("/kpi/weekly-average-break-time/:user_uuid/:start_date/:end_date", authMiddleware.RequireRoles("manager, admin"), kpiHandler.GetAverageBreakTime)
		// moyenne par shift par individu
		protected.GET("/kpi/average-time-per-shift/:user_uuid/:start_date/:end_date", authMiddleware.RequireRoles("manager, admin"), kpiHandler.GetAverageTimePerShift)

		protected.POST("/kpi/export", authMiddleware.RequireRoles("manager, admin"), kpiHandler.ExportKPIData)
		protected.GET("/kpi/files/:filename", authMiddleware.RequireRoles("all"), kpiHandler.DownloadKPIFile)
	}

	return r
}
