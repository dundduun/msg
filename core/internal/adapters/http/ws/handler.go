package ws

import (
	httpp "github.com/dundduun/msg/core/internal/adapters/http"
	"github.com/dundduun/msg/core/internal/adapters/http/middleware/auth"
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

const defaultGroup = "lobby"

type Handler struct {
	log            *slog.Logger
	upgrader       websocket.Upgrader
	hub            *Hub
	profileService httpp.ProfileService
}

func NewHandler(log *slog.Logger, profileService httpp.ProfileService) *Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // handle check origin for production

	return &Handler{
		log:            log,
		upgrader:       upgrader,
		hub:            newHub(log),
		profileService: profileService,
	}
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	const op = "http.ws.Handler.Connect"
	log := h.log.With(slog.String("op", op))

	claims, ok := auth.WithClaims(r.Context())
	if !ok {
		log.Error("can't get claims from context")
		return
	}
	log = log.With(slog.Any("uid", claims.UID))

	profile, err := h.profileService.GetProfile(r.Context(), claims.UID)
	if err != nil {
		log.Error("failed to get profile", logerr.Err(err))
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed to upgrade", logerr.Err(err))
		return
	}
	log.Info("connection open")
	defer func() {
		_ = conn.Close()
	}()

	client := &Client{
		clientID: uuid.New(),
		profile: ProfileData{
			ID:       claims.UID,
			Name:     profile.Name,
			Username: profile.Username,
		},
		group: defaultGroup,
		send:  make(chan Envelope, 256),
		conn:  conn,
		hub:   h.hub,
		log:   h.log,
	}
	h.hub.register <- client

	go client.writePump()
	client.readPump()
}
