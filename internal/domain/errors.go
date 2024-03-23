package domain

import "errors"

var (
	// ErrUserDoesNotExist is thrown when getting an user that does not exist.
	ErrUserDoesNotExist = errors.New("user does not exist")
	// ErrEmailAlreadyInUse is thrown when inserting new User with an Email that already exists.
	ErrEmailAlreadyInUse = errors.New("email is arealdy in use")
)
