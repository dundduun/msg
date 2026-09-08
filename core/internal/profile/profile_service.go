package profile

import (
	"context"
	"errors"
	"fmt"
	"github.com/dundduun/msg/core/pkg/logerr"
	"log/slog"
)

var ErrNoProfile = errors.New("profile not found")

type Repo interface {
	ExtractProfile(ctx context.Context, id int) (Profile, error)
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

func (s *Service) GetProfile(ctx context.Context, id int) (Profile, error) {
	const op = "profile.Service.GetProfile"
	log := s.log.With(slog.String("op", op), slog.Int("id", id))

	profile, err := s.cache.GetProfile(ctx, id)
	if err == nil {
		return profile, nil
	}
	if !errors.Is(err, ErrCacheMiss) {
		log.Error("failed to get profile from cache", logerr.Err(err))
	}

	profile, err = s.repo.ExtractProfile(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoProfile):
			return Profile{}, fmt.Errorf("%s: %w", op, err)
		default:
			log.Error("failed to get profile", logerr.Err(err))
			return Profile{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return profile, nil
}

// мб не стоит сюда смотреть, а написать по-новой чтоб ничего не упустить
//profile, err := s.cache.Profile(ctx, id)
//if err == nil { // maybe switch
//return profile ...
//} else {
//if err == ErrNoProfileCache {
// give it to go further, but remember to write to cache a record
//} else
//500 error, we should to do something
//}
//}
