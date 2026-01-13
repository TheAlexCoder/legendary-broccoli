package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"fired-calendar/middleware"
	"fired-calendar/models"
)

type UpdateProfileRequest struct {
	Username  string `json:"username"`
	FiredDate string `json:"fired_date"`
}

type ProfileResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FiredDate string `json:"fired_date"`
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	response := ProfileResponse{
		ID:       user.ID,
		Username: user.Username,
	}
	if user.FiredDate != nil {
		response.FiredDate = *user.FiredDate
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	user.Username = req.Username
	if req.FiredDate != "" {
		// Validate date format
		if _, err := time.Parse("2006-01-02", req.FiredDate); err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		user.FiredDate = &req.FiredDate
	}

	if err := models.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetRecoveryPhraseHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"recovery_phrase": user.RecoveryPhrase,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := models.SoftDeleteUser(userID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Clear session
	if err := middleware.ClearUserSession(w, r); err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RestoreUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if the user was soft-deleted within the last 7 days
	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	if user == nil || !user.IsDeleted {
		http.Error(w, "User not found or not deleted", http.StatusNotFound)
		return
	}

	// Check if deletion was within 7 days
	if user.DeletedAt != nil {
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		if user.DeletedAt.Before(sevenDaysAgo) {
			http.Error(w, "User deletion is older than 7 days, cannot restore", http.StatusBadRequest)
			return
		}
	}

	// Restore user
	if err := models.RestoreUser(userID); err != nil {
		http.Error(w, "Failed to restore user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
