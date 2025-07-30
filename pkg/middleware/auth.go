package middleware

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"simple_gin_server/configs"
	"simple_gin_server/pkg/jwt_stuff"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// ValidateMiddleware создает middleware для валидации
func ValidateAuthMiddleware(model interface{}) gin.HandlerFunc {
	validate := validator.New()

	return func(c *gin.Context) {
		// Создаем новый экземпляр структуры для валидации
		request := reflect.New(reflect.TypeOf(model).Elem()).Interface()

		// Парсим JSON
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})

			c.Abort()
			return
		}

		// Валидируем структуру
		if err := validate.Struct(request); err != nil {
			errors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				errors[err.Field()] = err.Tag()
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})

			c.Abort()
			return
		}
		// Сохраняем валидированные данные в контекст для использования в обработчике
		c.Set("validatedData", request)
		c.Next()
	}
}

// ---------------------------------------------------ПОКА В РАЗРАБОТКЕ-----------------------------------------------------
func AuthMiddleware(config *configs.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка
		authHeader := c.GetHeader("Authorization")

		// проверяем наличие токена в заготовке, если его нет, выдаём ошибку и не пускаем запрос дальше
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Проверяем формат "Bearer <token>"
		tokenString, err := CheckBearerFormat(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		//Парсим токен
		token, err := jwt_stuff.ParseTokenWithClaims(c, tokenString, config.Auth.SecretAcc)
		if err != nil {
			log.Println("Invalid token")
			return
		}

		// Проверяем claims
		if claims, ok := token.Claims.(*jwt_stuff.CustomClaims); ok && token.Valid {
			// Добавляем данные пользователя в контекст
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role) // Важно для RoleMiddleware
			c.Set("is_active", claims.IsActive)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		}
	}
}

// ---------------------------------------------------ПОКА В РАЗРАБОТКЕ-----------------------------------------------------

func CheckBearerFormat(authHeader string) (string, error) {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}
	return "", errors.New("Invalid authorization header format")
}
