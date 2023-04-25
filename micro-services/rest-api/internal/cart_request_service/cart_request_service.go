package cart_request_service

import (
	"log"
	"sync"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"

	// internal modules
	"big-corp-shopping/rest-api/internal/entity"
	"big-corp-shopping/rest-api/internal/kafka_service"
	"big-corp-shopping/rest-api/internal/config"
	"big-corp-shopping/rest-api/internal/utils"

)

var Service *CartRequestService = NewCartRequestService()

type CartRequestService struct {
	requests 		map[string]chan entity.CartRequestResponse
	requestChannel	chan kafka.Message
	responseChannel	chan kafka.Message
	mu 				sync.Mutex
}

func NewCartRequestService() *CartRequestService {
	requestChannel := make(chan kafka.Message, 10)
	responseChannel := make(chan kafka.Message, 10)

	return &CartRequestService{
		requests: make(map[string]chan entity.CartRequestResponse),
		requestChannel: requestChannel,
		responseChannel: responseChannel,
		mu: sync.Mutex{},
	}
}

// Run starts the CartRequestService, listening on incomming request (via invocation of NewCartRequest) and
// incoming responses from Kafka.
func(cs *CartRequestService) Run (ctx context.Context) {
	// setting up config that should be used by the kafka-communicator
	responseTopic := config.KafkaConfig.ResponseTopic
	requestTopic := config.KafkaConfig.RequestTopic
	kafkaBrokers := config.KafkaConfig.Brokers

	log.Printf("%s", requestTopic)
	// initializing kafka communicator
	kafkaCommunicator, err := kafka_service.NewKafkaCommunicator("cart-request-communicator", responseTopic, requestTopic, 0, cs.requestChannel, cs.responseChannel, kafkaBrokers)
	if err != nil {
		// ok to panic since this is a crucial step
		panic(err)
	}

	// starting the kafka communicator
	log.Println("[info] starting communicator")
	go func() {
		kafkaCommunicator.Run(ctx)
	}()
	
	// starting response-request
	log.Println("[info] start listening from cart-request-responses")
	cs.listenForResponse(ctx)
}


func (cs *CartRequestService) listenForResponse (ctx context.Context) {
	for {
		select {
		// for canceling
		case <-ctx.Done():
			log.Print("[info] recieved signal from ctx, shutting down cart-service")
			return
		
		// for consuming responses from kafka-communicator
		case msg := <-cs.responseChannel:
			log.Printf("[info] recieved message in response channel")
			// msg is of type *kafka.Message and needs to be converted to cartRequestResponse
			var cartRequestResponse entity.CartRequestResponse
			err := json.Unmarshal(msg.Value, &cartRequestResponse)
			if err != nil {
				log.Printf("[warning] could not deserialize message from kafka communicator")
				continue 
			} 

			cs.mu.Lock()
			if responseClientChannel, ok := cs.requests[cartRequestResponse.GetRequestID()]; ok {
				log.Printf("[info] got response of request %s", cartRequestResponse.GetRequestID())
				// publish the response to the corresponding channel
                responseClientChannel <- cartRequestResponse

				// removes entry
				delete(cs.requests, cartRequestResponse.GetRequestID()) 
            } else {
				log.Printf("[DEBUG] %s", cartRequestResponse)
				log.Printf("[warning] could not find request with request id %s", cartRequestResponse.GetRequestID())
			}
			cs.mu.Unlock()
		}

	}
}

func (cs *CartRequestService) NewCartRequest(userID string, action entity.CartRequestAction, product entity.Product, quantity int) (chan entity.CartRequestResponse, error) {
	// generate unique id 
	requestID, err := utils.GetUniqueID()
	if err != nil {
		log.Printf("[warning] could not generate unique id.")
		return nil, err 
	}

	// creates cart request object
	cartRequest := entity.NewCartRequest(userID, action, product, requestID, quantity)

	// makes request and gets response channel back
	ch, err := cs.makeRequest(*cartRequest)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (cs *CartRequestService) makeRequest(cartRequest entity.CartRequest) (chan entity.CartRequestResponse, error) {
	// makes channel for response
	ch := make(chan entity.CartRequestResponse, 1)

	// converts cartRequest to binary so it can be sent via kafka
	kafkaMsg, err := utils.ToKafkaMessage(cartRequest)
	if err != nil {
		return nil, err
	}

	// makes entry in request
	cs.mu.Lock()
	cs.requests[cartRequest.GetRequestID()] = ch
	cs.mu.Unlock()

	// publish kafka message to kafka communicator
	cs.requestChannel <- *kafkaMsg
	
	return ch, nil
}
