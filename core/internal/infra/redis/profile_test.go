package redis_test

import (
	"context"
	"encoding/json"
	"github.com/dundduun/msg/core/internal/infra/redis"
	"github.com/dundduun/msg/core/internal/profile"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCacheProfile(t *testing.T) {
	db, mock := redismock.NewClientMock()

	t.Run("get profile", func(t *testing.T) {
		v7, _ := uuid.NewV7()
		prof := profile.Profile{
			ID:       v7,
			Username: "1",
			Name:     "1",
		}

		profBytes, err := json.Marshal(prof)
		require.NoError(t, err)

		mock.ExpectGet(redis.PrefixId(prof.ID)).SetVal(string(profBytes))

		cache, err := redis.NewCache(db, 0)
		require.NoError(t, err)
		gotProfile, err := cache.Get(context.Background(), prof.ID)
		require.NoError(t, err)

		assert.Equal(t, prof, gotProfile)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("set profile", func(t *testing.T) {
		ttl := 10 * time.Second
		v7, _ := uuid.NewV7()
		prof := profile.Profile{
			ID:       v7,
			Username: "2",
			Name:     "2",
		}

		mock.ExpectSet(redis.PrefixId(prof.ID), prof, ttl).SetVal("ok")

		cache, err := redis.NewCache(db, ttl)
		require.NoError(t, err)
		err = cache.Set(context.Background(), prof)
		require.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
