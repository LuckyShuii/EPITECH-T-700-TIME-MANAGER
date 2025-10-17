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

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	database := db.ConnectPostgres()

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
		protected.POST("/users/register", authMiddleware.RequireRoles("admin"), userHandler.RegisterUser)

		protected.GET("/users", authMiddleware.RequireRoles("admin"), userHandler.GetUsers)

		protected.DELETE("/users/delete", authMiddleware.RequireRoles("admin"), userHandler.DeleteUser)

		/**
		 * Work Sessions & Breaks Routes
		 */
		protected.POST("/work-session/update-clocking", authMiddleware.RequireRoles("all"), workSessionHandler.UpdateWorkSessionClocking)
		protected.POST("/work-session/update-breaking", authMiddleware.RequireRoles("all"), breakHandler.UpdateBreak)

		protected.GET("/work-session/history", authMiddleware.RequireRoles("all"), workSessionHandler.GetWorkSessionHistory)

		protected.GET("/work-session/status", authMiddleware.RequireRoles("all"), workSessionHandler.GetWorkSessionStatus)
	}

	return r
}
