package middleware

import (
	"net/http"
	"strings"

	jwtservice "github.com/MaksimCpp/AvitoClone/internal/infrastructure/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(
	jwtService *jwtservice.JWTService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": "Unauthorized.",
			})

			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": "Invalid authorization header.",
			})

			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtService.Parse(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"detail": err.Error(),
			})

			c.Abort()
			return
		}
	
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
