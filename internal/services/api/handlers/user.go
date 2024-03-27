package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lyracampos/go-clean-architecture/internal/domain"
	"github.com/lyracampos/go-clean-architecture/internal/domain/usecases"
	"github.com/lyracampos/go-clean-architecture/internal/services/api"
	"go.uber.org/zap"
)

type userHandler struct {
	log           *zap.SugaredLogger
	listUseCase   usecases.ListUserUseCase
	getUseCase    usecases.GetUserUseCase
	createUseCase usecases.CreateUserUseCase
}

func NewUserHandler(
	log *zap.SugaredLogger,
	listUseCase usecases.ListUserUseCase,
	getUseCase usecases.GetUserUseCase,
	createUseCase usecases.CreateUserUseCase,
) *userHandler {
	return &userHandler{
		log:           log,
		listUseCase:   listUseCase,
		getUseCase:    getUseCase,
		createUseCase: createUseCase,
	}
}

// swagger:route GET /users users ListUsers
// Return a list of users from system
// responses:
//
//	200: userListResponse
//	501: internalServerErrorResponse
func (h *userHandler) ListUsers(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("userHandler.ListUsers - started")

	ctx := context.Background()
	rw.Header().Set("Content-type", "application/json")

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	emailFilter := r.URL.Query().Get("email")
	roleFilter := r.URL.Query().Get("role")

	listUserResult, err := h.listUseCase.Execute(ctx, usecases.ListUserInput{
		FirstName: firstName,
		LastName:  lastName,
		Email:     emailFilter,
		Role:      roleFilter,
	})

	if err != nil {
		errMsg, statusCode := h.handlerErrors(err)
		rw.WriteHeader(statusCode)
		_, err := rw.Write([]byte(errMsg))
		if err != nil {
			log.Printf("userHandler.ListUsers - write failed: %v", err)
		}

		return
	}

	h.log.Info("userHandler.ListUsers - finished successfully")

	if err := json.NewEncoder(rw).Encode(listUserResult); err != nil {
		log.Printf("userHandler.ListUsers - encode failed: %v", err)
	}
}

// swagger:route GET /users/{id} users GetUser
// Return an user from system
// responses:
//
//	200: userGetResponse
//	404: notFoundResponse
//	501: internalServerErrorResponse
func (h *userHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("userHandler.GetUser - started")

	ctx := context.Background()
	rw.Header().Set("Content-type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		errMsg, statusCode := h.handlerErrors(err)
		rw.WriteHeader(statusCode)
		_, err := rw.Write([]byte(errMsg))
		if err != nil {
			log.Printf("userHandler.GetUser - write failed: %v", err)
		}

		return
	}

	getUserResult, err := h.getUseCase.Execute(ctx, usecases.GetUserInput{ID: int64(id)})
	if err != nil {
		errMsg, statusCode := h.handlerErrors(err)
		rw.WriteHeader(statusCode)
		_, err := rw.Write([]byte(errMsg))
		if err != nil {
			log.Printf("userHandler.GetUser - write failed: %v", err)
		}

		return
	}

	h.log.Info("userHandler.GetUser - finished successfully")

	if err := json.NewEncoder(rw).Encode(getUserResult); err != nil {
		log.Printf("userHandler.GetUser - encode failed: %v", err)
	}
}

// swagger:route POST /users  users AddUser
// Add new user in the application
// responses:
//
//	201: userAddResponse
//	400: notFoundResponse
//	501: internalServerErrorResponse
func (h *userHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	rw.Header().Set("Content-type", "application/json")

	requestBody, _ := io.ReadAll(r.Body)
	var input usecases.CreateUserInput
	if err := json.Unmarshal(requestBody, &input); err != nil {
		errMsg, statusCode := h.handlerErrors(err)
		rw.WriteHeader(statusCode)
		_, err := rw.Write([]byte(errMsg))
		if err != nil {
			log.Printf("userHandler.CreateUser - json unmarshal failed: %v", err)
		}

		return
	}

	h.log.Info("userHandler.CreateUser - started")

	createUserResult, err := h.createUseCase.Execute(ctx, usecases.CreateUserInput{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Role:      input.Role,
	})
	if err != nil {
		errMsg, statusCode := h.handlerErrors(err)
		rw.WriteHeader(statusCode)
		_, err := rw.Write([]byte(errMsg))
		if err != nil {
			log.Printf("userHandler.CreateUser - write failed: %v", err)
		}

		return
	}

	h.log.Info("userHandler.CreateUser - finished successfully")

	rw.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(rw).Encode(createUserResult); err != nil {
		log.Printf("userHandler.CreateUser - encode failed: %v", err)
	}
}

func (h *userHandler) handlerErrors(err error) (string, int) {
	h.log.Error(err.Error())

	validationError := &domain.ValidationError{}

	switch {
	case errors.As(err, &validationError):
		return err.Error(), http.StatusBadRequest
	case errors.Is(err, domain.ErrEmailAlreadyInUse):
		return api.ErrEmailAlreadyInUse.Error(), http.StatusConflict
	case errors.Is(err, domain.ErrUserDoesNotExist):
		return api.ErrUserDoesNotExist.Error(), http.StatusNotFound
	default:
		return err.Error(), http.StatusInternalServerError
	}
}
