package http

import (
	"errors"
	"github.com/dundduun/msg/core/pkg/logerr"
	"github.com/dundduun/msg/core/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
	"strconv"
)

type ProfileService interface {
	//GetProfile(id int)
}

type Response struct {
	Profile Profile `json:"profile,omitempty"`
	response.Response
}

type Profile struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type ProfileHandler struct {
	service ProfileService
	conn    *pgx.Conn
	log     *slog.Logger
}

func NewProfileHandler(service ProfileService, conn *pgx.Conn, log *slog.Logger) *ProfileHandler {
	return &ProfileHandler{
		service: service,
		conn:    conn,
		log:     log,
	}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	const op = "/profile/{id}"
	log := h.log.With(slog.String("op", op))

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		log.Warn("bad id number", slog.Int("id", id))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Response{
			Response: response.Error("bad id"),
		})

		return
	}
	row := h.conn.QueryRow(r.Context(), "select id, username, name from profile where id = $1 limit 1", id)

	profile := Profile{}

	err = row.Scan(&profile.ID, &profile.Username, &profile.Name)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		log.Warn("profile not found", slog.Int("id", id))
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, Response{
			Response: response.Error("profile not found"),
		})

		return
	case !errors.Is(err, nil):
		log.Error("failed to get profile", logerr.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, Response{
			Response: response.Error("failed to get profile"),
		})

		return
	}

	render.JSON(w, r, Response{
		Profile:  profile,
		Response: response.OK(),
	})
}
