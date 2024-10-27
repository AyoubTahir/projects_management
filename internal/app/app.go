package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AyoubTahir/projects_management/config"
	"github.com/AyoubTahir/projects_management/internal/container"
	"github.com/AyoubTahir/projects_management/internal/server"
)

type App struct {
	cfg       *config.Config
	container *container.Container
	server    *server.Server
}

func New(cfg *config.Config) (*App, error) {
	// Initialize dependency container
	container, err := container.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize container: %w", err)
	}

	// Initialize server
	server, err := server.New(cfg, container)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server: %w", err)
	}

	return &App{
		cfg:       cfg,
		container: container,
		server:    server,
	}, nil
}

func (a *App) Start() error {
	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	log.Println("Starting server...")
	log.Printf("Listening on %s", a.cfg.Server.Port)
	log.Printf("Press Ctrl+C to gracefully shut down the server")
	go func() {
		if err := a.server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	return a.Shutdown()
}

func (a *App) Shutdown() error {
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Shutdown server
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	// Close container resources
	if err := a.container.Close(); err != nil {
		return fmt.Errorf("failed to close container: %w", err)
	}

	return nil
}
