package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"userID"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
