package jwt

import (
	"github.com/dundduun/msg/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

const (
	secret = "123"
	email  = "a@bbc.c"
	uid    = int64(12348)
)

func TestCreateToken(t *testing.T) {
	t.Run("expired token", func(t *testing.T) {
		tokenString := newToken(t, 0)
		_, err := parse(tokenString)
		assert.ErrorIs(t, err, jwt.ErrTokenExpired)
	})

	t.Run("token was updated", func(t *testing.T) {
		tokenString := newToken(t, time.Minute)
		tokenString += "ab"
		_, err := parse(tokenString)
		assert.ErrorIs(t, err, jwt.ErrTokenMalformed)
	})

	t.Run("happy", func(t *testing.T) {
		tokenString := newToken(t, time.Hour)
		token, err := parse(tokenString)
		require.NoError(t, err)
		require.True(t, token.Valid)

		t.Run("valid uid and email", func(t *testing.T) {
			claims, ok := token.Claims.(jwt.MapClaims)
			require.True(t, ok)
			assert.Equal(t, email, claims["email"])
			assert.Equal(t, float64(uid), claims["uid"])
		})
	})
}

func newToken(t *testing.T, duration time.Duration) string {
	hash, err := bcrypt.GenerateFromPassword([]byte("abc123"), bcrypt.MinCost)
	require.NoError(t, err)
	tokenString, err := CreateToken(models.User{
		ID:       uid,
		Email:    email,
		PassHash: hash,
	}, duration, secret)
	require.NoError(t, err)

	return tokenString
}

func parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
