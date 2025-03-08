package types

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	ID       string `json:"id"`
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
