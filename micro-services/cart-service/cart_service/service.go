package cart_service

import (
	"cart-service/entity"
	"cart-service/cart"
	"context"
	
	"time"
	"sync/atomic"
	"sync"
	"log"
	"fmt"
	"math"
	"github.com/google/uuid"
	"context"
)

var service = &CartService{
	carts: make(map[string]*entity.CartHandle)
}

type CartService struct {
	carts map[string]*entity.CartHandle
}

func (c *CartService) Run(ctx context.Context) error {
	log.Prtinf("[debug] initilizing the cart service...")
	// start the cart-request-communicator
	log.Printf("[debug] starting the cart-request-communicator")
	cartRequestInChannel := make(chan kafka.Message)
	cartRequestOutChannel := make(chan kafka.Message)
	cartRequestCommunicator := kafka_service.NewKafkaCommunicator("cart-request", config.KafkaCartRequest.TopicOut, config.KafkaCartRequest.TopicIn, config.KafkaCartRequst.ParitionIn, cartRequestInChannel, cartRequestOutChannel)
	go cartRequestCommunicator.Run()


	// start the inventory-request-communicator
	log.Printf("[debug] starting the inventory-request-communicator")
	inventoryRequestInChannel := make(chan kafka.Message)
	inventoryRequestOutChannel := make(chan kafka.Message)
	inventoryRequestCommunicator := kafka_service.NewKafkaCommunicator("inventory-request", config.KafkaInventoryRequest.TopicOut, config.KafkaInventoryRequest.TopicIn, config.KafkanventoryRequst.ParitionIn, inventoryRequestInChannel, inventoryRequestOutChannel)
	go inventoryRequestCommunicator.Run()

	// start listening on all channels
	log.Printf("[debug] starting the cart service...")
	for {
		select{
			case <-ctx.Done():
				log.Prtinf("[debug] got canceling signal from ctx, exiting cart-service")
				return

			case msg :=<- cartRequestOutChannel:
				log.Printf("[debug] recevied message from cart-request")
				kafkaMsg := c.handleCartRequest(msg)
				inventoryRequestInChannel <- kafkaMsg
			
			case msg :=<- inventoryRequestOutChannel:
				log.Printf("[debug] recevied message from inventory-request")

		}
	}
}



