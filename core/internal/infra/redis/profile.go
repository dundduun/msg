package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	profilep "github.com/dundduun/msg/core/internal/profile"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var profilePrefix = "cache:profile:"

type ProfileCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewCache(rdb *redis.Client, ttl time.Duration) (*ProfileCache, error) {
	const op = "profile.NewCache"

	if rdb == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrNilRDB)
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrBadTTL)
	}

	return &ProfileCache{rdb: rdb, ttl: ttl}, nil
}

func (c *ProfileCache) Get(ctx context.Context, id uuid.UUID) (profilep.Profile, error) {
	const op = "profile.ProfileCache.Check"

	profileJSON, err := c.rdb.Get(ctx, prefixID(id)).Result()
	if errors.Is(err, redis.Nil) {
		return profilep.Profile{}, profilep.ErrCacheMiss
	}
	if err != nil {
		return profilep.Profile{}, fmt.Errorf("%s: %w", op, err)
	}

	profile := profilep.Profile{}
	err = json.Unmarshal([]byte(profileJSON), &profile)
	if err != nil {
		return profilep.Profile{}, fmt.Errorf("%s: %w", op, err)
	}

	return profile, nil
}

func (c *ProfileCache) Set(ctx context.Context, profile profilep.Profile) error {
	const op = "profile.ProfileCache.GenTicket"

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = c.rdb.Set(ctx, prefixID(profile.ID), profileJSON, c.ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func prefixID(id uuid.UUID) string {
	return fmt.Sprintf("%s%v", profilePrefix, id)
}
