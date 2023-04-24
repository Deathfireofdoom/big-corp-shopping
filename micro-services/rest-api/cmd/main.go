package main


import (
	"net/http"
	"log"
	"context"

	"github.com/go-chi/chi"
	"big-corp-shopping/rest-api/internal/api/cart"
	"big-corp-shopping/rest-api/internal/api/health"
	"big-corp-shopping/rest-api/internal/api/migration"
	
	"big-corp-shopping/rest-api/internal/cart_request_service"
	_ "big-corp-shopping/rest-api/internal/config"
)


func main() {
	// creating context
	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	// starting cart_request_service
	log.Println("[info] starting cart-request-service")
	go func () {
		cart_request_service.Service.Run(ctx)
	}()

	r := chi.NewRouter()

	// Mounting test endpoints
	r.Mount("/test", health.NewRouter())

	// Mounting cart endpoints
	r.Mount("/cart", cart.NewRouter())
	
	// Mounting migration endpoints
	r.Mount("/migration", migration.NewRouter())

	log.Println("Starting server on port 8080...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}