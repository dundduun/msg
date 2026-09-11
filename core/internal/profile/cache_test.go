package profile_test

import (
	"context"
	"encoding/json"
	"github.com/dundduun/msg/core/internal/profile"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCacheProfile(t *testing.T) {
	db, mock := redismock.NewClientMock()

	t.Run("get profile", func(t *testing.T) {
		prof := profile.Profile{
			ID:       1,
			Username: "1",
			Name:     "1",
		}

		profBytes, err := json.Marshal(prof)
		require.NoError(t, err)

		mock.ExpectGet(profile.PrefixedID(prof.ID)).SetVal(string(profBytes))

		cache := profile.NewCache(db, 0)
		gotProfile, err := cache.GetProfile(context.Background(), prof.ID)
		require.NoError(t, err)

		assert.Equal(t, prof, gotProfile)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("set profile", func(t *testing.T) {
		ttl := 10 * time.Second
		prof := profile.Profile{
			ID:       2,
			Username: "2",
			Name:     "2",
		}

		mock.ExpectSet(profile.PrefixedID(2), prof, ttl).SetVal("ok")

		cache := profile.NewCache(db, ttl)
		err := cache.SetProfile(context.Background(), prof)
		require.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
