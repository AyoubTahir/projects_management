package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AyoubTahir/projects_management/config"
	"github.com/AyoubTahir/projects_management/internal/container"
	"github.com/AyoubTahir/projects_management/internal/routes"
)

type Server struct {
	cfg    *config.Config
	server *http.Server
	//router *Router
}

func New(cfg *config.Config, container *container.Container) (*Server, error) {

	// Initialize router
	r := routes.NewRouter(container)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
	}

	return &Server{
		cfg:    cfg,
		server: server,
	}, nil
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
