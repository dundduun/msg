package jwt

import (
	"github.com/dundduun/msg/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateToken(user models.User, duration time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(duration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
