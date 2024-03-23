package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lyracampos/go-clean-architecture/internal/domain/usecases"
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

func (h *userHandler) ListUsers(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("userHandler.ListUsers - started")

	ctx := context.Background()
	rw.Header().Set("Content-type", "application/json")

	listUserResult, err := h.listUseCase.Execute(ctx, usecases.ListUserInput{
		Email: "",
		Role:  "",
	})

	if err != nil {
		handledError, statusCode := h.handlerErrors(err)
		rw.Write([]byte(handledError.Error()))
		rw.WriteHeader(statusCode)

		return
	}

	h.log.Info("userHandler.ListUsers - finished successfully")

	json.NewEncoder(rw).Encode(listUserResult)
}

func (h *userHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("userHandler.GetUser - started")

	ctx := context.Background()
	rw.Header().Set("Content-type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handledError, statusCode := h.handlerErrors(err)
		rw.Write([]byte(handledError.Error()))
		rw.WriteHeader(statusCode)

		return
	}

	getUserResult, err := h.getUseCase.Execute(ctx, usecases.GetUserInput{ID: int64(id)})
	if err != nil {
		handledError, statusCode := h.handlerErrors(err)
		rw.Write([]byte(handledError.Error()))
		rw.WriteHeader(statusCode)

		return
	}

	h.log.Info("userHandler.GetUser - finished successfully")

	json.NewEncoder(rw).Encode(getUserResult)
}

func (h *userHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("userHandler.CreateUser - started")

	ctx := context.Background()
	rw.Header().Set("Content-type", "application/json")

	requestBody, _ := io.ReadAll(r.Body)
	var input usecases.CreateUserInput
	json.Unmarshal(requestBody, &input)

	createUserResult, err := h.createUseCase.Execute(ctx, usecases.CreateUserInput{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Role:      input.Role,
	})
	if err != nil {
		handledError, statusCode := h.handlerErrors(err)
		rw.Write([]byte(handledError.Error()))
		rw.WriteHeader(statusCode)

		return
	}

	h.log.Info("userHandler.CreateUser - finished successfully")

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(createUserResult)
}

func (h *userHandler) handlerErrors(err error) (error, int) {
	return nil, http.StatusNotFound
}
