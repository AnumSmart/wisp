package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Тестовая структура для валидации
type TestAuthRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

// тест middleware валидации входящих данных
func TestValidateAuthMiddleware(t *testing.T) {
	// Создаем тестовый роутер Gin
	router := gin.Default()
	middleware := ValidateAuthMiddleware(TestAuthRequest{})

	// Тестовый обработчик
	router.POST("/test", middleware, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	tests := []struct {
		name         string
		payload      string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid request",
			payload:      `{"email":"test@example.com","password":"validPass123!"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid email",
			payload:      `{"email":"invalid-email","password":"validPass123!"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Validation failed","details":{"Email":"email"}}`,
		},
		{
			name:         "Missing password",
			payload:      `{"email":"test@example.com"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Validation failed","details":{"Password":"required"}}`,
		},
		{
			name:         "Invalid JSON",
			payload:      `invalid-json`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Invalid request body"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			// Выполняем запрос
			router.ServeHTTP(w, req)

			// Проверяем результаты
			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
