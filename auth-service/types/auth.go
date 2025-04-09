package types

import (
	"context"
	"time"
)

type AuthStore interface {
	GetRefreshTokenByUserID(ctx context.Context, userID string) (*RefreshToken, error)
	UpsertRefreshToken(ctx context.Context, payload UpdateRefreshTokenPayload) error
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
	UserID    string    `db:"user_id"`
	Jti       string    `db:"jti"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type UpdateRefreshTokenPayload struct {
	UserID    string    `db:"user_id"`
	Jti       string    `db:"jti"`
	ExpiresAt time.Time `db:"expires_at"`
}

type UpdateRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
