package profile

import (
	"context"
	"errors"
	"fmt"
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/google/uuid"
	"log/slog"
)

var (
	ErrNoProfile     = errors.New("profile not found")
	ErrUsernameTaken = errors.New("username is taken")
	ErrCacheMiss     = errors.New("profile not found in cache")
	ErrBadId         = errors.New("bad id")
)

type Repo interface {
	ExtractProfile(ctx context.Context, id uuid.UUID) (Profile, error)
	InsertProfile(ctx context.Context, prof Profile) error
}

type Cache interface {
	Get(ctx context.Context, id uuid.UUID) (Profile, error)
	Set(ctx context.Context, profile Profile) error
}

type Service struct {
	cache Cache
	repo  Repo
	log   *slog.Logger
}

func NewService(cache Cache, repo Repo, log *slog.Logger) *Service {
	return &Service{
		cache: cache,
		repo:  repo,
		log:   log,
	}
}

func (s *Service) GetProfile(ctx context.Context, uid uuid.UUID) (Profile, error) {
	const op = "profile.Service.Check"
	log := s.log.With(slog.String("op", op), slog.String("uid", uid.String()))

	if uid == uuid.Nil {
		return Profile{}, ErrBadId
	}

	profile, err := s.cache.Get(ctx, uid)
	if err == nil {
		log.Info("cache hit!")
		return profile, nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		log.Error("failed to get profile from cache", logerr.Err(err))
	}

	log.Info("cache miss")
	profile, err = s.repo.ExtractProfile(ctx, uid)
	if err != nil {
		if !errors.Is(err, ErrNoProfile) {
			log.Error("failed to get profile", logerr.Err(err))
		}
		return Profile{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.cache.Set(ctx, profile)
	if err != nil {
		log.Error("failed to set profile to cache", logerr.Err(err))
	} else {
		log.Info("wrote to cache")
	}

	return profile, nil
}

func (s *Service) CreateProfile(ctx context.Context, username string, name string) error {
	const op = "profile.Service.CreateProfile"
	uid, _ := uuid.NewV7()
	profile := Profile{
		ID:       uid,
		Username: username,
		Name:     name,
	}
	log := s.log.With(slog.String("op", op), slog.String("uid", uid.String()))

	err := s.repo.InsertProfile(ctx, profile)
	if err != nil {
		if errors.Is(err, ErrUsernameTaken) {
			return err
		}

		log.Error("failed to set profile", logerr.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
