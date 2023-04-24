package cart_service

import (
	"github.com/segmentio/kafka-go"
	"cart-service/internal/entity"
	"cart-service/internal/config"
	"cart-service/internal/kafka_service"
	

	"log"
	"context"
	"time"
)

var Service = &CartService{
	carts: make(map[string]*entity.CartHandle),
}

type CartService struct {
	carts map[string]*entity.CartHandle
}

func (c *CartService) Run(ctx context.Context) error {
	log.Printf("[debug] initilizing the cart service...")
	// start the cart-request-communicator
	log.Printf("[debug] starting the cart-request-communicator")
	cartRequestInChannel := make(chan kafka.Message, 10)
	cartRequestOutChannel := make(chan kafka.Message, 10)
	cartRequestCommunicator, err := kafka_service.NewKafkaCommunicator("cart-request", config.KafkaCartRequest.RequestTopic, config.KafkaCartRequest.ResponseTopic, 0, cartRequestInChannel, cartRequestOutChannel, config.KafkaCartRequest.Brokers)
	if err != nil {
		log.Println("[error] could not start cart-request-communicator")
		panic(err)
	}
	go cartRequestCommunicator.Run(ctx)


	// start the inventory-request-communicator
	log.Printf("[debug] starting the inventory-request-communicator")
	inventoryRequestInChannel := make(chan kafka.Message, 10)
	inventoryRequestOutChannel := make(chan kafka.Message, 10)
	inventoryRequestCommunicator, err := kafka_service.NewKafkaCommunicator("inventory-request", config.KafkaInventoryRequest.ResponseTopic, config.KafkaInventoryRequest.RequestTopic, 0, inventoryRequestInChannel, inventoryRequestOutChannel, config.KafkaInventoryRequest.Brokers)
	if err != nil {
		log.Println("[error] could not start inventory-request-communicator")
		panic(err)
	}
	go inventoryRequestCommunicator.Run(ctx)


	// ticker channel to get some logging, probably not useful in prod, but nice to see what is happening
	ticker := time.NewTicker(10 * time.Second)
	tickerChannel := ticker.C

	// start listening on all channels
	log.Printf("[debug] starting the cart service...")
	for {
		select{
			case <-ctx.Done():
				log.Printf("[debug] got canceling signal from ctx, exiting cart-service")
				return nil
			
			case <-tickerChannel:
				log.Println("[debug] system info:")
				log.Printf("%s", c)

			case msg :=<- cartRequestOutChannel:
				log.Printf("[debug] recevied message from cart-request")
				kafkaMsg, err := c.handleCartRequest(msg)
				if err != nil {
					log.Printf("[error] %s", err)
					continue
				}
				inventoryRequestInChannel <- *kafkaMsg
			
			case msg :=<- inventoryRequestOutChannel:
				log.Printf("[debug] recevied message from inventory-request")
				kafkaMsg, err := c.handleInventoryResponse(msg)
				if err != nil {
					log.Printf("[error] %s", err)
					continue
				}
				cartRequestInChannel <- *kafkaMsg
		}
	}
}



