package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppHost string
	SwaggerHost string

	Host string
	Port string
	User string
	Password string
	DB string

	JWTSecret string
	JWTExpiresHours string

	MinioEndpoint string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket string
	MinioUseSSL bool
}

func Load() *Config {
	_ = godotenv.Load()
	
	cfg := &Config{
		AppPort: os.Getenv("APP_PORT"),
		AppHost: os.Getenv("APP_HOST"),
		SwaggerHost: os.Getenv("SWAGGER_HOST"),
		Host: os.Getenv("POSTGRES_HOST"),
		Port: os.Getenv("POSTGRES_PORT"),
		User: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DB: os.Getenv("POSTGRES_DB"),

		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTExpiresHours: os.Getenv("JWT_EXPIRES_HOURS"),

		MinioEndpoint: os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
		MinioBucket: os.Getenv("MINIO_BUCKET"),
		MinioUseSSL: os.Getenv("MINIO_USE_SSL") == "true",
	}

	return cfg
}

func (c *Config) ConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
	)
}
