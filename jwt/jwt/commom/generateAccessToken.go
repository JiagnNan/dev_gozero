package commom

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func GenerateAccessToken(secret string, userId int64, role string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     now.Add(time.Hour * 24).Unix(),
		"iat":     now.Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func GenerateRefreshToken(secret string, refreshToken string, userId int64) (string, error) {
	now := time.Now()
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("refresh Token 已失效或非法")
	}
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     now.Add(time.Hour * 24).Unix(),
		"iat":     now.Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func CheckPassword(password, OldPassword string) bool {
	if password == OldPassword {
		return true
	}
	return false
}
