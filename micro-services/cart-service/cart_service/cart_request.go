package cart_service


func (c *CartService) handleCartRequest(msg kafka.Message) kafka.Message {
	// unmarshal kafka-message to cartRequest struct
	var cartRequest entity.CartRequest
	err := json.Unmarshal(msg.Value, &cartRequest)
	if err != nil {
		log.Printf("[warning] could not unmarshal cart-request")
		return nil 
	}

	// case statement to check what kind of action is requested
	switch cartRequest.Action {
		case entity.Add:
			log.Printf("[debug] got request to add product to cart with id %s", cartRequest.UserID)
			kafkaMsg := c.addProduct(cartRequest)
			return kafkaMsg
		case entity.Remove:
			panic("Implement this")
		default:
			log.Printf("[warning] unknown action %s in cart request %s", cartRequest.Action, cartRequest.requestID)
			return nil 
	}
}

// returns nil if no communication is needed. 
func (c *CartService) addProduct(cartRequest entity.CartRequest) kafka.Message {
	// gets the cart for the user, this function will create a cart if no exist for the user
	cart, _, err := c.getCartByUserID(cartRequest.UserID)
	if err != nil {
		log.Printf("[error] could handle request %v", cartRequest)
		return nil 
	}

	// check if cart is active, meaning that all products are on hold, if not, making sure it is on hold.
	// this would happen if the user have been inactive for a while.
	cart.checkStatus()

	// creates inventoryRequest and cartTransaction
	inventoryRequest := NewInventoryRequest(entity.CartRequest)

	// add inventory request to cart as pending request
	cart.pendingRequests[inventoryRequest.requestID] = inventoryRequest

	// makes it into a kafka message and return it
	kafkaMsg, err := c.structToKafkaMsg(inventoryRequest)
	if err != nil {
		log.Printf("[warning] could not convert inventoryRequest to kafka message, id %s", cartRequest.requestID)
		return nil 
	}

	return kafkaMsg	
}

func (c *CartService) removeProduct(cartRequest entity.CartRequest) kafka.Message {
	panic("implement this")
}

func (c *CartService) deleteCart(cartRequest entity.CartRequest) []kafka.Message {
	panic("implement this")
}
