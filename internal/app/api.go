package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/lyracampos/go-clean-architecture/config"
	"github.com/lyracampos/go-clean-architecture/internal/domain"
	"github.com/lyracampos/go-clean-architecture/internal/domain/usecases"
	"github.com/lyracampos/go-clean-architecture/internal/gateways/postgres"
	"github.com/lyracampos/go-clean-architecture/internal/services/api/handlers"
	"go.uber.org/zap"
)

func RunAPI(config *config.Config) {
	router := mux.NewRouter()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("failed to defer logger sync: %v", err)
		}
	}()

	sugar := logger.Sugar()

	// database
	postgresClient, err := postgres.NewClient(sugar, config)
	if err != nil {
		log.Fatalf("can't initialize postgres client: %v", err)
	}

	userDatabaseGateway := postgres.NewUserDatabase(postgresClient)

	validator := domain.NewValidatorService()

	// use cases
	listUserUseCase := usecases.NewListUserUseCase(userDatabaseGateway)
	getUserUseCase := usecases.NewGetUserUseCase(userDatabaseGateway)
	createUserUseCase := usecases.NewCreateUserUseCase(userDatabaseGateway, validator)

	// health handler
	healthHandler := handlers.NewHealthHandler(sugar)
	healthRouter := router.Methods(http.MethodGet).Subrouter()
	healthRouter.HandleFunc("/health", healthHandler.CheckStatus)

	// user handlers
	userHandler := handlers.NewUserHandler(sugar, listUserUseCase, getUserUseCase, createUserUseCase)

	listUserRouter := router.Methods(http.MethodGet).Subrouter()
	listUserRouter.HandleFunc("/users", userHandler.ListUsers)

	getUserRouter := router.Methods(http.MethodGet).Subrouter()
	getUserRouter.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser)

	createUserRouter := router.Methods(http.MethodPost).Subrouter()
	createUserRouter.HandleFunc("/users", userHandler.CreateUser)

	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)

	http.Handle("/", router)

	address := fmt.Sprintf("%s:%d", config.API.Host, config.API.Port)

	server := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		sugar.Infof("running API HTTP server at: %s", address)

		if err := server.ListenAndServe(); err != nil {
			sugar.Errorln("Error starting server: %w", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.Shutdown(ctx); err != nil {
		sugar.Errorln("Error shutting down server: %w", err)
	}

	log.Println("shutting down")
}
