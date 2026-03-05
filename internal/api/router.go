package api

import (
	"allmystuff/internal/imgstore"
	"allmystuff/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(s store.Store, imgs *imgstore.ImageStore, apiKey string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	ih := &ItemHandler{store: s, imgs: imgs}
	th := &TagHandler{store: s}
	imgh := &ImageHandler{store: s, imgs: imgs}

	r.Route("/api", func(r chi.Router) {
		r.Use(APIKeyAuth(apiKey))
		r.Get("/items", ih.List)
		r.Post("/items", ih.Create)
		r.Get("/items/{id}", ih.Get)
		r.Put("/items/{id}", ih.Update)
		r.Delete("/items/{id}", ih.Delete)

		r.Post("/items/{id}/images", imgh.Upload)
		r.Put("/items/{id}/images/order", imgh.Reorder)

		r.Get("/images/{id}", imgh.Serve)
		r.Delete("/images/{id}", imgh.Delete)

		r.Get("/tags", th.List)
	})

	return r
}
