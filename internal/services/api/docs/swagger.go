// nolint:unused
// Package classification User API
//
// Documentation for User API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package docs

import (
	"github.com/lyracampos/go-clean-architecture/internal/domain/entities"
	"github.com/lyracampos/go-clean-architecture/internal/domain/usecases"
)

// Internal server error message returned as a string
// swagger:response internalServerErrorResponse
type internalServerErrorResponseWrapper struct {
	// error description
	// in: body
	Body MessageError
}

// Not found message error returned as string
// swagger:response notFoundResponse
type errorNotFoundResponseWrapper struct {
	// error description
	// in: body
	Body MessageError
}

// BadRequest message error returned as string
// swagger:response badRequestResponse
type badRequestResponseWrapper struct {
	// error description
	// in: body
	Body MessageError
}

type MessageError struct {
	Message string `json:"message"`
}

// Data structure representing an user
// swagger:response userGetResponse
type userGetResponseWrapper struct {
	// in: body
	Body entities.User
}

// Data structure representing a list of user
// swagger:response userListResponse
type userListResponseWrapper struct {
	// Newly created product
	// in: body
	Body []entities.User
}

// Data structure representing user added
// swagger:response userAddResponse
type userAddResponseWrapper struct {
	// in: body
	Body usecases.CreateUserOutput
}

// swagger:parameters AddUser
type userAddCommandWrapper struct {
	// Payload to add new user in application
	// in: body
	// required: true
	Body usecases.CreateUserInput
}
