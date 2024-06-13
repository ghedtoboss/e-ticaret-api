package models

// User represents a user in the system.
// @Description Kullanıcı modelini temsil eder
type User struct {
	ID       int    `json:"id" example:"1"`                   
	Email    string `json:"email" example:"user@example.com"` 
	Password string `json:"password"`                         
	Name     string `json:"name" example:"John Doe"`          
	Role     string `json:"role" example:"seller"`            
}
