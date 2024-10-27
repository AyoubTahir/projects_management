package handlers

import (
	"net/http"
	"strconv"

	"github.com/AyoubTahir/projects_management/internal/services"
	"github.com/AyoubTahir/projects_management/pkg/types"
	"github.com/AyoubTahir/projects_management/pkg/validator"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service   *services.Service
	Validator *validator.Validator
}

func NewUserHandler(service *services.Service) UserHandlerI {
	return &UserHandler{
		service:   service,
		Validator: validator.New(),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.CreateUserPayload

	if err := ParseJSON(r, &user); err != nil {
		JsonResponse(w, http.StatusBadRequest, types.RouteResponse{
			Status:  false,
			Message: "Missing request body",
			Errors:  err.Error(),
		})
		return
	}

	if err := h.Validator.Validate(user); err != nil {
		JsonResponse(w, http.StatusUnprocessableEntity, types.RouteResponse{
			Status:  false,
			Message: "Validation error",
			Errors:  h.Validator.GetErrors(),
		})
		return
	}

	data, err := h.service.User.CreateUser(r.Context(), &user)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		JsonResponse(w, http.StatusInternalServerError, types.RouteResponse{
			Status:  false,
			Message: "Something went wrong",
			Errors:  err.Error(),
		})
		return
	}

	JsonResponse(w, http.StatusCreated, types.RouteResponse{
		Status:  true,
		Message: "User created successfully",
		Data:    data,
	})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		JsonResponse(w, http.StatusBadRequest, types.RouteResponse{
			Status:  false,
			Message: "Invalid user ID",
			Errors:  err.Error(),
		})
		return
	}

	user, err := h.service.User.GetUserByID(r.Context(), id)
	if err != nil {
		JsonResponse(w, http.StatusInternalServerError, types.RouteResponse{
			Status:  false,
			Message: "Failed to get user",
			Errors:  err.Error(),
		})
		return
	}

	JsonResponse(w, http.StatusOK, types.RouteResponse{
		Status:  true,
		Message: "User retrieved successfully",
		//add a fake map to data
		Data: user,
	})
}
