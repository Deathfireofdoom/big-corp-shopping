// This handler should consume messages from kafka bus to check for incoming requests to add to
// cart, also, keeping track of responses on messages that has been sent to inventory request. 

import (
	"sync"
	"context"
	"encoding/json"
	"time"
)

var ResponseHandler = NewKafkaHandler()

type Handler struct {
	RequestEntries map[string]chan entity.InventoryRequestResponse
	mu sync.Mutex
	consumer *kafka.Reader
	producer *kafka.Writer
}


func NewKafkaHandler() *Handler {
	// create a producer for publishing messages to kafka
	producer, err := NewKafkaProducer(config.KafkaProducerConfig.Brokers, config.KafkaProducerConfig.Topic)
	if err != nil {
		panic(err)
	}

	// creates a consumer for consuming messages from kafka
	consumer := NewKafkaConsumer(config.KafkaConsumerConfig.Topic, config.KafkaConsumerConfig.Partition, config.KafkaConsumerConfig.Brokers)

	return &Handler{
		RequestEntries: make(map[string]chan entity.InventoryRequestResponse), 
		mu: sync.Mutex{}, 
		producer: producer,
		consumer: consumer,
	}
}


func (h *Handler) ListenForResponses(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Listen for canceling of context and returning.
			return
		default:
			// Tries to consume message in kafka bus. 
			msg, err := h.consumer.ReadMessage(ctx)
			if err != nil {
				log.Printf("[warning] could not consume message: %v", err)
				continue 
			}

			// Deserialize the kafka message to CartRequestResponse
			var inventoryRequestResponse entity.InventoryRequestResponse
			err = json.Unmarshal(msg.Value, &inventoryRequestResponse)
			if err != nil {
				log.Printf("[warning] could not deserialize message from kafka bus")
				continue 
			} 

			// Gets the destination channel for the message consumed by 
			h.mu.Lock()
			if responseChannel, ok := h.RequestEntries[inventoryRequestResponse.GetRequestID()]; ok {
				// publish the response to the corresponding channel
                responseChannel <- inventoryRequestResponse

				// Removes entry
				delete(h.RequestEntries, inventoryRequestResponse.GetRequestID()) 
            }
			h.mu.Unlock()
		}
	}
}

func (h *Handler) SendInventoryRequest(request entity.InventoryRequest) (chan entity.InventoryRequestResponse, error) {
	// converts struct to kafka msg
	kafkaMsg, err := ToKafkaMessage(config.KafkaProducerConfig.Topic, request)
	if err != nil {
		log.Printf("[warning] could not serialize message from kafka bus")
		return nil, err 
	}

	// creates channel that wil be used to respond 
	responseChannel := make(chan entity.InventoryRequestResponse, 1)

	// add entry to RequestEntries to keep track of requests
	handler.mu.Lock()
	handler.RequestEntries[request.RequestID] = responseChannel
	handle.mu.Unlock()

	// publish messages with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	handler.producer.WriteMessages(ctx, *kafkaMsg)

	return ch, nil
}
