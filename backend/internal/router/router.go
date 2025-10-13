package router

import (
	"app/internal/app/user/handler"
	"app/internal/app/user/repository"
	"app/internal/app/user/service"
	"app/internal/db"

	authH "app/internal/app/auth/handler"
	authS "app/internal/app/auth/service"
	authM "app/internal/middleware"

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
	 * Protected Routes
	 */
	protected := r.Group("/api")
	protected.Use(authMiddleware.AuthenticationMiddleware)
	{
		/**
		 * Public Routes
		 */
		r.POST("/api/auth/login", authHandler.LoginHandler)

		protected.GET("/auth/me", authHandler.MeHandler)
		protected.POST("/auth/logout", authHandler.LogoutHandler)

		/**
		 * User Management Routes
		 */
		protected.GET("/users", authMiddleware.RequireRoles("user_manager"), userHandler.GetUsers)
		protected.POST("/users/register", authMiddleware.RequireRoles("user_manager"), userHandler.RegisterUser)

		/**
		 * Work Sessions Routes
		 */

	}

	return r
}
