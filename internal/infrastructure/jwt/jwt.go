package jwtservice

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret string
	ttl time.Duration
}

type Claims struct {
	UserID int
	jwt.RegisteredClaims
}

func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{
		secret: secret,
		ttl: ttl,
	}
}

func (j *JWTService) Generate(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(j.ttl),
			),
		},
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTService) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	
	return claims, nil
}
