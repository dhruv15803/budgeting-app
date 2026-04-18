package handlers

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data"`
}

type ApiError struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func writeJsonResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string) error {
	return writeJsonResponse(w, status, ApiError{Success: false, Message: message, StatusCode: status})
}
