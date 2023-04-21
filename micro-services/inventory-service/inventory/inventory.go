// This file contains the logic used to handle incoming request to the Inventory service.
package Inventory

import (
	"fmt"
	"inventory/entity"
	"context"
)

var Service = &Inventory{}

type Inventory struct {
}

// Todo make chan one way
func (i *Inventory) HandleRequest(ctx context.Context, in chan entity.InventoryRequest) chan entity.InventoryResult {
	out := make(chan entity.InventoryResult)
	const workers = 10

	// Spinning up workers.
	for k := 0; k < workers; k++ {
		go i.worker(ctx, in, out)
	}

	return out
}

func (i *Inventory) worker(ctx context.Context, in chan entity.InventoryRequest, out chan entity.InventoryResult) {
	for {
		select {
		case req := <- in:
			fmt.Println("Handling incoming request.")
			// Mock - TODO CHANGE THIS.
			result := entity.InventoryResult{Status: entity.Success, RequestID: req.RequestID, Request: req}
			out <- result
		
		case <-ctx.Done():
			// Gracefully escapes
			return 
		}
	}

}