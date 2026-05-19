package main

import (
	"log"
	"net/http"

	"github.com/MaksimCpp/AvitoClone/internal/config"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	pool, err := database.NewPostgresPool(cfg.ConnStr())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	err = router.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal(err)
	}
}