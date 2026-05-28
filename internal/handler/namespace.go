package handler

import (
	"net/http"

	"github.com/brobridge/sentinel/internal/k8s"
)

func listNamespaces(store *k8s.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		names, err := store.ListNamespaces(r.Context())
		if err != nil {
			http.Error(w, `{"error":"failed to list namespaces"}`, http.StatusInternalServerError)
			return
		}
		if names == nil {
			names = []string{}
		}
		writeJSON(w, http.StatusOK, names)
	}
}
