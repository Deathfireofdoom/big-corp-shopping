package cart_service

import (
	"cart-service/internal/entity"
	"log"
)

// getCartHandleByUserID gets the cartHandle that belongs to the UserID, if cartHandle does not exist
// it will be created and returned.
func (c *CartService) getCartHandleByUserID(UserID string) (*entity.CartHandle, bool, error) {
	cartHandle, exist := c.carts[UserID]
	if exist {
		return cartHandle, false, nil
	} else {
		// If cart does not exists, it will be created.
		log.Printf("[debug] cart does not exist, creating.")
		cartHandle = entity.NewCartHandle(UserID)
		c.carts[UserID] = cartHandle
		return cartHandle, true, nil 
	}
}

func (c *CartService) getCartByUserID(UserID string) (*entity.Cart) {
	cartHandle, _, _ := c.getCartHandleByUserID(UserID)
	return &cartHandle.Cart
}

// updateCartHandle updates cart handle
func (c *CartService) updateCartHandle(cartHandle *entity.CartHandle) {
	c.carts[cartHandle.UserID] = cartHandle
}


func (c *CartService) InventoryRequestResponseToCartRequestResponse(inventoryRequestResponse entity.InventoryRequest) {
	// converts inventoryRequest to a cartRequestResponse
	panic("implement this")
}