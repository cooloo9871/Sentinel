package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/brobridge/sentinel/internal/auth"
	k8sclient "github.com/brobridge/sentinel/internal/k8s"
)

func loginHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
			return
		}

		creds, err := k8sclient.GetCredentials(r.Context(), cfg.DynClient, cfg.Namespace)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
			return
		}

		hash, ok := creds[req.Username]
		if !ok {
			// Run bcrypt against a dummy hash to prevent timing-based username enumeration.
			hash = "$2a$10$YzJDJjIwMjYwNTI4aGFzaGVzaGVzAAAAAAAAAAAAAAAAAAAAAAAAAA"
		}
		if auth.VerifyPassword(req.Password, hash) != nil || !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}

		token, err := auth.GenerateToken(req.Username, cfg.JWTSecret, cfg.TokenTTL)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "sentinel_token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(cfg.TokenTTL / time.Second),
		})
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "sentinel_token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		})
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}
