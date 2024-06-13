package models

import "github.com/dgrijalva/jwt-go"

// Claims represents the JWT claims.
// @Description JWT iddialarını temsil eder
type Claims struct {
	Username string `json:"username" example:"user@example.com"` 
	UserID   int    `json:"userID" example:"1"`                  
	Role     string `json:"role" example:"seller"`               
	jwt.StandardClaims
}
