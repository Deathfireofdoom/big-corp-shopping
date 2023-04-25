package cart_service

import (
	"encoding/json"
	"log"
	"fmt"

	"github.com/segmentio/kafka-go"

	"cart-service/internal/entity"
	"cart-service/internal/utils"
)

func (c *CartService) handleCartRequest(msg kafka.Message) (*kafka.Message, bool, error) {
	// unmarshal kafka-message to cartRequest struct
	var cartRequest entity.CartRequest
	err := json.Unmarshal(msg.Value, &cartRequest)
	if err != nil {
		log.Printf("[warning] could not unmarshal cart-request")
		return &kafka.Message{}, false, err 
	}

	// case statement to check what kind of action is requested
	switch cartRequest.Action {
		case entity.CartRequestAdd:
			log.Printf("[debug] got request to add product to cart with id %s", cartRequest.UserID)
			kafkaMsg, _, err := c.addPendingProductEntry(cartRequest)
			return kafkaMsg, false, err
			
		case entity.CartRequestDelete:
			log.Printf("[debug] got request to remove product to cart with id %s", cartRequest.UserID)
			kafkaMsg, inventoryRequest, err := c.addPendingProductEntry(cartRequest)
			err = c.applyInventoryRequest(inventoryRequest)
			if err != nil {
				log.Printf("[warning] could not apply inventory request on delete: %s", err)
				return nil, false, fmt.Errorf("not implemented")
			}

			return kafkaMsg, false, nil
		
		case entity.CartRequestCheck:
			log.Printf("[debug] got request to check cart")
			kafkaMsg, err := c.handleCartRequestCheck(cartRequest)
			return kafkaMsg, true, err


		default:
			log.Printf("[warning] unknown action %s in cart request %s", cartRequest.Action, cartRequest.RequestID)
			return &kafka.Message{}, false, fmt.Errorf("not implemented")  
	}
}

func (c *CartService) handleCartRequestCheck(cartRequest entity.CartRequest) (*kafka.Message, error) {
	// get cart by userID
	cart := c.getCartByUserID(cartRequest.UserID)

	// convert cartRequest to cartRequest response
	cartRequestResponse := entity.NewCartRequestResponseFromCartRequest(cartRequest, *cart)

	// converts to kafka-message
	kafkaMsg, err := utils.ToKafkaMessage(cartRequestResponse)
	if err != nil {
		log.Printf("[warning] could not convert cartRequestResponse to kafka message, id %s", cartRequest.RequestID)
		return &kafka.Message{}, err
	}
	return kafkaMsg, nil
}


// returns nil if no communication is needed. 
func (c *CartService) addPendingProductEntry(cartRequest entity.CartRequest) (*kafka.Message, entity.InventoryRequest, error) {
	// gets the cart-handle, and creates a new one if it does not exsists
	cartHandle, _, err := c.getCartHandleByUserID(cartRequest.UserID)
	if err != nil {
		log.Printf("[error] could handle request %v", cartRequest)
		return &kafka.Message{}, entity.InventoryRequest{}, err 
	}

	// checks out lock to avoid races
	log.Printf("[info] waiting for lock on cartHandle")
	cartHandle.Mu.Lock()
	defer func () {
		log.Printf("[info] releasing lock on cartHandle")
		cartHandle.Mu.Unlock()
	}()

	// implement a check, to make sure the cart is on hold.
	//cartHandle.cart.checkStatus()
	log.Printf("[todo] implement checkstatus on cart object")

	// converts cartRequest to inventoryRequest
	inventoryRequest := entity.NewInventoryRequestFromCartRequest(cartRequest)

	// add request to pending.
	cartHandle.Cart.PendingRequests[inventoryRequest.RequestID] = inventoryRequest

	// makes it into a kafka message and return it
	kafkaMsg, err := utils.ToKafkaMessage(inventoryRequest)
	if err != nil {
		log.Printf("[warning] could not convert inventoryRequest to kafka message, id %s", cartRequest.RequestID)
		return &kafka.Message{}, entity.InventoryRequest{}, err 
	}
	return kafkaMsg, inventoryRequest, nil
}
