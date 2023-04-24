package migration

import (
	"github.com/go-chi/chi"

)


func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/inventory", InitializeInventoryHandler)

	return r
}