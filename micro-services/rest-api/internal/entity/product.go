package entity

type Product struct {
	Name 		string `json:"product_name"`
	Code 		string `json:"product_code"`
}


type ProductEntry struct {
	product 	Product	`json:"product"`
	quantity 	int		`json:"quantity"`
	hold		bool	`json:"hold"`
}


type ProductPayload struct {
	Product 	Product `json:"product"`
	Quantity 	int 	`json:"quantity"`
}