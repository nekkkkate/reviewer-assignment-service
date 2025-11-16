package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reviewer-assignment-service/internal/app/config"
	"reviewer-assignment-service/internal/app/routes"
	"reviewer-assignment-service/internal/domain/services/impl"
	"reviewer-assignment-service/internal/infrastructure/database"
	"reviewer-assignment-service/internal/infrastructure/persistence/postgres"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserDataBase(db)
	teamRepo := postgres.NewTeamDataBase(db)
	pullRequestRepo := postgres.NewPullRequestDataBase(db)

	userService := impl.NewUserService(userRepo)
	teamService := impl.NewTeamService(teamRepo)
	pullRequestService := impl.NewPullRequestService(pullRequestRepo)

	router := routes.SetupRouter(userService, pullRequestService, teamService)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
