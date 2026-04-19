package handlers

import (
	"net/http"
	"strconv"

	"github.com/dhruv15803/budgeting-app/internal/auth"
	"github.com/dhruv15803/budgeting-app/internal/middleware"
	"github.com/dhruv15803/budgeting-app/internal/services"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func claimsFromRequest(r *http.Request) *auth.Claims {
	return middleware.ClaimsFromContext(r.Context())
}

func (h *Handler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.services.Users.DeleteUserById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := writeJsonResponse(w, http.StatusOK, ApiResponse{Success: true, Message: "Health check successful"}); err != nil {
		writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}
