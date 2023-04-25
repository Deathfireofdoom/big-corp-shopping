package cart

import (
	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/update-product", UpdateProductToCart)
	return r 
}