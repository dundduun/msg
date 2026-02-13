package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("email is taken")
)
