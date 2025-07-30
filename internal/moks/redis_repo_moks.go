package moks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRedisRepo struct {
	mock.Mock
}

func (m *MockRedisRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisRepo) Exists(ctx context.Context, redisKey string) (bool, error) {
	args := m.Called(ctx, redisKey)
	return args.Get(0).(bool), args.Error(1)
}
