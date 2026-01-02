package utils

import (
	"errors"
	"os"
	"time"

	"job-portal-api/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("jwt secret not configured")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
