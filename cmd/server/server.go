package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	userHandler "user.com/m/internal/handler/user"
	"user.com/m/internal/middleware"
	userRepo "user.com/m/internal/repository/user"
	userService "user.com/m/internal/service/user"
)

func main() {
	router := mux.NewRouter()
	userRepository := userRepo.NewUserRepository()
	userService := userService.NewUserService(userRepository)
	userHandler := userHandler.NewUserHandler(userService)

	router.HandleFunc("/api/users", userHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/api/users/{id}", userHandler.GetUserById).Methods(http.MethodGet)
	router.HandleFunc("/api/users", userHandler.GetUsers).Methods(http.MethodGet)

	handlerChain := middleware.RequestIDMiddleware(
		middleware.RecoveryMiddleware(
			middleware.LoggingMiddleware(
				middleware.TimingMiddleware(router),
			),
		),
	)

	srv := &http.Server{
		Handler:      handlerChain,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting server on %s", srv.Addr)
	log.Printf("Available endpoints:")
	log.Printf("  POST /api/users - Create new user")
	log.Printf("  GET  /api/users - Get all users")
	log.Printf("  GET  /api/users/{id} - Get user by ID")
	log.Fatal(srv.ListenAndServe())
}
