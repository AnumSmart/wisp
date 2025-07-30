package moks

import (
	"context"
	"simple_gin_server/internal/users"

	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) AddUser(ctx context.Context, email, hashed_pass, role string, is_active bool) error {
	args := m.Called(ctx, email, hashed_pass, role, is_active)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*users.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepo) GetEmailLIst(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUserRepo) CheckIfInBaseByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepo) AddRefreshToken(ctx context.Context, email, refreshToken string) error {
	args := m.Called(ctx, email, refreshToken)
	return args.Error(0)
}

func (m *MockUserRepo) ClearRefreshToken(ctx context.Context, claimsEmail string) error {
	args := m.Called(ctx, claimsEmail)
	return args.Error(0)
}
