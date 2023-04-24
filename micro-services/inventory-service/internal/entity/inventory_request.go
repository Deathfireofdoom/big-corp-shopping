package entity

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