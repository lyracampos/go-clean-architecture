package api

import "errors"

var (
	ErrUserDoesNotExist  = errors.New("no user was found for this ID. Please check the ID and try again")
	ErrEmailAlreadyInUse = errors.New("email in use by another user")
)
