package api

import (
	"net/http"

	"allmystuff/internal/store"
)

type TagHandler struct {
	store store.Store
}

func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	tags, err := h.store.ListTags(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tags)
}
