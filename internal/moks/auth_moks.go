package moks

import (
	"context"
	"simple_gin_server/internal/users"
	"simple_gin_server/pkg/jwt_stuff"

	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0) // Возвращаем error (может быть nil)
}

func (m *MockAuthService) Login(ctx context.Context, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0) // Возвращаем error (может быть nil)
}

func (m *MockAuthService) GetUserList(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1) // Возвращаем слайс и error
}

func (m *MockAuthService) AddRefreshTokenToDb(ctx context.Context, email, refreshToken string) error {
	args := m.Called(ctx, email, refreshToken)
	return args.Error(0) // Возвращаем error (может быть nil)
}

func (m *MockAuthService) InvalidateRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0) // Возвращаем error (может быть nil)
}

func (m *MockAuthService) ExistsInBlackList(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx)
	return args.Get(0).(bool), args.Error(1) // Возвращаем слайс и error
}

func (m *MockAuthService) GetUserByClaims(ctx context.Context, claims jwt_stuff.CustomClaims) (*users.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(*users.User), args.Error(1) // Возвращаем слайс и error
}

func (m *MockAuthService) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	args := m.Called(ctx)
	return args.Get(0).(*users.User), args.Error(1) // Возвращаем слайс и error
}
