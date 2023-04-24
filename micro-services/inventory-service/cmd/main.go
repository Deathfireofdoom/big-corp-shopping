package main

import (
	"log"
	"context"

	// internal
	"inventory/internal/inventory_service"
)

func main() {
	log.Println("[info] creating inventory service")
	inventoryService := inventory_service.NewInventoryService()
	log.Println("[info] creating context")
	ctx := context.Background()
	log.Println("[info] starting inventory service")
	inventoryService.Run(ctx)
}