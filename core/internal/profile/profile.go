package profile

import "github.com/google/uuid"

type Profile struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
}
