package auth

import (
	"net/http"
	"net/http/httptest"
	"simple_gin_server/configs"
	"simple_gin_server/internal/moks"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// функция настройки тестового окружения
func setUpAuthHandlerTest(t *testing.T) (*moks.MockAuthService, *AuthHandler) {
	mockService := new(moks.MockAuthService)
	conf := configs.LoadConfig()
	return mockService, NewAuthHandler(mockService, conf)
}

// тест проверяет "счастливый путь" (happy path) регистрации пользователя
func TestRegisterHandler_ValidCredentials(t *testing.T) {
	// Настройка теста
	mockService, handler := setUpAuthHandlerTest(t)
	gin.SetMode(gin.TestMode)

	// Ожидаемый вызов
	mockService.On("Register", mock.Anything, "test@example.com", "va123*tro").Return(nil)

	// Создание тестового контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Устанавливаем validatedData в контекст
	c.Set("validatedData", &RegisterRequest{
		Email:    "test@example.com",
		Password: "va123*tro",
	})

	// Вызов хендлера
	handler.RegisterHandler(c)

	// Проверки
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "user registered"}`, w.Body.String())
	mockService.AssertExpectations(t)
}

// тест проверяет разные случаи, когда валидация не будет пройдена
func TestRegisterHandler_ValidationErrors(t *testing.T) {
	// создаём слайс структур тестов
	tests := []struct {
		name        string          // имя теста (описание)
		input       RegisterRequest // входные данные для валидатора
		wantStatus  int             // желаемый статус
		wantMessage string          // желаемое сообщение
	}{
		{
			name: "Empty email",
			input: RegisterRequest{
				Email:    "",
				Password: "validPass123!",
			},
			wantStatus:  http.StatusBadRequest,
			wantMessage: `{"error":"email is required"}`,
		},
		{
			name: "Invalid email format",
			input: RegisterRequest{
				Email:    "not-an-email",
				Password: "validPass123!",
			},
			wantStatus:  http.StatusBadRequest,
			wantMessage: `{"error":"invalid email format"}`,
		},
		{
			name: "Weak password",
			input: RegisterRequest{
				Email:    "test@example.com",
				Password: "123",
			},
			wantStatus:  http.StatusBadRequest,
			wantMessage: `{"error":"password must be at least 8 characters"}`,
		},
		{
			name: "Password without special chars",
			input: RegisterRequest{
				Email:    "test@example.com",
				Password: "Password123",
			},
			wantStatus:  http.StatusBadRequest,
			wantMessage: `{"error":"password must contain special character"}`,
		},
	}

	for _, tt := range tests {
		// t.Run запустит тест с имененм в отдельной гоурутине, надо проверить, нужна ли синхронизация через waitgroup
		t.Run(tt.name, func(t *testing.T) {
			mockService, handler := setUpAuthHandlerTest(t)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("validatedData", &tt.input)

			handler.RegisterHandler(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantMessage, w.Body.String())
			mockService.AssertNotCalled(t, "Register")
		})
	}
}
