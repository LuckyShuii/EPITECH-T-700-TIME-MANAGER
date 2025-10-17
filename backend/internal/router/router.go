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
	workSessionService := workSessionS.NewWorkSessionService(workSessionRepo, userService)
	workSessionHandler := workSessionH.NewWorkSessionHandler(workSessionService)

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
		protected.GET("/auth/me", authHandler.MeHandler)
		protected.POST("/auth/logout", authHandler.LogoutHandler)

		/**
		 * User Management Routes
		 */
		protected.GET("/users", authMiddleware.RequireRoles("user_manager"), userHandler.GetUsers)
		protected.POST("/users/register", authMiddleware.RequireRoles("admin"), userHandler.RegisterUser)

		/**
		 * Work Sessions Routes
		 */
		// authMiddleware.RequireRoles("user")
		protected.POST("/work-session/update-clocking", workSessionHandler.UpdateWorkSessionClocking)
	}

	return r
}
