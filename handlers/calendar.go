package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"fired-calendar/middleware"
	"fired-calendar/models"
	"fired-calendar/utils"
)

type CheckDayRequest struct {
	Date string `json:"date"`
}

type CalendarStats struct {
	TotalWorkingDays int `json:"total_working_days"`
	CheckedDays      int `json:"checked_days"`
	RemainingDays    int `json:"remaining_days"`
}

func GetCalendarDaysHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get current month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	entries, err := models.GetCalendarEntries(userID, startOfMonth.Format("2006-01-02"), endOfMonth.Format("2006-01-02"))
	if err != nil {
		http.Error(w, "Failed to get calendar entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func CheckDayHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CheckDayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := models.CheckDay(userID, req.Date); err != nil {
		http.Error(w, "Failed to check day", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UncheckDayHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CheckDayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := models.UncheckDay(userID, req.Date); err != nil {
		http.Error(w, "Failed to uncheck day", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetStatsHandler(w http.ResponseWriter, r *http.Request) {
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

	var firedDate time.Time = time.Now();
	// Тоже самое????
	// firedDate := time.Now();

	if user.FiredDate != nil {
		// http.Error(w, "Fired date not set", http.StatusBadRequest)
		// return

		//firedDate, err = time.Parse("2006-01-02", *user.FiredDate)
		if err != nil {
			//http.Error(w, "Invalid fired date", http.StatusBadRequest)
			//return
		}
	}

	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Calculate total working days from today to fired date
	totalWorkingDays := utils.CalculateWorkingDays(startOfToday, firedDate)

	// Get checked days count
	checkedDays, err := models.GetCheckedDaysCount(userID, startOfToday.Format("2006-01-02"), firedDate.Format("2006-01-02"))
	if err != nil {
		http.Error(w, "Failed to get checked days count", http.StatusInternalServerError)
		return
	}

	stats := CalendarStats{
		TotalWorkingDays: totalWorkingDays,
		CheckedDays:      checkedDays,
		RemainingDays:    totalWorkingDays - checkedDays,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
