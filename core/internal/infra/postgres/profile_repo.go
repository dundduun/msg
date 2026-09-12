package postgres

import (
	"context"
	"errors"
	"fmt"
	prof "github.com/dundduun/msg/core/internal/profile"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProfileRepo struct {
	conn *pgx.Conn
}

func NewProfileRepo(conn *pgx.Conn) *ProfileRepo {
	return &ProfileRepo{conn: conn}
}

func (p *ProfileRepo) ExtractProfile(ctx context.Context, id uuid.UUID) (prof.Profile, error) {
	const op = "postgres.ProfileRepo.ExtractProfile"

	row := p.conn.QueryRow(ctx, "select id, username, name from profile where id = $1", id)

	profile := prof.Profile{}
	err := row.Scan(&profile.ID, &profile.Username, &profile.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return prof.Profile{}, prof.ErrNoProfile
		} else {
			return prof.Profile{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	return profile, nil
}

var CodeUniqueViolation = "23505"

func (p *ProfileRepo) InsertProfile(ctx context.Context, profile prof.Profile) error {
	const op = "postgres.ProfileRepo.InsertProfile"

	_, err := p.conn.Exec(ctx,
		`insert into profile (id, username, name) values ($1, $2, $3)`,
		profile.ID,
		profile.Username,
		profile.Name,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == CodeUniqueViolation {
			return prof.ErrUsernameTaken
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
