package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string) (accessToken, refreshToken string, err error) {
	accessToken = "fake_access_token"
	refreshToken = "fake_refresh_token"
	return
}

func ParseToken(tokenString string) (*Claims, error) {
	return &Claims{}, nil
}
