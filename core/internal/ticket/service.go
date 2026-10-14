package ticket

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/google/uuid"
	"log/slog"
)

type Cache interface {
	Get(ctx context.Context, ticket string) (uuid.UUID, error)
	Set(ctx context.Context, ticket string, uid uuid.UUID) error
}

var (
	ErrTicketNotFound = errors.New("ticket doesn't exist")
	ErrInvalidTicket  = errors.New("ticket is invalid")
	ErrInvalidUID     = errors.New("uid is invalid")
)

type Service struct {
	cache Cache
	log   *slog.Logger
}

func NewService(cache Cache, log *slog.Logger) *Service {
	return &Service{
		cache: cache,
		log:   log,
	}
}

func (s *Service) Get(ctx context.Context, ticket string) (uuid.UUID, error) {
	const op = "ticket.Service.Get"
	log := s.log.With(slog.String("op", op))

	if ticket == "" {
		return uuid.Nil, ErrInvalidTicket
	}

	uid, err := s.cache.Get(ctx, ticket)
	if err != nil {
		if errors.Is(err, ErrTicketNotFound) {
			log.Warn("ticket not found")
			return uuid.Nil, err
		}
		log.Error("failed to get ticket", logerr.Err(err))
		return uuid.Nil, err
	}

	return uid, nil
}

func (s *Service) GenTicket(ctx context.Context, uid uuid.UUID) (string, error) {
	const op = "ticket.Service.GenTicket"
	log := s.log.With(slog.String("op", op), slog.String("uid", uid.String()))

	if uid == uuid.Nil {
		log.Warn("invalid uid revealed in service")
		return "", ErrInvalidUID
	}

	ticket := genTicket()
	err := s.cache.Set(ctx, ticket, uid)
	if err != nil {
		log.Error("failed to set ticket", logerr.Err(err))
		return "", err
	}

	return ticket, nil
}

func genTicket() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
