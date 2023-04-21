package entity

type CartRequest struct {
	requestID 	string 		`json:"request_id"`
	userID 		string 		`json:"user_id"`
	action		Action		`json:"action"` 
	product 	Product 	`json:"product"`
	quantity	int			`json:"quantity"`
}

func NewCartRequest(userID string, action Action, product Product, requestID string) *CartRequest {
	return &CartRequest{
		requestID: requestID,
		action: action,
		userID: userID,
		product: product,
	}
}

func(cr *CartRequest) GetRequestID() string {
	return cr.requestID
}

type Action string
const (
	Add 	Action = "add"
	Delete 	Action = "delete"
	Check 	Action = "check"
)