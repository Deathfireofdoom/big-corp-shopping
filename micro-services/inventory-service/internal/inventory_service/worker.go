package inventory_service

import (
	"context"
	"log"
	"fmt"

	// internal
	"inventory/internal/entity"
	"inventory/internal/db"
	"inventory/internal/redis_service"

)


type InventoryWorker struct {
	id 				string
	requestChannel 	chan entity.InventoryRequest
	responseChannel	chan entity.InventoryRequestResponse
	db				*db.DbService
}

func NewInventoryWorker(id string, requestChannel chan entity.InventoryRequest, responseChannel	chan entity.InventoryRequestResponse, redisService *redis_service.RedisService) *InventoryWorker {
	db, err := db.NewDbService(redisService)
	if err != nil {
		log.Printf("[error] could not connect to database: %s", err)
		panic(err)
	}
	return &InventoryWorker{
		id: id,
		requestChannel: requestChannel,
		responseChannel: responseChannel,
		db: db,
	}
}

func (iw *InventoryWorker) Run(ctx context.Context) {
	log.Printf("[info] starting worker %s", iw.id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[info] got cancling signal on worker %s, shutting down", iw.id)
			return
		case request :=<- iw.requestChannel:
			log.Printf("[info] got request %s on worker %s", request.RequestID, iw.id)
			iw.handleRequest(request)
			log.Printf("[info] request %s handled by worker %s", request.RequestID, iw.id)
		}
	}
}


func (iw *InventoryWorker) handleRequest(request entity.InventoryRequest) {
    switch request.Action {
    case entity.InventoryRequestHold:
		log.Printf("[debug] go hold request on worker %s", iw.id)
        iw.putOnHold(request)
    
	case entity.InventoryRequestRelease:
		log.Printf("[debug] go release request on worker %s", iw.id)
		iw.putOnHold(request)
    
	case entity.InventoryRequestFinalize:
		log.Printf("[debug] go finalize request on worker %s", iw.id)
        log.Printf("[TODO] finalize is not implemented")
    
	default:
        // Handle unknown action
        log.Printf("[warning] unknown action: %v", request.Action)
    }
}
 
func (iw *InventoryWorker) putOnHold(request entity.InventoryRequest) {
	// checks db if hold can be satisfied
	_, err := iw.db.PutOnHold(request.Product.Code, request.Quantity, request.UserID)
	
	var message string
	var statusCode int
	if err != nil {
		log.Printf("[warning] could not fullfil request %s", request.RequestID)
		message = fmt.Sprintf("%v", err)
		statusCode = 500
	} else {
		log.Printf("[debug] could fullfil request %s", request.RequestID)
		message = "ok"
		statusCode = 200
	}

	// creates struct to return response to inventory-service
	inventoryRequestResponse := entity.InventoryRequestResponse{
		StatusCode: statusCode,
		Message: message,
		Request: request,
	}
	iw.responseChannel <- inventoryRequestResponse
	log.Printf("[debug] request %s handled by %s", request.RequestID, iw.id)
}

func (iw *InventoryWorker) releaseHold() {
	panic("implement this")
}

func (iw *InventoryWorker) finalizeHold() {
	panic("implement this")
}