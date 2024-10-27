package routes

import (
	"github.com/AyoubTahir/projects_management/internal/container"
	"github.com/gorilla/mux"
)

func NewRouter(container *container.Container) *mux.Router {
	r := mux.NewRouter()

	RegisterUserRoutes(r, container.Handler)
	// Register other routes here (e.g., order routes)

	return r
}
