package main

import (
	"context"
	"errors"
	"github.com/go-chi/render"

	//"github.com/dundduun/msg/core/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"net/http"
)

type Response struct {
	Profile Profile `json:"profile,omitempty"`
	Error   string  `json:"error,omitempty"`
}

type Profile struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

func main() {
	//_ = config.MustLoad()

	conn, err := pgx.Connect(context.Background(), "postgres://admin:example@localhost:5433/msg")
	if err != nil {
		panic("failed to connect: " + err.Error())
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/profile/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		row := conn.QueryRow(r.Context(), "select id, username, name from profile where id = $1 limit 1", id)

		profile := Profile{}

		err = row.Scan(&profile.ID, &profile.Username, &profile.Name)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, Response{
				Error: "no matches",
			})
			return
		case !errors.Is(err, nil):
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, Response{
				Error: err.Error(), // don't show the real error
			})
			return
		}

		render.JSON(w, r, Response{
			Profile: profile,
		})
	})

	http.ListenAndServe(":3000", r)
}
