package entity

type CartRequest struct{
	RequestID 		string				`json:"request_id"`
	UserID 			string				`json:"user_id"`
	Product 		Product				`json:"product"`
	Quantity 		int					`json:"quantity"`
	Action			CartRequestAction 	`json:"action"`
}


type CartRequestResponse struct {
	RequestID 	string	`json:"request_id"`
	StatusCode 	int		`json:"status_code"`
	Message 	string	`json:"message"`
	Cart		Cart	`json:"cart"`
}

type CartRequestAction string
const (
	CartRequestAdd 		CartRequestAction = "add"
	CartRequestOrder 	CartRequestAction = "order"
	CartRequestCheck	CartRequestAction = "check"
	CartRequestDelete	CartRequestAction = "delete"
)

func NewCartRequestResponseFromInventoryResponse(inventoryRequestResponse InventoryRequestResponse, cart Cart) CartRequestResponse {
	return CartRequestResponse{
		RequestID: inventoryRequestResponse.Request.RequestID,
		StatusCode: inventoryRequestResponse.StatusCode,
		Message: inventoryRequestResponse.Message,
		Cart: cart,
	}
}