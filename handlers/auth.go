package handlers

import (
	"encoding/json"
	"net/http"

	"fired-calendar/middleware"
	"fired-calendar/models"
	"fired-calendar/utils"
)

type RegisterRequest struct {
	Username string `json:"username"`
}

type RegisterResponse struct {
	RecoveryPhrase string `json:"recovery_phrase"`
}

type LoginRequest struct {
	RecoveryPhrase string `json:"recovery_phrase"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	recoveryPhrase := utils.GenerateRecoveryPhrase()

	user, err := models.CreateUser(req.Username, recoveryPhrase)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Set session
	if err := middleware.SetUserSession(w, r, user.ID); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{RecoveryPhrase: recoveryPhrase}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if !utils.IsValidRecoveryPhrase(req.RecoveryPhrase) {
		http.Error(w, "Invalid recovery phrase", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByRecoveryPhrase(req.RecoveryPhrase)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Set session
	if err := middleware.SetUserSession(w, r, user.ID); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := middleware.ClearUserSession(w, r); err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
