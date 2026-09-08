package profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var prefix = "cache:profile"

var ErrCacheMiss = errors.New("cache has no profile with this id")

type Cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

func (c *Cache) GetProfile(ctx context.Context, id int) (Profile, error) {
	const op = "profile.Cache.GetProfile"

	profileJSON, err := c.rdb.Get(ctx, fmt.Sprintf("%s:%d", prefix, id)).Result()
	if errors.Is(err, redis.Nil) {
		return Profile{}, ErrCacheMiss
	} else if err != nil {
		return Profile{}, fmt.Errorf("%s: %w", op, err)
	}

	profile := Profile{}
	err = json.Unmarshal([]byte(profileJSON), &profile)
	if err != nil {
		return Profile{}, fmt.Errorf("%s: %w", op, err)
	}

	return profile, nil
}

//func (c *Cache) SetProfile(ctx context.Context, profile Profile) (int, error) {
//
//}
