package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"k8s.io/client-go/dynamic"

	"github.com/brobridge/sentinel/internal/auth"
	"github.com/brobridge/sentinel/internal/k8s"
)

// Config holds dependencies for all handlers.
type Config struct {
	Store     *k8s.Store
	DynClient dynamic.Interface
	JWTSecret []byte
	Namespace string        // k8s namespace where Sentinel runs
	TokenTTL  time.Duration // JWT lifetime
}

// New builds the HTTP handler tree.
func New(cfg Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/auth/login", loginHandler(cfg))
	r.Post("/api/auth/logout", logoutHandler())

	r.Group(func(r chi.Router) {
		r.Use(jwtMiddleware(cfg.JWTSecret))
		r.Get("/api/policies", listPolicies(cfg.Store))
		r.Post("/api/policies", createPolicy(cfg.Store))
		r.Post("/api/policies/preview", previewPolicy)
		r.Get("/api/policies/{name}", getPolicy(cfg.Store))
		r.Put("/api/policies/{name}", updatePolicy(cfg.Store))
		r.Delete("/api/policies/{name}", deletePolicy(cfg.Store))
		r.Get("/api/namespaces", listNamespaces(cfg.Store))
		r.Get("/api/mode", getMode(cfg.Store))
		r.Put("/api/mode", setMode(cfg.Store))
	})

	return r
}

func jwtMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("sentinel_token")
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
				return
			}
			if _, err := auth.ValidateToken(cookie.Value, secret); err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
