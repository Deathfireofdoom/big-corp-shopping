package entity

type Product struct {
	Name 		string `json:"product_name"`
	Code 		string `json:"product_code"`
}

type ProductEntry struct {
	Product Product
	Quantity int 
}
