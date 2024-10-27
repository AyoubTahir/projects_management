package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AyoubTahir/projects_management/internal/services"
	"github.com/AyoubTahir/projects_management/pkg/types"
)

type Handler struct {
	Service *services.Service
	User    UserHandlerI
	// Add other service dependencies as needed
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Service: service,
		User:    NewUserHandler(service),
	}
}

type UserHandlerI interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	// Add other user-related methods as needed
}

func JsonResponse(w http.ResponseWriter, status int, response types.RouteResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

// add commen func to Handlers
