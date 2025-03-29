package types

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthStore interface {
	GetRefreshTokenByUserID(ctx context.Context, userID int) (*RefreshToken, error)
	UpsertRefreshToken(ctx context.Context, payload UpdateRefreshTokenPayload) error
}

type CustomClaims struct {
	ID       string `json:"id"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type UserLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Jti       string    `db:"jti"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type UpdateRefreshTokenPayload struct {
	UserID    int       `db:"user_id"`
	Jti       string    `db:"jti"`
	ExpiresAt time.Time `db:"expires_at"`
}

type UpdateRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
