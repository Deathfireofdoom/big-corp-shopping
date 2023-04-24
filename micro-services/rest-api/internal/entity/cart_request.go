package entity

type CartRequest struct {
	RequestID 	string 		`json:"request_id"`
	UserID 		string 		`json:"user_id"`
	Action		Action		`json:"action"` 
	Product 	Product 	`json:"product"`
	Quantity	int			`json:"quantity"`
}

func NewCartRequest(userID string, action Action, product Product, requestID string, quantity int) *CartRequest {
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

type Action string
const (
	Add 	Action = "add"
	Delete 	Action = "delete"
	Check 	Action = "check"
)