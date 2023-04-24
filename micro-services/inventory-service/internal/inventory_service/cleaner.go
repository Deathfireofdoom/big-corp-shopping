package inventory_service

import (
	"context"

	// internal
	"inventory/internal/entity"
)


type cleaner struct {
	responseChannel chan entity.InventoryRequestResponse
}

func (c *cleaner) Run(ctx context.Context) {
	panic("implement this")
}

// cleanUp will run periodaclly and release any old holds and communicate this back to
// the cart service.
func (c *cleaner) cleanUp() {
	panic("implement this")
}