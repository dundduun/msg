package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dundduun/msg/sso/internal/domain/models"
	"github.com/dundduun/msg/sso/internal/lib/werr"
	"github.com/dundduun/msg/sso/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const cost = 12

var ErrInvalidCredentials = errors.New("invalid credentials")

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
	User(ctx context.Context, email string) (user models.User, err error)
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

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "auth.Login"

	logger := a.logger.With(
		zap.String("op", op),
		zap.String("email", email),
	)
	logger.Info("attempting to log in user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			logger.Warn("user not found", zap.Error(err))
			return "", ErrInvalidCredentials
		}

		logger.Error("failed to provide user", zap.Error(err))
		return "", werr.WrapError(op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		logger.Info("invalid credentials", zap.Error(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	panic("not implemented")
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	logger := a.logger.With(
		zap.String("op", op),
		zap.String("email", email),
	)
	logger.Info("registering user")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		logger.Error("failed to generate hash", zap.Error(err))
		return 0, werr.WrapError(op, err)
	}

	uid, err := a.userSaver.SaveUser(ctx, email, hash)
	if err != nil {
		if errors.Is(err, storage.ErrEmailTaken) {
			logger.Warn("email taken", zap.Error(err))
			return 0, werr.WrapError(op, err)
		}
		logger.Error("failed to save user", zap.Error(err))
		return 0, werr.WrapError(op, err)
	}

	logger.Info("user registered", zap.Int64("uid", uid))

	return uid, nil
}
