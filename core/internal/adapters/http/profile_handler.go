package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	prof "github.com/dundduun/msg/core/internal/profile"
	"github.com/dundduun/msg/core/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type ProfileService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (prof.Profile, error)
	CreateProfile(ctx context.Context, username string, name string) error
}

type ProfileHandler struct {
	service  ProfileService
	validate *validator.Validate
}

func NewProfileHandler(service ProfileService) *ProfileHandler {
	return &ProfileHandler{
		service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

type GetProfileResponse struct {
	Profile prof.Profile `json:"profile,omitzero"`
	response.Response
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "id")
	id, err := uuid.Parse(param)
	if err != nil {
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

	render.JSON(w, r, GetProfileResponse{
		Profile:  profile,
		Response: response.OK(),
	})
}

type CreateProfileRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
}

func (h *ProfileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var req CreateProfileRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errRes(w, r, http.StatusBadRequest, fmt.Sprintf("bad request body: %s", err.Error()))
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)
	err = h.validate.Struct(req)
	if err != nil {
		errRes(w, r, http.StatusBadRequest, fmt.Sprintf("bad request body: %s", err.Error()))
		return
	}

	err = h.service.CreateProfile(r.Context(), req.Username, req.Name)
	if err != nil {
		if errors.Is(err, prof.ErrUsernameTaken) {
			errRes(w, r, http.StatusConflict, err.Error())
		} else {
			errRes(w, r, http.StatusInternalServerError, "failed to create profile")
		}

		return
	}

	render.Status(r, http.StatusCreated)
}

func errRes(w http.ResponseWriter, r *http.Request, status int, err string) {
	render.Status(r, status)
	render.JSON(w, r, GetProfileResponse{
		Response: response.Error(err),
	})
}
