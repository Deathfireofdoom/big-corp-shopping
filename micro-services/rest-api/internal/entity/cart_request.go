package entity

type CartRequest struct {
	RequestID 	string 				`json:"request_id"`
	UserID 		string 				`json:"user_id"`
	Action		CartRequestAction	`json:"action"` 
	Product 	Product 			`json:"product"`
	Quantity	int					`json:"quantity"`
}

func NewCartRequest(userID string, action CartRequestAction, product Product, requestID string, quantity int) *CartRequest {
	return &CartRequest{
		RequestID: requestID,
		Action: action,
		UserID: userID,
		Product: product,
		Quantity: quantity,
	}
}

func(cr *CartRequest) GetRequestID() string {
	return cr.RequestID
}

type CartRequestAction string
const (
	CartRequestAdd 		CartRequestAction = "add"
	CartRequestOrder 	CartRequestAction = "order"
	CartRequestCheck	CartRequestAction = "check"
	CartRequestDelete	CartRequestAction = "delete"
)