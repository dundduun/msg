package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCacheSetProfile(t *testing.T) {
	db, mock := redismock.NewClientMock()

	id := 1
	profile := Profile{
		ID:       id,
		Username: "1",
		Name:     "1",
	}

	profBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	mock.ExpectGet(fmt.Sprintf("%s:%d", prefix, 1)).SetVal(string(profBytes))

	cache := NewCache(db)
	gotProfile, err := cache.GetProfile(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, profile, gotProfile)
}
