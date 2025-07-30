package auth

import (
	"context"
	"simple_gin_server/internal/moks"
	"simple_gin_server/internal/users"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// функция настройки тестового окружения
func setUpServiceTest(t *testing.T) (*moks.MockUserRepo, *moks.MockRedisRepo, *AuthService) {
	mockUserRepo := new(moks.MockUserRepo)
	mockReddisRepo := new(moks.MockRedisRepo)
	return mockUserRepo, mockReddisRepo, NewAuthService(mockUserRepo, mockReddisRepo)
}

// тест для метода Register у слоя Service
func TestAuthService_Register(t *testing.T) {

	t.Run("Successfull registration", func(t *testing.T) {
		mockUserRepo, _, service := setUpServiceTest(t)

		mockUserRepo.On("CheckIfInBaseByEmail", mock.Anything, "test@example.com").Return(false, nil)
		mockUserRepo.On("AddUser", mock.Anything, "test@example.com", mock.Anything).Return(nil)

		err := service.Register(context.Background(), "test@example.com", "password123")

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockUserRepo, _, service := setUpServiceTest(t)

		mockUserRepo.On("CheckIfInBaseByEmail", mock.Anything, "existing@example.com").Return(true, nil)

		err := service.Register(context.Background(), "existing@example.com", "password123")

		assert.Error(t, err)
		assert.Equal(t, "user with such Email is in base", err.Error())
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("context canceled", func(t *testing.T) {
		_, _, service := setUpServiceTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := service.Register(ctx, "test@example.com", "password123")

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

// тест для метода Login у слоя Service
func TestAuthService_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mockUserRepo, _, service := setUpServiceTest(t)

		// Хешированный пароль для теста
		hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMy.Mrq1V8H3M3kL6h7pW1pJ5Qn6T7XzB1O"

		mockUser := &users.User{
			Email:    "test@example.com",
			HashPass: hashedPassword,
		}

		mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)

		err := service.Login(context.Background(), "test@example.com", "password123")

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockUserRepo, _, service := setUpServiceTest(t)

		// Хешированный пароль для теста
		hashedPassword := "$2a$10$N9qo8uLOickgx2ZMRZoMy.Mrq1V8H3M3kL6h7pW1pJ5Qn6T7XzB1O"

		mockUser := &users.User{
			Email:    "test@example.com",
			HashPass: hashedPassword,
		}

		mockUserRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)

		err := service.Login(context.Background(), "test@example.com", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, users.ErrWrongCredentials, err.Error())
		mockUserRepo.AssertExpectations(t)
	})
}
