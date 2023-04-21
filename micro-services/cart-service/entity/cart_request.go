

type CartRequest struct{
	requestID 		string`json:"request_id"`
	UserID 			string`json:"user_id"`
	Product 		entity.Product`json:"product"`
	quantity 		int`json:"quantity"`
	Action			CartRequestResponse `json:"action"`
}


type CartRequestResponse struct {
	requestID 	string`json:"request_id"`
	statusCode 	int`json:"status_code"`
	message 	string`json:"message"`
	cart		cart.Cart`json:"cart"`
}

type CartRequestAction string
const (
	CartRequestUpdate 		CartRequestAction = "update"
	CartRequestOrder 		CartRequestAction = "order"
	CartRequestCheck	 	CartRequestAction = "check"
	CartRequestDelete	 	CartRequestAction = "delete"
)