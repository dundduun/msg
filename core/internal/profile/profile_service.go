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
)

type Repo interface {
	ExtractProfile(ctx context.Context, id uuid.UUID) (Profile, error)
	InsertProfile(ctx context.Context, prof Profile) error
}

type Service struct {
	cache *Cache
	repo  Repo
	log   *slog.Logger
}

func NewService(cache *Cache, repo Repo, log *slog.Logger) *Service {
	return &Service{
		cache: cache,
		repo:  repo,
		log:   log,
	}
}

func (s *Service) GetProfile(ctx context.Context, id uuid.UUID) (Profile, error) {
	const op = "profile.Service.GetProfile"
	log := s.log.With(slog.String("op", op), slog.Any("id", id))

	profile, err := s.cache.GetProfile(ctx, id)
	if err == nil {
		log.Info("cache hit!")
		return profile, nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		log.Error("failed to get profile from cache", logerr.Err(err))
	}

	log.Info("cache miss")
	profile, err = s.repo.ExtractProfile(ctx, id)
	if err != nil {
		if !errors.Is(err, ErrNoProfile) {
			log.Error("failed to get profile", logerr.Err(err))
		}
		return Profile{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.cache.SetProfile(ctx, profile)
	if err != nil {
		log.Error("failed to set profile to cache", logerr.Err(err))
	} else {
		log.Info("wrote to cache")
	}

	return profile, nil
}

func (s *Service) CreateProfile(ctx context.Context, username string, name string) error {
	const op = "profile.Service.CreateProfile"
	id, _ := uuid.NewV7()
	profile := Profile{
		ID:       id,
		Username: username,
		Name:     name,
	}
	log := s.log.With(slog.String("op", op), slog.Any("id", id))

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
