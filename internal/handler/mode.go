package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brobridge/sentinel/internal/k8s"
)

func getMode(store *k8s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mode, err := store.GetMode(r.Context())
		if err != nil {
			http.Error(w, `{"error":"failed to get mode"}`, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"mode": mode})
	}
}

func setMode(store *k8s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Mode string `json:"mode"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}
		if req.Mode != "Monitoring" && req.Mode != "Protect" {
			http.Error(w, `{"error":"mode must be Monitoring or Protect"}`, http.StatusBadRequest)
			return
		}
		if err := store.SetMode(r.Context(), req.Mode); err != nil {
			http.Error(w, `{"error":"failed to set mode"}`, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"mode": req.Mode})
	}
}
