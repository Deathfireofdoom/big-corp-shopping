package inventory_service

import (
	"context"
	"log"
	"encoding/json"

	"github.com/segmentio/kafka-go"


	// internal
	"inventory/internal/entity"
	"inventory/internal/kafka_service"
	"inventory/internal/config"
	"inventory/internal/utils"
)

// todo 21 april
// fill database with data
// rewrite the run-function in this file
// star the service in main.go

type InventoryService struct {
	workerPool 				map[string]*InventoryWorker
	workerRequestChannel 	chan entity.InventoryRequest
	workerResponseChannel 	chan entity.InventoryRequestResponse
}

func NewInventoryService() *InventoryService{
	workerRequestChannel := make(chan entity.InventoryRequest, 10)
	workerResponseChannel := make(chan entity.InventoryRequestResponse, 10)


	workerPool := map[string]*InventoryWorker{}
	for i:=0; i < 5; i++ {
		id, _ := utils.GetUniqueID()
		workerPool[id] = NewInventoryWorker(id, workerRequestChannel, workerResponseChannel)
	}

	return &InventoryService{
		workerPool: workerPool,
		workerRequestChannel: workerRequestChannel,
		workerResponseChannel: workerResponseChannel, 
	}
}


func (is *InventoryService) Run(ctx context.Context) error {
	log.Printf("[debug] initilizing the inventory-service...")
	// start the cart-request-communicator
	log.Printf("[debug] starting the inventory-request-communicator")
	inventoryRequestInChannel := make(chan kafka.Message, 10)
	inventoryRequestOutChannel := make(chan kafka.Message, 10)
	inventoryRequestCommunicator, err := kafka_service.NewKafkaCommunicator("inventory-request", config.KafkaConfig.RequestTopic, config.KafkaConfig.ResponseTopic,  0, inventoryRequestInChannel, inventoryRequestOutChannel, config.KafkaConfig.Brokers)
	if err != nil {
		log.Printf("[error] could not start inventory-request-communicator")
		return err
	}
	go inventoryRequestCommunicator.Run(ctx)


	// starting workers
	for _, worker := range is.workerPool {
		go worker.Run(ctx)
	}

	// start listening on all channels
	log.Printf("[debug] starting the inventory service...")
	for {
		select{
			case <-ctx.Done():
				log.Printf("[debug] got canceling signal from ctx, exiting cart-service")
				return nil

			case msg :=<- is.workerResponseChannel:
				log.Printf("[debug] recieved response from worker")
				kafkaMsg, err := utils.ToKafkaMessage(msg)
				if err != nil {
					log.Printf("[warning] could not convert struct to kafka message")
					continue
				}
				inventoryRequestInChannel <- *kafkaMsg
			
			case msg :=<- inventoryRequestOutChannel:
				log.Printf("[debug]  recevied inventory-request")
				// unmarshal and send to next channel
				var inventoryRequest entity.InventoryRequest
				err := json.Unmarshal(msg.Value, &inventoryRequest)
				if err != nil {
					log.Printf("[error] failed to unmarshal inventory request: %s", err)
					continue
				}
				is.workerRequestChannel <- inventoryRequest
		}
	}
}
