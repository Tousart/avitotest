package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/avitotest/internal/api"
	"github.com/tousart/avitotest/internal/repository/postgres"
	"github.com/tousart/avitotest/internal/server"
	"github.com/tousart/avitotest/internal/usecase/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error, 1)

	// repository

	address := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSLMODE"),
	)

	teamsRepo, err := postgres.NewTeamsRepository(address)
	if err != nil {
		log.Fatalf("failed to create teams repository")
	}

	usersRepo, err := postgres.NewUsersRepository(address)
	if err != nil {
		log.Fatalf("failed to create users repository")
	}

	pullRequestsRepo, err := postgres.NewPullRequestsRepository(address)
	if err != nil {
		log.Fatalf("failed to create users repository")
	}

	// usecase

	teamsService := service.NewTeamsService(teamsRepo)

	usersService := service.NewUsersService(usersRepo)

	pullRequestsService := service.NewPullRequestsService(pullRequestsRepo)

	// api

	r := chi.NewRouter()

	teamsAPI := api.CreateTeamsAPI(teamsService)
	teamsAPI.WithTeamsHandlers(r)

	usersAPI := api.CreateUsersAPI(usersService)
	usersAPI.WithUsersHandlers(r)

	pullRequestsAPI := api.CreatePullRequestsAPI(pullRequestsService)
	pullRequestsAPI.WithPullRequestsHandlers(r)

	// Запуск сервера

	serv := server.CreateAndRunServer(r, os.Getenv("SERVER_PORT"), errChan)

	select {
	case err := <-errChan:
		log.Printf("server has been stopped: %v\n", err)
	case <-ctx.Done():
		log.Println("starting graceful shutdown...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := serv.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v\n", err)
		} else {
			log.Println("server stopped gracefully")
		}
	}
}
