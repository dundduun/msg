package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type role int

var (
	ErrInvalidToken = errors.New("invalid ticket")

	validMethods = jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})
)

type Claims struct {
	UID uuid.UUID
	//Role   role
	jwt.RegisteredClaims
}

type Service struct {
	secret []byte
}

func NewService(secret []byte) Service {
	return Service{secret: secret}
}

func (s *Service) verifyKey(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return s.secret, nil
}

func (s *Service) VerifyToken(tokenString string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, s.verifyKey, validMethods)
	if err != nil || !token.Valid {
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}
