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
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}

		creds, err := k8sclient.GetCredentials(r.Context(), cfg.DynClient, cfg.Namespace)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}

		hash, ok := creds[req.Username]
		if !ok || auth.VerifyPassword(req.Password, hash) != nil {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(req.Username, cfg.JWTSecret, cfg.TokenTTL)
		if err != nil {
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "sentinel_token",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
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
