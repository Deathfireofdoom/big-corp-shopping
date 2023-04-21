package cart_service 


func (c *CartService) handleInventoryResponse(msg kafka.Message) kafka.Message {
	// convert kafka message to struct
	var inventoryRequest InventoryRequest
	err := json.Unmarshal(msg.Value, &inventoryRequest)
	if err != nil {
		log.Printf("[warning] could not unmarshal inventory-request")
		return nil 
	}

	// check if response was successful
	if inventoryRequest.response.status != 200 {
		log.Printf("[debug] request %s was not successful due to %s", inventoryRequest.requestID, inventoryRequest.response.message)

	} else {
		switch inventoryRequest.Action {
		case entity.InventoryRequestHold:
			log.Printf("[debug] got response from inventory-service for request %s for user %s", inventoryRequest.requestID, inventoryRequest.UserID)
			err := c.applyInventoryRequest(inventoryRequest)
			if err != nil {
				log.Printf("[warning] could not apply inventory request %s for user %s", inventoryRequest.requestID, inventoryRequest.UserID)
				inventoryRequest.response.status_code = 501
				inventoryRequest.response.message = err
			}
		
		case entity.InventoryRequestRelease:
			panic("Implement this")
		
		default:
			log.Printf("[warning] unknown action %s in inventory request response %s", inventoryRequest.Action, inventoryRequest.requestID)
			return nil 
	}

	}

	// convert inventoryRequest to cartRequestResponse

	// serialize as kafka and return


}


func (c *CartService) applyInventoryRequest(inventoryRequest entity.InventoryRequest) error {
	// get cart by user id
	cart := c.getCartByUserID(inventoryRequest.userID)

	// takes transaction from transaction register and applies it to the cart
	log.Printf("[debug] applying inventoryRequest %s on cart %s", inventoryRequest.requestID, cart.UserID)
	inventoryRequest, ok := cart.pendingRequests[requestID]; ok {
		productEntry, exist := cart.productEntries[inventoryRequest.Product.Code]; exist {
			productEntry.Quantity = productEntry.Quantity + inventoryRequest.Quantity 
		} else {
			productEntry = entity.ProductEntry{
				Product: inventoryRequest.Product,
				Quantity: inventoryRequest.Quantity,
				Hold: true,
			}
		}
		cart.productEntries[productEntry.Product.Code] = productEntry
		delete(requestID, c.productEntries)
		return nil

	} else {
		log.Printf("[warning] could not find request %s in cart %s", inventoryRequest.requestID, cart.UserID)
		return fmt.Errorf("could not find request %s in cart %s", inventoryRequest.requestID, cart.userID)
	}
}