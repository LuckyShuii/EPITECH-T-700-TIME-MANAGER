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

	// ✅ éviter le conflit de nom
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
	 * Public Routes
	 */
	r.POST("/api/authenticate", authHandler.LoginHandler)

	// ✅ Protected Routes
	protected := r.Group("/api")
	protected.Use(authMiddleware.AuthenticationMiddleware)
	{
		protected.POST("/logout", authHandler.LogoutHandler)
		protected.GET("/me", authHandler.MeHandler)
		protected.GET("/users", authMiddleware.RequireRoles("user_manager"), userHandler.GetUsers)
		protected.POST("/users/register", authMiddleware.RequireRoles("user_manager"), userHandler.RegisterUser)
	}

	return r
}
