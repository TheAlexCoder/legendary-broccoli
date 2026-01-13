package middleware

import (
	"context"
	"net/http"

	"fired-calendar/config"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(config.SessionKey))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "fired-calendar-session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID.(int))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromSession(r *http.Request) (int, error) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return 0, http.ErrNoCookie
	}
	return userID.(int), nil
}

func SetUserSession(w http.ResponseWriter, r *http.Request, userID int) error {
	session, err := store.Get(r, "fired-calendar-session")
	if err != nil {
		return err
	}

	session.Values["user_id"] = userID
	return session.Save(r, w)
}

func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "fired-calendar-session")
	if err != nil {
		return err
	}

	delete(session.Values, "user_id")
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
