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

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	database := db.ConnectPostgres()
	db.ConnectRedis()

	/**
	 * User Module
	 */
	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	/**
	 * Auth Module
	 */
	authService := authS.NewAuthService(userService)
	authHandler := authH.NewAuthHandler(authService)
	authMiddleware := &authM.AuthHandler{Service: authService}

	/**
	* Work Sessions Routes
	 */
	workSessionRepo := workSessionR.NewWorkSessionRepository(database)
	workSessionService := workSessionS.NewWorkSessionService(workSessionRepo, userService, BreakR.NewBreakRepository(database))
	workSessionHandler := workSessionH.NewWorkSessionHandler(workSessionService)

	breakRepo := BreakR.NewBreakRepository(database)
	breakService := BreakS.NewBreakService(breakRepo, workSessionRepo)
	breakHandler := BreakH.NewBreakHandler(breakService)

	teamRepo := TeamR.NewTeamRepository(database)
	teamService := TeamS.NewTeamService(teamRepo, userService)
	teamHandler := TeamH.NewTeamHandler(teamService, userService)

	weeklyRateRepo := WeeklyRatesR.NewWeeklyRateRepository(database)
	weeklyRateService := WeeklyRatesS.NewWeeklyRateService(weeklyRateRepo)
	weeklyRateHandler := WeeklyRatesH.NewWeeklyRateHandler(weeklyRateService)

	/**
	* Public Routes
	 */
	r.POST("/api/auth/login", authHandler.LoginHandler)

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
		// TODO: on user creation add the weekly rate in the form to check
		protected.POST("/users/register", authMiddleware.RequireRoles("admin"), userHandler.RegisterUser)

		protected.PUT("/users/update-status", authMiddleware.RequireRoles("admin"), userHandler.UpdateUserStatus)
		// TODO: on update a user add the weekly rate in the form to check
		protected.PUT("/users", authMiddleware.RequireRoles("admin"), userHandler.UpdateUser)

		protected.GET("/users", authMiddleware.RequireRoles("admin"), userHandler.GetUsers)
		protected.GET("/users/:uuid", authMiddleware.RequireRoles("all"), userHandler.GetUserByUUID)
		protected.GET("/users/weekly-rates", authMiddleware.RequireRoles("all"), weeklyRateHandler.GetAll)
		protected.POST("/users/weekly-rates/create", authMiddleware.RequireRoles("admin"), weeklyRateHandler.Create)
		protected.PUT("/users/weekly-rates/:uuid/update", authMiddleware.RequireRoles("admin"), weeklyRateHandler.Update)
		// TODO: protected.DELETE("/users/weekly-rates/:uuid/delete
		// TODO: protected.POST("/users/weekly-rates/:weekly_rate_uuid/assign-to-user/:user_uuid

		protected.DELETE("/users/delete", authMiddleware.RequireRoles("admin"), userHandler.DeleteUser)

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
	}

	return r
}
