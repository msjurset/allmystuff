package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"

	"allmystuff/internal/imgstore"
	"allmystuff/internal/model"
	"allmystuff/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ImageHandler struct {
	store store.Store
	imgs  *imgstore.ImageStore
}

func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	// Verify item exists
	if _, err := h.store.GetItem(r.Context(), itemID); err != nil {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20MB
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageID := uuid.New()
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	relPath, err := h.imgs.Save(itemID, imageID, ext, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get current max sort_order
	existing, _ := h.store.ListImages(r.Context(), itemID)
	sortOrder := len(existing)

	img := model.Image{
		ID:        imageID,
		ItemID:    itemID,
		Filename:  header.Filename,
		Filepath:  relPath,
		SortOrder: sortOrder,
		CreatedAt: time.Now(),
	}

	created, err := h.store.CreateImage(r.Context(), img)
	if err != nil {
		h.imgs.Delete(relPath)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *ImageHandler) Serve(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	img, err := h.store.GetImage(r.Context(), id)
	if err != nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, h.imgs.Path(img.Filepath))
}

func (h *ImageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	img, err := h.store.GetImage(r.Context(), id)
	if err != nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	if err := h.store.DeleteImage(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.imgs.Delete(img.Filepath)
	w.WriteHeader(http.StatusNoContent)
}

func (h *ImageHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	itemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	var req struct {
		ImageIDs []uuid.UUID `json:"image_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.ReorderImages(r.Context(), itemID, req.ImageIDs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
