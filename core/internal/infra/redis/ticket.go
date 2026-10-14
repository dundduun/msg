package redis

import (
	"context"
	"errors"
	"fmt"
	ticketp "github.com/dundduun/msg/core/internal/ticket"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

const prefix = "ws:ticket:"

var (
	ErrNilRDB = errors.New("rdb should not be nil")
	ErrBadTTL = errors.New("ttl should be bigger than 0")
)

type TicketCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewTicketCache(rdb *redis.Client, ttl time.Duration) (*TicketCache, error) {
	const op = "redis.NewTicketCache"

	if rdb == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrNilRDB)
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("%s: %w", op, ErrBadTTL)
	}

	return &TicketCache{rdb: rdb, ttl: ttl}, nil
}

func (t *TicketCache) Get(ctx context.Context, ticket string) (uuid.UUID, error) {
	const op = "redis.TicketCache.Get"

	uidStr, err := t.rdb.Get(ctx, prefixTicket(ticket)).Result()
	if errors.Is(err, redis.Nil) {
		return uuid.Nil, ticketp.ErrTicketNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: invalid cached UUID: %w", op, err)
	}

	return uid, nil
}

func (t *TicketCache) Set(ctx context.Context, ticket string, uid uuid.UUID) error {
	const op = "redis.TicketCache.GenTicket"

	err := t.rdb.Set(ctx, prefixTicket(ticket), uid.String(), t.ttl).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func prefixTicket(ticket string) string {
	return prefix + ticket
}
