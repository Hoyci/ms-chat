package types

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
