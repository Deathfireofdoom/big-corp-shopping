package entity 

type CartRequestResponse struct {
	requestID string 	`json:"request_id"`
	statusCode int 		`json:"status_code"`
	message string 		`json:"message"`
	cart	Cart		`json:"cart"`
}

func (response *CartRequestResponse) GetRequestID() string {
    return response.requestID
}

func (response *CartRequestResponse) GetStatusCode() int {
    return response.statusCode
}

func (response *CartRequestResponse) GetMessage() string {
    return response.message
}