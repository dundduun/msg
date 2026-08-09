package postgres

import (
	"context"
	"errors"
	"fmt"
	prof "github.com/dundduun/msg/core/internal/profile"
	"github.com/jackc/pgx/v5"
)

type ProfileRepo struct {
	conn *pgx.Conn
}

func NewProfileRepo(conn *pgx.Conn) *ProfileRepo {
	return &ProfileRepo{conn: conn}
}

func (p *ProfileRepo) Profile(ctx context.Context, id int) (prof.Profile, error) {
	const op = "postgres.ProfileRepo.Profile"

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
