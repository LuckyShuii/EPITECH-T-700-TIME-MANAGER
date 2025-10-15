package router

import (
	"app/docs"
	"app/internal/app/user/handler"
	"app/internal/app/user/repository"
	"app/internal/app/user/service"
	"app/internal/db"
	"net/http"

	authH "app/internal/app/auth/handler"
	authS "app/internal/app/auth/service"
	authM "app/internal/middleware"

	workSessionH "app/internal/app/work-session/handler"
	workSessionR "app/internal/app/work-session/repository"
	workSessionS "app/internal/app/work-session/service"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger setup
	docs.SwaggerInfo.Title = "Time Manager API"
	docs.SwaggerInfo.Description = "API documentation for Time Manager backend"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = "localhost:8081"

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
	// @Summary Login user
	// @Description Authenticate a user with email and password
	// @Tags Auth
	// @Accept json
	// @Produce json
	// @Param credentials body LoginRequest true "User credentials"
	// @Success 200 {object} LoginResponse
	// @Failure 401 {object} ErrorResponse
	// @Router /auth/login [post]
	r.POST("/api/auth/login", authHandler.LoginHandler)

	/**
	 * Protected Routes
	 */
	protected := r.Group("/api")
	// Route par défaut Swagger UI
	protected.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	protected.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Time Manager API Docs</title>
				<!-- ReDoc -->
				<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
			</head>
			<body>
				<redoc spec-url='/api/swagger/doc.json'></redoc>
			</body>
			</html>
		`)
	})

	protected.Use(authMiddleware.AuthenticationMiddleware)
	{

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
		// authMiddleware.RequireRoles("user")
		protected.POST("/work-session/update-clocking", workSessionHandler.UpdateWorkSessionClocking)
	}

	return r
}
