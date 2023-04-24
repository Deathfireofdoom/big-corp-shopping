package cart_service

import (
	"encoding/json"
	"log"
	"fmt"

	"github.com/segmentio/kafka-go"

	"cart-service/internal/entity"
	"cart-service/internal/utils"
)

func (c *CartService) handleCartRequest(msg kafka.Message) (*kafka.Message, error) {
	// unmarshal kafka-message to cartRequest struct
	var cartRequest entity.CartRequest
	err := json.Unmarshal(msg.Value, &cartRequest)
	if err != nil {
		log.Printf("[warning] could not unmarshal cart-request")
		return &kafka.Message{}, err 
	}

	// case statement to check what kind of action is requested
	switch cartRequest.Action {
		case entity.CartRequestAdd:
			log.Printf("[debug] got request to add product to cart with id %s", cartRequest.UserID)
			kafkaMsg, err := c.addProduct(cartRequest)
			return kafkaMsg, err
			
		case entity.CartRequestDelete:
			log.Printf("[debug] got request to add product to cart with id %s", cartRequest.UserID)
			log.Printf("[TODO] implement cartRequest.Delete")
			return &kafka.Message{}, fmt.Errorf("not implemented") 

		default:
			log.Printf("[warning] unknown action %s in cart request %s", cartRequest.Action, cartRequest.RequestID)
			return &kafka.Message{}, fmt.Errorf("not implemented")  
	}
}

// returns nil if no communication is needed. 
func (c *CartService) addProduct(cartRequest entity.CartRequest) (*kafka.Message, error) {
	// gets the cart-handle, and creates a new one if it does not exsists
	cartHandle, _, err := c.getCartHandleByUserID(cartRequest.UserID)
	if err != nil {
		log.Printf("[error] could handle request %v", cartRequest)
		return &kafka.Message{}, err 
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
		return &kafka.Message{} ,err 
	}
	return kafkaMsg, nil
}


func (c *CartService) removeProduct(cartRequest entity.CartRequest) kafka.Message {
	panic("implement this")
}

func (c *CartService) deleteCart(cartRequest entity.CartRequest) []kafka.Message {
	panic("implement this")
}
