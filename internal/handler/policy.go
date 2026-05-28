package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"sigs.k8s.io/yaml"

	"github.com/brobridge/sentinel/internal/k8s"
	"github.com/brobridge/sentinel/internal/policy"
)

type createPolicyRequest struct {
	Source  string                  `json:"source"`           // "form" or "yaml"
	Form    *policy.PolicyFormInput `json:"form,omitempty"`
	Action  string                  `json:"action,omitempty"` // "Post" or "Sigkill"
	RawYAML string                  `json:"rawYaml,omitempty"`
}

func listPolicies(store *k8s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := store.List(r.Context())
		if err != nil {
			http.Error(w, `{"error":"failed to list policies"}`, http.StatusInternalServerError)
			return
		}
		if records == nil {
			records = []k8s.PolicyRecord{}
		}
		writeJSON(w, http.StatusOK, records)
	}
}

func getPolicy(store *k8s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		namespace := r.URL.Query().Get("namespace")

		record, err := store.Get(r.Context(), name, namespace)
		if err != nil {
			http.Error(w, `{"error":"policy not found"}`, http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, record)
	}
}

func createPolicy(store *k8s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPolicyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}

		if req.Source == "yaml" {
			if err := store.ApplyRaw(r.Context(), req.RawYAML); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
			return
		}

		if req.Form == nil {
			http.Error(w, `{"error":"form required"}`, http.StatusBadRequest)
			return
		}
		action := req.Action
		if action == "" {
			action = policy.ActionPost
		}
		tp, err := policy.Build(*req.Form, action)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := store.Apply(r.Context(), tp); err != nil {
			http.Error(w, `{"error":"failed to apply policy"}`, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
	}
}

func updatePolicy(store *k8s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPolicyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}

		if req.Source == "yaml" {
			if err := store.ApplyRaw(r.Context(), req.RawYAML); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
			return
		}

		if req.Form == nil {
			http.Error(w, `{"error":"form required"}`, http.StatusBadRequest)
			return
		}
		action := req.Action
		if action == "" {
			action = policy.ActionPost
		}
		tp, err := policy.Build(*req.Form, action)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := store.Apply(r.Context(), tp); err != nil {
			http.Error(w, `{"error":"failed to apply policy"}`, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
	}
}

func deletePolicy(store *k8s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		namespace := r.URL.Query().Get("namespace")

		if err := store.Delete(r.Context(), name, namespace); err != nil {
			http.Error(w, `{"error":"failed to delete policy"}`, http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

func previewPolicy(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Form   policy.PolicyFormInput `json:"form"`
		Action string                 `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	action := req.Action
	if action == "" {
		action = policy.ActionPost
	}
	tp, err := policy.Build(req.Form, action)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	b, err := yaml.Marshal(tp)
	if err != nil {
		http.Error(w, `{"error":"marshal failed"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"yaml": string(b)})
}
