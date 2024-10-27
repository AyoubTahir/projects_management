package routes

import (
	"github.com/AyoubTahir/projects_management/internal/handlers"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router, handler *handlers.Handler) {
	r.HandleFunc("/users", handler.User.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handler.User.GetUser).Methods("GET")
	// Add other user-related routes here
}
