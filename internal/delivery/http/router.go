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
	itemHandler *handler.ItemHandler,
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
		users.GET("/:user_id/items", itemHandler.ListByUserID)
	}

	items := router.Group("/api/v1/items")
	{
		items.GET("/:id", itemHandler.GetByID)
		items.GET("", itemHandler.List)
	}
	itemsProtected := items.Group("")
	itemsProtected.Use(middleware.AuthMiddleware(jwtService))
	{
		itemsProtected.POST("", itemHandler.Create)
		itemsProtected.DELETE("/:id", itemHandler.Delete)
		itemsProtected.GET("/me", itemHandler.ListMyItems)
	}
}