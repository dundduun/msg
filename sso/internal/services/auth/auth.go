package auth

import (
	"context"
	"github.com/dundduun/msg/sso/internal/domain/models"
	"go.uber.org/zap"
	"time"
)

type Auth struct {
	logger       *zap.Logger
	userSaver    UserSaver
	userProvider UserProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	ProvideUser(ctx context.Context, email string) (user models.User, err error)
}

func New(
	logger *zap.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		logger:       logger,
		userSaver:    userSaver,
		userProvider: userProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string) (token string, err error) {
	panic("not implemented")
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (uid int64, err error) {
	panic("not implemented")
}
