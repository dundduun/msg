package profile

import (
	"context"
	"errors"
	"fmt"
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

var ErrNoProfile = errors.New("profile doesn't exist")

type Service struct {
	conn *pgx.Conn // убрать
	log  *slog.Logger
}

func NewService(conn *pgx.Conn, log *slog.Logger) *Service {
	return &Service{
		conn: conn,
		log:  log,
	}
}

func (s *Service) GetProfile(ctx context.Context, id int) (Profile, error) {
	const op = "profile.Service.GetProfile"
	log := s.log.With(slog.String("op", op), slog.Int("id", id))

	row := s.conn.QueryRow(ctx, "select id, username, name from profile where id = $1 limit 1", id)

	profile := Profile{}

	err := row.Scan(&profile.ID, &profile.Username, &profile.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Profile{}, ErrNoProfile
		} else {
			log.Error("failed to get profile", logerr.Err(err))
			return Profile{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return profile, err
}
