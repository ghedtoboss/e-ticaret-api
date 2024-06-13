package models

// Cart represents a shopping cart.
// @Description Alışveriş sepetini temsil eder
type Cart struct {
	ID     int `json:"id" example:"1"`      
	UserID int `json:"user_id" example:"1"` 
}

// CartItem represents an item in the shopping cart.
// @Description Sepet öğesi modelini temsil eder
type CartItem struct {
	ID        int     `json:"id" example:"1"`         
	CartID    int     `json:"cart_id" example:"1"`    
	ProductID int     `json:"product_id" example:"1"` 
	Quantity  int     `json:"quantity" example:"1"`   
	Price     float64 `json:"price" example:"19.99"`  
}
