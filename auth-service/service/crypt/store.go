package crypt

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHandler interface {
	HashPassword(ctx context.Context, password string) (string, error)
	CheckPassword(ctx context.Context, hashedPassword, password string) error
}

type BcryptPasswordStore struct{}

func (b *BcryptPasswordStore) HashPassword(ctx context.Context, password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (b *BcryptPasswordStore) CheckPassword(ctx context.Context, hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
