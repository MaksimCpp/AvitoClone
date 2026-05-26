// @title AvitoClone API
// @version 1.0
// @description Avito clone API
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"time"

	docs "github.com/MaksimCpp/AvitoClone/docs"
	"github.com/MaksimCpp/AvitoClone/internal/config"
	httpdelivery "github.com/MaksimCpp/AvitoClone/internal/delivery/http"
	"github.com/MaksimCpp/AvitoClone/internal/delivery/http/handler"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/database"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/hash"
	jwtservice "github.com/MaksimCpp/AvitoClone/internal/infrastructure/jwt"
	"github.com/MaksimCpp/AvitoClone/internal/repository/postgresql"
	itemusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/item"
	userusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()
	docs.SwaggerInfo.Host = cfg.SwaggerHost
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

	itemRepo := postgresql.NewPostgreSQLItemRepository(pool)
	createItemUseCase := itemusecase.NewPostgreSQLCreateItemUseCase(itemRepo)
	deleteItemUseCase := itemusecase.NewPostgreSQLDeleteItemUseCasee(itemRepo)
	getItemByIDUseCase := itemusecase.NewPostgreSQLGetItemByIDUseCase(itemRepo)
	listItemsByUserIDUseCase := itemusecase.NewPostgreSQLListItemsByUserIDUseCase(itemRepo)
	listItemsUseCase := itemusecase.NewPostgreSQLListItemsUseCase(itemRepo)
	itemHandler := handler.NewItemHandler(
		createItemUseCase, deleteItemUseCase, getItemByIDUseCase,
		listItemsByUserIDUseCase, listItemsUseCase,
	)

	router := gin.Default()
	err = router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal(err)
	}
	
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
	httpdelivery.SetupRoutes(
		router,
		userHandler,
		itemHandler,
		jwtService,
	)

	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}