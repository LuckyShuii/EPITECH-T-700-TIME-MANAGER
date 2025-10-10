package router

import (
	"app/internal/app/user/handler"
	"app/internal/app/user/repository"
	"app/internal/app/user/service"
	"app/internal/db"

	authHandler "app/internal/app/auth/handler"
	authService "app/internal/app/auth/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	db := db.ConnectPostgres()

	/**
	 * User Module
	 */
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	/**
	 * Auth Module
	 */
	authService := authService.NewAuthService(userService)
	authHandler := authHandler.NewAuthHandler(authService)

	api := r.Group("/api")
	{
		/**
		 * Auth Routes
		 */
		api.POST("/authenticate", authHandler.LoginHandler)
		api.GET("/me", authHandler.MeHandler)

		/**
		 * User Routes
		 */
		api.GET("/users", userHandler.GetUsers)
		api.POST("/users/register", userHandler.RegisterUser)
	}

	return r
}
