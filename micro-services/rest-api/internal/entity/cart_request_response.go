package entity 

type CartRequestResponse struct {
	RequestID string 	`json:"request_id"`
	StatusCode int 		`json:"status_code"`
	Message string 		`json:"message"`
	Cart	Cart		`json:"cart"`
}

func (response *CartRequestResponse) GetRequestID() string {
    return response.RequestID
}

func (response *CartRequestResponse) GetStatusCode() int {
    return response.StatusCode
}

func (response *CartRequestResponse) GetMessage() string {
    return response.Message
}