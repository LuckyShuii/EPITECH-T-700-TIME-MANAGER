package router

import (
	"app/internal/app/user/handler"
	"app/internal/app/user/repository"
	"app/internal/app/user/service"
	"app/internal/db"

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

	api := r.Group("/api")
	{
		/**
		 * User Routes
		 */
		api.GET("/users", userHandler.GetUsers)
	}

	return r
}
