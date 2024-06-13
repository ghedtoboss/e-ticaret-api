package models

// Product represents a product in the system.
// @Description Ürün modelini temsil eder
type Product struct {
	ID          int     `json:"id" example:"1"`                    
	Name        string  `json:"name" example:"Product Name"`       
	Description string  `json:"description" example:"Description"` 
	Quantity    int     `json:"quantity" example:"100"`            
	Price       float64 `json:"price" example:"19.99"`             
	SellerID    int     `json:"seller_id" example:"1"`             
	Category    string  `json:"category" example:"Electronics"`    
	ImageURL    string  `json:"image_url" example:"http://..."`    
}
