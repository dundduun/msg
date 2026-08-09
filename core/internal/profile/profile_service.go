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
	Profile(ctx context.Context, id int) (Profile, error)
}

type Service struct {
	repo Repo
	log  *slog.Logger
}

func NewService(repo Repo, log *slog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) GetProfile(ctx context.Context, id int) (Profile, error) {
	const op = "profile.Service.GetProfile"
	log := s.log.With(slog.String("op", op), slog.Int("id", id))

	profile, err := s.repo.Profile(ctx, id)
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
