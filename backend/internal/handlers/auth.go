package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dhruv15803/budgeting-app/internal/services"
)

type registerRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Username *string `json:"username"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	err := h.services.Users.Register(req.Email, req.Password, req.Username)
	if errors.Is(err, services.ErrEmailTaken) {
		_ = writeJsonError(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := writeJsonResponse(w, http.StatusCreated, ApiResponse{
		Success: true,
		Message: "Registration successful. Check your email to verify your account.",
		Data:    nil,
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	token, err := h.services.Users.Login(req.Email, req.Password)
	if errors.Is(err, services.ErrInvalidCredentials) {
		_ = writeJsonError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if errors.Is(err, services.ErrNotVerified) {
		_ = writeJsonError(w, http.StatusForbidden, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var resp authTokenResponse
	resp.Success = true
	resp.Message = "Login successful"
	resp.Data.Token = token
	if err := writeJsonResponse(w, http.StatusOK, resp); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	out, err := h.services.Users.VerifyEmail(token)
	if errors.Is(err, services.ErrInvalidOrExpiredToken) {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var resp authTokenResponse
	resp.Success = true
	resp.Message = "Email verified"
	resp.Data.Token = out
	if err := writeJsonResponse(w, http.StatusOK, resp); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}
