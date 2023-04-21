package entity


// TODO: This should be moved to another package I think.

type Status int

const (
	Success = iota
	Fail
)


type InventoryRequest struct {
	Products 	[]Product `json:"products"`
	Action 		string `json:"action"`
	RequestID 	string `json:"request_id"`
}

type Product struct {
	Name 		string `json:"product_name"`
	Code 		string `json:"product_code"`
	Quantity	int `json:"quantity"`
}

type InventoryResult struct {
	Status Status
	Request InventoryRequest
	RequestID string 
}
