package cart

import (
	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.HandleFunc("/add-product", AddProductToCart)
	r.Delete("/", RemoveProductFromCart)

	return r 
}