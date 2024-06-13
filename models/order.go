package models

import "time"

// Order represents an order in the system.
// @Description Sipariş modelini temsil eder
type Order struct {
	ID         int       `json:"id" example:"1"`               
	UserID     int       `json:"user_id" example:"1"`          
	TotalPrice float64   `json:"total_price" example:"199.99"` 
	CreatedAt  time.Time `json:"created_at"`                   
	Status     string    `json:"status" example:"pending"`     
}

// OrderItem represents an item in an order.
// @Description Sipariş öğesi modelini temsil eder
type OrderItem struct {
	ID        int     `json:"id" example:"1"`         
	OrderID   int     `json:"order_id" example:"1"`   
	ProductID int     `json:"product_id" example:"1"` 
	Quantity  int     `json:"quantity" example:"2"`   
	Price     float64 `json:"price" example:"99.99"`  
}
