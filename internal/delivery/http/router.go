package httpdelivery

import (
	"github.com/MaksimCpp/AvitoClone/internal/delivery/http/handler"
	"github.com/MaksimCpp/AvitoClone/internal/delivery/http/middleware"
	jwtservice "github.com/MaksimCpp/AvitoClone/internal/infrastructure/jwt"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	jwtService *jwtservice.JWTService,
) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	users := router.Group("/api/v1/users")
	users.Use(middleware.AuthMiddleware(jwtService))
	{
		users.GET("/me", userHandler.GetMe)
	}
}