package mocks

import (
	"context"

	"github.com/hoyci/ms-chat/auth-service/types"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) Create(ctx context.Context, user types.CreateUserDatabasePayload) (*types.UserResponse, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*types.UserResponse), args.Error(1)
}

func (m *MockUserStore) GetByID(ctx context.Context, userID string) (*types.UserResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*types.UserResponse), args.Error(1)
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*types.GetByEmailResponse, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*types.GetByEmailResponse), args.Error(1)
}

func (m *MockUserStore) UpdateByID(
	ctx context.Context, userID string, newUser types.UpdateUserPayload,
) (*types.UserResponse, error) {
	args := m.Called(ctx, userID, newUser)
	return args.Get(0).(*types.UserResponse), args.Error(1)
}

func (m *MockUserStore) DeleteByID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
