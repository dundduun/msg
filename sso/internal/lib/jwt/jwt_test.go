package jwt

import (
	"fmt"
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
	t.Run("expired ticket", func(t *testing.T) {
		tokenString := newToken(t, 0)
		_, err := parse(tokenString)
		assert.ErrorIs(t, err, jwt.ErrTokenExpired)
	})

	t.Run("ticket was updated", func(t *testing.T) {
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

func TestRealToken(t *testing.T) { // временное решение из-за отсутствия регистрации
	const (
		secret = "yeasmWPaUgcVjLbdl/R9+CvSSHOOUzIJx/5AQsZnrrA="
		uid    = "01a0963e-1db1-7600-9a6b-59385e9e01ab"
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	fmt.Println("ticket: ", tokenString)
}

// sweep password from this function, he is not used in ticket creation
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
