package cart_service

// GetCartByUserID gets the cart that belongs to the UserID, if cart does not exist
// it will be created and returned.
func (c *CartService) getCartByUserID(UserID string) (*Cart, bool, error) {
	cart, exist := c.carts[UserID]
	if exist {
		return cart, false, nil
	} else {
		// If cart does not exists, it will be created.
		log.Printf("[debug] cart does not exist, creating.")
		cart, err := c.NewCart(UserID)
		if err {
			log.Printf("[error] could not create for user %s", UserID)
			return cart, true, err
		} else {
			return cart, true, nil
		}
	}
}

func (c *CartService) newCart(UserID string) (*Cart, error) {
	// check if userid is alredy present, return error if.
	_, exist := c.carts[UserID]
	if exist {
		log.Printf("[warning] tried to create a cart for a user that already exists")
		return &Cart{}, fmt.Errorf("cart already exists for user %s", UserID)
	}

	// Add new cart to map with UserID as key. return nil
	cart := NewCart(UserID)
	c.carts[UserID] = cart
	return cart, nil
}

func (c *CartService) structToKafkaMsg(object interface{}) (kafka.Message, error) {
	// converts struct to json.
	jsonData, err := json.Marshal(object)
	if err != nil {
		log.Printf("[warning] could not marshal object to json for kafka-msg-conversion. %v", err)
		return nil, err 
	}

	// creates and returns kafka message
	message := kafka.Message{
		Value: jsonData
	}
	return message, nil 
}

func (c *CartService) InventoryRequestResponseToCartRequestResponse(inventoryRequestResponse entity.InventoryRequest) {
	// converts inventoryRequest to a cartRequestResponse
}