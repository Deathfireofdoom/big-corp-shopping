package cart 

import (
	"cart-service/entity"
	"time"
	"sync/atomic"
	"sync"
	"log"
	"fmt"
	"math"
	"github.com/google/uuid"
	"context"
)


var counter uint64 = 0

type Cart struct {
	id uint64
	userID string
	productEntries map[string]entity.ProductEntry
	lastActivity time.Time
	mu sync.Mutex 
	pendingRequests map[string]InventoryRequest
}

func (c *Cart) checkStatus() {
	// checks if all products are on hold, this should be a blocking operation. Maybe with a channel or something.
	log.Printf("[warning] TODO implement check status function.")
}

func NewCart(UserID string) *Cart {
	return &Cart{
		id: atomi.AddUint64(&counter, 1)
		userID: UserID,
		productEntries: map[string]entity.ProductEntry{}
		lastActivity: time.Now(),
		mu: sync.Mutex{},
		pendingRequests: map[string]InventoryRequest{}
	}
}




// 
func NewInventoryRequest(cartRequest entity.CartRequest) InventoryRequest {
	// generate unique id, generate a inventory request and return it
	requestID := getUniqueRequestID()
	return &InventoryRequest{
		requestID: requestID,
		productCode: cartRequest.ProductCode,
		quantity: cartRequest.Quantity,
		UserID: cartRequest.UserID,
		Action: cartRequest.Action,
	}
} 






// OoooooLLLLDDD



func (c *CartService) DeleteCartByUserID(UserID string) error {
	// Check if userID has cart. Return error if not.
	_, exist := c.carts[UserID]
	if exist {
		log.Printf("[debug] deleting cart for UserID %s", UserID)
		delete(c.carts, UserID)
		return nil
	} else {
		log.Printf("[warning] trying to delete non-existent cart for UserID %s", UserID)
		return fmt.Errorf("no cart exist for UserID %s", UserID)
	}
}





func (c *Cart) AddProduct(product entity.Product, quantity int) error {
	// Check if cart has put hold on everything.
	log.Prtinf("TODO - implement a check so everything is on hold.")

	// Send message to make hold on product. 
	err := c.sendInventoryRequest(product, quantity, entity.Hold)
	if err {
		fmt.Printf("[warning] could not hold %v of product %s for UserID %s", quantity, product.Name, c.ownerID)
		return err
	}

	// If not make a product entry and add it to the map.
	productEntry, exist := c.productEntries[product.Code]
	if exist {
		log.Printf("[debug] product already exist in cart, adding quantity")
		productEntry.Quantity =  productEntry.Quantity + quantity
		productEntry.HoldQuantity = productEntry.HoldQuantity + quantity
	} else {
		log.Printf("[debug] product does not exist in cart, adding entry")
		productEntry = entity.ProductEntry{
			Product: product,
			Quantity: quantity,
			HoldQuantity, quantity,
		}
	}

	// Update the entry
	c.productEntries[productEntry.Product.Code] = productEntry
	return nil
}

func (c *Cart) sendInventoryRequest(product entity.Product, quantity int, entity.Action) (error) {	
	// converts input to inventory request
	// generates unique id.
	requestID, err := getUniqueRequestID()
	if err != nil {
		log.Printf("[warning] could not generate uniqueID: %v", err)
		return err 
	}

	// creates struct that will be sent to kafka-handler.
	inventoryRequest := &InventoryRequest{
		Action: action,
		CartID: c.GetCartID(), 
		ProductCode: product.Code, 
		Quantity: quantity
	}

	// send request to kafka-handler
	responseChannel, err := kafka.Handler.SendInventoryRequest(inventoryRequest)
	if err != nil {
		log.Printf("[warning] could not make inventory-request: %v", err)
		return err
	}

	// waiting for response
	inventoryRequest :=<-responseChannel
	if inventoryRequest.Error != nil {
		log.Printf("[warning] inventory request was denied: %v", err)
		return err
	}

	// returns nil to indicate the request was ok.
	return nil
}

func (c *Cart) RemoveProduct(product entity.Product, quantity int) error {
	// Check if product code is in map, if not return error
	productEntry, exist := c.productEntries[product.Code]
	if !exist {
		log.Printf("[warning] %s is not in cart with id %s, can't delete", product.Name, c.GetCartID())
		return fmt.Errorf("product %s does not exist in cart %s", product.Name, c.GetCartID())
	}

	// Check if number of removed products are more then product in cart, return error if
	if productEntry.Quantity < quantity {
		log.Printf("[warning] trying do delete %v of product %s but only %v exist, no delete", quantity, product.Name, productEntry.Quantity)
		return fmt.Errorf("trying do delete %v of product %s but only %v exist, no delete", quantity, product.Name, productEntry.Quantity)
	}

	// Check if removed products are less than product in cart, decrease amount.
	productEntry.Quantity = productEntry.Quantity - quantity

	// send request inventory to remove hold quantity, the user does not care about resopnse so will not wait on response.
	var releaseQuantity int
	if quantity > productEntry.HoldQuantity {
		releaseQuantity = prodcutEntry.HoldQuantity
	} else {
		releaseQuantity = quantity
	}
	go sendInventoryRequest(product, releaseQuantity, entity.Release)

	// Check if removed product are equal to product in cart, remove entry.
	if productEntry.Quantity == 0 {
		log.Printf("[debug] new amount for product %s is 0, removing entry", product.Name)
		return nil
	} else {
		log.Printf("[debug] new amount for product %s is %v, updating entry", product.Name, productEntry.Quantity)
		c.productEntries[product.Code] = productEntry
		return nil
	}
}

func (c *Cart) RemoveCart() error {
	// Remove all products from cart and send message to inventory service to remove them.
	// Return error if not successful. Maybe delete itself if possible.
	log.Printf("[debug] removing cart for UserID %s", c.ownerID)
	var err error
	for _, productEntry := range c.productEntries {
		err = c.RemoveProduct(productEntry.Product, productEntry.Quantity)
		err != nil {
			log.Printf("[warning] could not delete product %s from cart UserID %s", productEntry.Product.Name, c.ownerID)
		}
	}
	return err 
}

func (c *Cart) Buy() error {
	// Sends the cart to be bought.
	panic("Implement this.")
}

func (c *Cart) GetCartID() uint {
	return c.id
}

func (c *Cart) GetOwnerID() string {
	return c.ownerID
}

func (c *Cart) SetOwnerID(string UserID) error {
	// TODO set up checks.
	c.ownerID = UserID
	return nil
}


func getUniqueRequestID() (string, error) {
	return UUID.New(), nil 
}