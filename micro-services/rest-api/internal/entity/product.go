package entity

type Product struct {
	Name 		string `json:"product_name"`
	Code 		string `json:"product_code"`
}


type ProductPayload struct {
	Product 	Product `json:"product"`
	Quantity 	int 	`json:"quantity"`
}


type ProductEntry struct {
	Product 	Product	`json:"product"`
	Quantity 	int		`json:"quantity"`
	Hold		bool	`json:"hold"`
}

