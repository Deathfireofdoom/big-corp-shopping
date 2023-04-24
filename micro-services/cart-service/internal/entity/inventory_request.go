package entity

import (
	"fmt"
)

type InventoryRequest struct {
	RequestID 		string`json:"request_id"`
	UserID 			string`json:"user_id"`
	Product 		Product`json:"product"`
	Quantity 		int`json:"quantity"`
	Action			InventoryRequestAction `json:"action"`
}

type InventoryRequestResponse struct {
	StatusCode 	int`json:"status_code"`
	Message 	string`json:"message"`
	Request		InventoryRequest `json:"request"`
}
 
type InventoryRequestAction string
const (
	InventoryRequestHold 		InventoryRequestAction = "hold"
	InventoryRequestRelease 	InventoryRequestAction = "release"
	InventoryRequestFinalize 	InventoryRequestAction = "finalize"
)

func NewInventoryRequestFromCartRequest(cartRequest CartRequest) InventoryRequest {
	var inventoryRequestAction InventoryRequestAction
	switch cartRequest.Action {
		case CartRequestAdd:
			inventoryRequestAction = InventoryRequestHold
		case CartRequestDelete:
			inventoryRequestAction = InventoryRequestRelease
		case CartRequestOrder:
			panic("implement order this")
	default:
		panic(fmt.Sprintf("unknown cartRequest.Action %s", cartRequest.Action)) // todo handle this
	}
	return InventoryRequest{
		RequestID: cartRequest.RequestID,
		UserID: cartRequest.UserID,
		Product: cartRequest.Product,
		Quantity: cartRequest.Quantity,
		Action: inventoryRequestAction,
	}
}