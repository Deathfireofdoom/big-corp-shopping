package cart_service 

import(
	"encoding/json"
	"log"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"cart-service/internal/entity"
	"cart-service/internal/utils"
)

func (c *CartService) handleInventoryResponse(msg kafka.Message) (*kafka.Message, error) {
	// convert kafka message to struct
	var inventoryRequestResponse entity.InventoryRequestResponse
	err := json.Unmarshal(msg.Value, &inventoryRequestResponse)
	if err != nil {
		log.Printf("[warning] could not unmarshal inventory-request-response")
		return &kafka.Message{}, err
	}

	// check if response was successful
	if inventoryRequestResponse.StatusCode != 200 {
		log.Printf("[debug] request %s was not successful due to %s", inventoryRequestResponse.Request.RequestID, inventoryRequestResponse.Message)
	} else {
		switch inventoryRequestResponse.Request.Action {
		case entity.InventoryRequestHold:
			log.Printf("[debug] got sucessful response from inventory-service for request %s for user %s", inventoryRequestResponse.Request.RequestID, inventoryRequestResponse.Request.UserID)
			err := c.applyInventoryRequest(inventoryRequestResponse)
			if err != nil {
				log.Printf("[warning] something whent wrong applying inventory request %s", err)
				inventoryRequestResponse.StatusCode = 501
				inventoryRequestResponse.Message = fmt.Sprintf("%s", err)
			}

		case entity.InventoryRequestRelease:
			log.Printf("[info] got a response from a release request, nothing to do.")
			return &kafka.Message{}, fmt.Errorf("release request, not an error")
		
		default:
			log.Printf("[warning] unknown action %s in inventory request response %s", inventoryRequestResponse.Request.Action, inventoryRequestResponse.Request.RequestID)
			inventoryRequestResponse.StatusCode = 501
			inventoryRequestResponse.Message = fmt.Sprintf("%s", err)
		}
	}

	// gets current cart
	cart := c.getCartByUserID(inventoryRequestResponse.Request.UserID)

	// convert inventoryRequestResponse to cartRequestResponse
	cartRequestResponse := entity.NewCartRequestResponseFromInventoryResponse(inventoryRequestResponse, *cart)
	
	// serialize as kafka and return
	kafkaMsg, err := utils.ToKafkaMessage(cartRequestResponse)
	if err != nil {
		log.Printf("[warning] could not convert struct to kafka: %s", err)
		return &kafka.Message{}, err
	}

	// returns kafka message
	return kafkaMsg, nil 
}

// applyInventoryRequest applies inventory request and sends back kafka-message if the new quantity is above 
// 0. The logic behind this is that the user will not be waiting for realesing-holds, only applying holds.
func (c *CartService) applyInventoryRequest(inventoryRequestResponse entity.InventoryRequestResponse) (error) {
	// extracts original inventoryRequest
	inventoryRequest := inventoryRequestResponse.Request
	
	// get cart by user id
	cartHandle, _, err := c.getCartHandleByUserID(inventoryRequest.UserID)
	if err != nil {
		return err
	}

	// locking cartHandle to avoid races, broom broom
	log.Printf("[info] waiting for lock on cart %s", cartHandle.Cart.UserID)
	cartHandle.Mu.Lock()
	defer func () {
		cartHandle.Mu.Unlock()
		log.Printf("[info] released lock on cart %s", cartHandle.Cart.UserID)
	}()
	log.Printf("[info] obtained lock on cart %s, checking out cart", cartHandle.Cart.UserID)
	cart := cartHandle.Cart

	log.Printf("[debug] applying inventoryRequest %s on cart %s", inventoryRequest.RequestID, cart.UserID)
	// checks if this request is registered on the cart, just to make sure we don't apply on wrong cart
	_, ok := cart.PendingRequests[inventoryRequest.RequestID]
	if !ok {
		log.Printf("[warning] could not find request %s in cart %s", inventoryRequest.RequestID, cart.UserID)
		return fmt.Errorf("could not find request %s in cart %s", inventoryRequest.RequestID, cart.UserID)
	}

	// checks if there already is a product entry for the product, else creating one
	log.Printf("[debug] fetching product-entry")
	productEntry, ok := cart.ProductEntries[inventoryRequest.Product.Code]
	if ok {
		// updating existing
		productEntry.Quantity = productEntry.Quantity + inventoryRequest.Quantity
	} else {
		// creating new
		productEntry = entity.ProductEntry{
			Product: inventoryRequest.Product,
			Quantity: inventoryRequest.Quantity,
		}
	}

	// checks if new quantity is zero or below, if so, removing entry
	log.Printf("[debug] updating product entry")
	if productEntry.Quantity <= 0 {
		log.Printf("[info] new quantity is zero or below, removing entry")
		if !ok {
			log.Printf("[warning] request tries to delete products that does not have a entry in cart, requestID %s", inventoryRequest.RequestID)
		}
		delete(cart.ProductEntries, productEntry.Product.Code, ) 
	} else {
		// adding / updating the product entry
		cart.ProductEntries[productEntry.Product.Code] = productEntry
	}

	// updates last-acitvity # TODO put this logic on the cart object, not here
	log.Printf("[debug] updating cart")
	cart.LastActivity = time.Now()

	// updating cart handle with new cart
	cartHandle.Cart = cart

	// updating the cart-service with the new updated cart
	log.Printf("[debug] updating cart handle")
	c.updateCartHandle(cartHandle)
	return nil
}
