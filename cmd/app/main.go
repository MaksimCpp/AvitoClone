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
	"strconv"
	"time"

	docs "github.com/MaksimCpp/AvitoClone/docs"
	"github.com/MaksimCpp/AvitoClone/internal/config"
	httpdelivery "github.com/MaksimCpp/AvitoClone/internal/delivery/http"
	"github.com/MaksimCpp/AvitoClone/internal/delivery/http/handler"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/database"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/hash"
	jwtservice "github.com/MaksimCpp/AvitoClone/internal/infrastructure/jwt"
	miniostorage "github.com/MaksimCpp/AvitoClone/internal/infrastructure/storage"
	"github.com/MaksimCpp/AvitoClone/internal/repository/cached"
	"github.com/MaksimCpp/AvitoClone/internal/repository/postgresql"
	avitoredis "github.com/MaksimCpp/AvitoClone/internal/repository/redis"
	itemusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/item"
	itemimageusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/item_image"
	userusecase "github.com/MaksimCpp/AvitoClone/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

	redisDB, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(
		&redis.Options{
			Addr: cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB: redisDB,
		},
	)
	defer rdb.Close()

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

	cache := avitoredis.NewCache(rdb)

	itemRepo := postgresql.NewPostgreSQLItemRepository(pool)
	cachedItemRepo := cached.NewCachedItemRepository(itemRepo, cache)
	createItemUseCase := itemusecase.NewPostgreSQLCreateItemUseCase(cachedItemRepo)
	deleteItemUseCase := itemusecase.NewPostgreSQLDeleteItemUseCasee(cachedItemRepo)
	getItemByIDUseCase := itemusecase.NewPostgreSQLGetItemByIDUseCase(cachedItemRepo)
	listItemsByUserIDUseCase := itemusecase.NewPostgreSQLListItemsByUserIDUseCase(cachedItemRepo)
	listItemsUseCase := itemusecase.NewPostgreSQLListItemsUseCase(cachedItemRepo)

	imageRepo := postgresql.NewPostgreSQLItemImageRepository(pool)
	imageStorage, err := miniostorage.NewMinIOStorage(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioAccessKey,
		cfg.MinioBucket,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Fatal(err)
	}

	uploadImageUseCase := itemimageusecase.NewPostgreSQLUploadImageUseCase(
		imageRepo, itemRepo, imageStorage, cfg,
	)
	deleteImageUseCase := itemimageusecase.NewPostgreSQLDeleteImageUseCase(
		imageRepo, itemRepo, imageStorage,
	)
	listImagerByItemIDUseCase := itemimageusecase.NewPostgreSQLListImagesByItemIDUseCase(
		imageRepo,
	)

	itemHandler := handler.NewItemHandler(
		createItemUseCase, deleteItemUseCase, getItemByIDUseCase,
		listItemsByUserIDUseCase, listItemsUseCase,
		uploadImageUseCase, deleteImageUseCase, listImagerByItemIDUseCase,
		cfg,
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