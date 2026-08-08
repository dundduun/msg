package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	prof "github.com/dundduun/msg/core/internal/profile"
	"github.com/dundduun/msg/core/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ProfileService interface {
	GetProfile(ctx context.Context, id int) (prof.Profile, error)
}

type ProfileHandler struct {
	service ProfileService
}

func NewProfileHandler(service ProfileService) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

type Response struct {
	Profile prof.Profile `json:"profile,omitzero"`
	response.Response
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		errRes(w, r, http.StatusBadRequest, "bad id")
		return
	}

	profile, err := h.service.GetProfile(r.Context(), id)
	if err != nil {
		if errors.Is(err, prof.ErrNoProfile) {
			errRes(w, r, http.StatusNotFound, "profile not found")
		} else {
			errRes(w, r, http.StatusInternalServerError, "failed to get profile")
		}

		return
	}

	render.JSON(w, r, Response{
		Profile:  profile,
		Response: response.OK(),
	})
}

func errRes(w http.ResponseWriter, r *http.Request, status int, err string) {
	render.Status(r, status)
	render.JSON(w, r, Response{
		Response: response.Error(err),
	})
}
