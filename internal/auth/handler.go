package auth

import (
	"fmt"
	"log"
	"net/http"
	"simple_gin_server/configs"
	"simple_gin_server/pkg/jwt_stuff"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	service ServiceInterface
	config  *configs.Config
}

func NewAuthHandler(service ServiceInterface, config *configs.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		config:  config,
	}
}

// Health Check
func (h *AuthHandler) Check(c *gin.Context) {
	c.Writer.Write([]byte("\nHealth check completed successfully"))
}

// Хэндлер валидации и регистрации нового пользователя, доступ к сервисному слою
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	validatedData, exists := c.Get("validatedData")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation data not found"})
		return
	}

	// Приведение типа с проверкой
	user, ok := validatedData.(*RegisterRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid request type",
		})
		return
	}

	//пробуем регистрировать пользователя
	err := h.service.Register(c, user.Email, user.Password)
	if err != nil {
		log.Println("Error during attempt of registration")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// в случае успешной регистрации - выдаём ответ
	c.JSON(http.StatusOK, gin.H{"message": "user registered"})
}

// Хэндлер логина, доступ к сервисному слою
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	//проверяем, есть ли в контексте валидированные данные
	validatedData, exists := c.Get("validatedData")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation data not found"})
		return
	}

	// Приведение типа с проверкой
	user, ok := validatedData.(*LoginRequest)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid request type",
		})
		return
	}

	//log.Printf("email from login request: %v", user.Email)
	//log.Printf("password from login request: %v", user.Password)

	//пробуем залогировать пользователя
	err := h.service.Login(c, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерация JWT токена

	jwtObject := jwt_stuff.NewJWT(
		h.config.Auth.SecretAcc,
		h.config.Auth.SecretRef,
		h.config.Auth.AccessTokenExp,
		h.config.Auth.RefreshTokenExp,
	)

	regUser, err := h.service.GetUserByEmail(c, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении Id зарегестрированного пользователя"})
		return
	}

	//генерируем access и refresh токены
	accessToken, refreshToken, err := jwtObject.GenerateTokens(user.Email, strconv.Itoa(regUser.Id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	err = h.service.AddRefreshTokenToDb(c, user.Email, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка записи refreshToken в БД"})
		return
	}

	// Отправка токена клиенту
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Хэндлер получения email зарегестрированных пользователей, доступ к сервиному слою
func (h *AuthHandler) ListHandler(c *gin.Context) {
	list, err := h.service.GetUserList(c)
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"user_emails": list})
}

// Хэндлер генерации нового access токена, при предоставлении валидного refresh токена
func (h *AuthHandler) ProcessRefreshTokenHandler(c *gin.Context) {
	//Проверка того, что JSON из запроса мапится в нужную структуру refresh токена
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Используем контекст из запроса
	ctx := c.Request.Context()

	// Проверяем не отменён ли контекст
	select {
	case <-ctx.Done():
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "request cancelled"})
		return
	default:
	}

	reqRefToken, err := jwt_stuff.ParseTokenWithClaims(c, req.RefreshToken, h.config.Auth.SecretRef)
	if err != nil {
		log.Println("Wrong refresh token")
		return
	}

	if !reqRefToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		log.Println("Invalid refresh token")
		return
	}

	// Извлекаем claims из токена
	claims, ok := reqRefToken.Claims.(*jwt_stuff.CustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Проверяем, что токен является refresh токеном
	if claims.TokenType != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not a refresh token"})
		return
	}

	// проверка в Reddis (черный список)
	redisKey := fmt.Sprintf("refresh_token:%s", claims.ID) // claims.ID = jti из токена
	exists, err := h.service.ExistsInBlackList(c, redisKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to check token in Redis"})
		return
	}
	if exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"}) // Токен в черном списке!
		return
	}

	// Проверка наличия данного refresh токена в БД
	user, err := h.service.GetUserByClaims(c, *claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user was not found by refresh token claims in BD"})
		return
	}

	if user.RefreshToken == "" || user.RefreshToken != req.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Извлекаем user Email из claims
	email := claims.Email
	userId := claims.ID

	// Создаем новый access токен
	newAccessTokenClaims := jwt_stuff.NewClaims(h.config.Auth.AccessTokenExp, email, userId, "access", "my_app")
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessTokenClaims)

	// Подписываем токен
	accessTokenString, err := newAccessToken.SignedString([]byte(h.config.Auth.SecretAcc))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Возвращаем новый access токен
	c.JSON(http.StatusOK, gin.H{"AccessToken": accessTokenString})
}

// Хэндлер для функции LogOut, инвалидация refresh токена
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	//Проверка того, что JSON из запроса мапится в нужную структуру refresh токена
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Используем контекст из запроса
	ctx := c.Request.Context()

	// Проверяем не отменён ли контекст
	select {
	case <-ctx.Done():
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "request cancelled"})
		return
	default:
	}

	// Инвалидируем токен в хранилище (Redis/БД)
	if err := h.service.InvalidateRefreshToken(c, req.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Refresh token has been added to black list and removed from DB"})
}

// Хэндлер для удаления юзера по его ID, только с админскими прававами(проверка прав через middleware)
func (h *AuthHandler) DeleteUserHandler(c *gin.Context) {

}
