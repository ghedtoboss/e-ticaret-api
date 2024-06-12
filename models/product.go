package models

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	SellerID    int     `json:"seller_id"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}
