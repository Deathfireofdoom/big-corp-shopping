type InventoryRequest struct {
	requestID 		string`json:"request_id"`
	UserID 			string`json:"user_id"`
	Product 		entity.Product`json:"product"`
	quantity 		int`json:"quantity"`
	Action			InventoryRequestAction `json:"action"`
}

type InventoryRequestResponse struct {
	statusCode 	int`json:"status_code"`
	message 	string`json:"message"`
	request		InventoryRequest `json:"request"`
}
 
type InventoryRequestAction string
const (
	InventoryRequestHold 		InventoryRequestAction = "hold"
	InventoryRequestRelease 	InventoryRequestAction = "release"
	InventoryRequestFinalize 	InventoryRequestAction = "finalize"
)