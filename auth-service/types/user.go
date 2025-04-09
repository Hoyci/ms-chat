package types

import (
	"context"
	"time"
)

type UserStore interface {
	Create(ctx context.Context, user CreateUserDatabasePayload) (*UserResponse, error)
	GetByID(ctx context.Context, userID string) (*UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*GetByEmailResponse, error)
	UpdateByID(ctx context.Context, userID string, user UpdateUserPayload) (*UserResponse, error)
	DeleteByID(ctx context.Context, userID string) error
}

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"passwordHash"`
	CreatedAt    time.Time  `json:"createdAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type GetByEmailResponse struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"passwordHash"`
	CreatedAt    time.Time  `json:"createdAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type UserResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type CreateUserRequestPayload struct {
	Username        string `json:"username" validate:"required,min=5"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

type CreateUserResponse struct {
	Message string `json:"message"`
}

type CreateUserDatabasePayload struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwor_hash"`
}

type UpdateUserPayload struct {
	Username string `json:"username" validate:"required,min=5"`
	Email    string `json:"email" validate:"required,email"`
}

type DeleteUserByIDResponse struct {
	ID string `json:"id"`
}
