// @title AvitoClone API
// @version 1.0
// @description Avito clone API
// @host localhost:8000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"time"

	_ "github.com/MaksimCpp/AvitoClone/docs"
	"github.com/MaksimCpp/AvitoClone/internal/config"
	httpdelivery "github.com/MaksimCpp/AvitoClone/internal/delivery/http"
	"github.com/MaksimCpp/AvitoClone/internal/delivery/http/handler"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/database"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/hash"
	jwtservice "github.com/MaksimCpp/AvitoClone/internal/infrastructure/jwt"
	"github.com/MaksimCpp/AvitoClone/internal/repository/postgresql"
	userusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()
	pool, err := database.NewPostgresPool(cfg.ConnStr())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	hasher := hash.NewBcryptHasher()
	jwtService := jwtservice.NewJWTService(
		cfg.JWTSecret,
		24 * time.Hour,
	)

	userRepo := postgresql.NewPostgreSQLUserRepository(pool)
	registerUserUseCase := userusecase.NewPostgreSQLRegisterUserUseCase(userRepo, hasher)
	loginUserUseCase := userusecase.NewPostgreSQLLoginUserUseCase(userRepo, hasher, jwtService)
	getMeUseCase := userusecase.NewPostgreSQLGetMeUseCase(userRepo)

	userHandler := handler.NewUserHandler(registerUserUseCase, loginUserUseCase, getMeUseCase)

	router := gin.Default()
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
	httpdelivery.SetupRoutes(
		router,
		userHandler,
		jwtService,
	)

	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}