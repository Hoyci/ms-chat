package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockPasswordHandler struct {
	mock.Mock
}

func (m *MockPasswordHandler) HashPassword(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHandler) CheckPassword(ctx context.Context, hashedPassword, password string) error {
	args := m.Called(ctx, hashedPassword, password)
	return args.Error(0)
}
