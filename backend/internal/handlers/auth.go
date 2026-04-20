package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/services"
)

type meResponse struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Username  *string    `json:"username"`
	ImageURL  *string    `json:"image_url"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type registerRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Username *string `json:"username"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type googleAuthRequest struct {
	Credential string `json:"credential"`
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

func (h *Handler) GoogleOAuth(w http.ResponseWriter, r *http.Request) {
	var req googleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	token, err := h.services.Users.LoginWithGoogle(r.Context(), req.Credential)
	if errors.Is(err, services.ErrGoogleAuthDisabled) {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, services.ErrInvalidGoogleToken) {
		_ = writeJsonError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if errors.Is(err, services.ErrGoogleEmailNotVerified) {
		_ = writeJsonError(w, http.StatusForbidden, err.Error())
		return
	}
	if errors.Is(err, services.ErrGoogleAccountConflict) {
		_ = writeJsonError(w, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		log.Printf("GoogleOAuth: %v", err)
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

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.services.Users.GetMe(claims.UserID)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	data := meResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		ImageURL:  user.ImageURL,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	type meResp struct {
		Success bool       `json:"success"`
		Message string     `json:"message"`
		Data    meResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, meResp{Success: true, Message: "ok", Data: data}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}
