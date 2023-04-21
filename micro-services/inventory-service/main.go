package main

import (
	"context"
	"inventory/entity"
	"inventory/inventory"
	"inventory/pubsub"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setting up channels
	inRequest := make(chan entity.InventoryRequest)

	// Connecting to Kafka-bridge
	brokers := []string{"localhost:9092"}
	const topicRequest = "inventory-request"
	const partition = 0 
	consumerKafka := pubsub.NewConsumerKafka(brokers, topicRequest, partition)

	// Starting consumer to read from Kafka-stream and publish to channel
	go consumerKafka.Consume(ctx, inRequest, false)
	
	// Starting request handler and getting output channel
	out := Inventory.Service.HandleRequest(ctx, inRequest)

	// Connecting publisher to Kafka-bridge
	const topicResult = "inventory-result"
	producerKafka := pubsub.NewProducerKafka(brokers, topicResult)

	// Starting producer
	go producerKafka.Produce(ctx, out, false)

	go sendMockRequest(inRequest)
	for {

	}	
}

func sendMockRequest(in chan entity.InventoryRequest){
	mockRequests := []entity.InventoryRequest{
		{
			Products: []entity.Product{
				{Name: "Product A", Quantity: 10},
				{Name: "Product B", Quantity: 5},
			},
			Action:    "hold",
			RequestID: "123",
		},
		{
			Products: []entity.Product{
				{Name: "Product C", Quantity: 7},
				{Name: "Product D", Quantity: 3},
			},
			Action:    "hold",
			RequestID: "456",
		},
		{
			Products: []entity.Product{
				{Name: "Product C", Quantity: 7},
				{Name: "Product D", Quantity: 3},
			},
			Action:    "hold",
			RequestID: "789",
		},
		{
			Products: []entity.Product{
				{Name: "Product C", Quantity: 7},
				{Name: "Product D", Quantity: 3},
			},
			Action:    "hold",
			RequestID: "101112",
		},
		{
			Products: []entity.Product{
				{Name: "Product C", Quantity: 7},
				{Name: "Product D", Quantity: 3},
			},
			Action:    "hold",
			RequestID: "131415",
		},
		{
			Products: []entity.Product{
				{Name: "Product C", Quantity: 7},
				{Name: "Product D", Quantity: 3},
			},
			Action:    "hold",
			RequestID: "161718",
		},
	}

	for _, request := range mockRequests {
		in <- request
		time.Sleep(3 * time.Second)
	}

}

