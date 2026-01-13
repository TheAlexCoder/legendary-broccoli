package main

import (
	"log"
	"net/http"
	"time"

	"fired-calendar/config"
	"fired-calendar/handlers"
	"fired-calendar/middleware"
	"fired-calendar/models"

	"github.com/gorilla/mux"
)

// LoggingMiddleware adds logging for each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture response details
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		requestSize := r.ContentLength
		if requestSize < 0 {
			requestSize = 0
		}

		log.Printf("[%s] %s %s %d %dms %dB -> %dB",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration.Milliseconds(),
			requestSize,
			wrapped.size,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.size == 0 {
		rw.size = len(b)
	} else {
		rw.size += len(b)
	}
	return rw.ResponseWriter.Write(b)
}

func main() {
	// Initialize database
	models.InitDB()

	// Create router
	r := mux.NewRouter()

	// Apply logging middleware to all routes
	r.Use(LoggingMiddleware)

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Auth routes (no auth required)
	api.HandleFunc("/auth/register", handlers.RegisterHandler).Methods("POST")
	api.HandleFunc("/auth/login", handlers.LoginHandler).Methods("POST")

	// Protected routes
	protected := api.NewRoute().Subrouter()
	protected.Use(middleware.AuthMiddleware)

	protected.HandleFunc("/auth/logout", handlers.LogoutHandler).Methods("POST")
	protected.HandleFunc("/calendar/days", handlers.GetCalendarDaysHandler).Methods("GET")
	protected.HandleFunc("/calendar/check", handlers.CheckDayHandler).Methods("POST")
	protected.HandleFunc("/calendar/uncheck", handlers.UncheckDayHandler).Methods("PUT")
	protected.HandleFunc("/calendar/stats", handlers.GetStatsHandler).Methods("GET")
	protected.HandleFunc("/profile", handlers.GetProfileHandler).Methods("GET")
	protected.HandleFunc("/profile", handlers.UpdateProfileHandler).Methods("PUT")
	protected.HandleFunc("/profile/recovery", handlers.GetRecoveryPhraseHandler).Methods("GET")
	protected.HandleFunc("/profile/restore", handlers.RestoreUserHandler).Methods("POST")
	protected.HandleFunc("/user/delete", handlers.DeleteUserHandler).Methods("DELETE")

	// Serve index.html for all other routes
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// Start server
	log.Printf("Server starting on port %s", config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, r))
}
