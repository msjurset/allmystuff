package api

import (
	"encoding/json"
	"net/http"

	"allmystuff/internal/imgstore"
	"allmystuff/internal/model"
	"allmystuff/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ItemHandler struct {
	store store.Store
	imgs  *imgstore.ImageStore
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := store.ItemFilter{
		Query:     r.URL.Query().Get("q"),
		Tag:       r.URL.Query().Get("tag"),
		Condition: r.URL.Query().Get("condition"),
	}
	items, err := h.store.ListItems(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	item, err := h.store.GetItem(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.ItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := h.store.CreateItem(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Fetch existing item to merge with partial input
	existing, err := h.store.GetItem(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Build input from existing item
	input := model.ItemInput{
		Name:           existing.Name,
		Description:    existing.Description,
		Brand:          existing.Brand,
		Model:          existing.Model,
		SerialNumber:   existing.SerialNumber,
		PurchasePrice:  existing.PurchasePrice,
		EstimatedValue: existing.EstimatedValue,
		Condition:      existing.Condition,
		Notes:          existing.Notes,
	}
	if existing.PurchaseDate != nil {
		s := existing.PurchaseDate.Format("2006-01-02")
		input.PurchaseDate = &s
	}
	for _, t := range existing.Tags {
		input.Tags = append(input.Tags, t.Name)
	}

	// Overlay with request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.store.UpdateItem(r.Context(), id, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Get images to clean up files
	images, _ := h.store.ListImages(r.Context(), id)
	if err := h.store.DeleteItem(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Clean up image files
	if len(images) > 0 {
		_ = h.imgs.DeleteItemDir(id)
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
