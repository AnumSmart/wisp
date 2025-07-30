package profile

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	service ServiceInterface
}

func NewProfileHandler(service ServiceInterface) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

// Хэндлер для создания нового профиля
func (p *ProfileHandler) CreateNewProfileHandler(c *gin.Context) {
	var profile Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// пробуем извлеч Email из контектста
	email, exists := c.Get("user_email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}

	// Делаем приведение типа к string
	emailStr, ok := email.(string) // Приведение типа
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong Email type"})
		return
	}

	// Создаём новый профиль
	err := p.service.CreateNewProfile(c, &profile, emailStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new profile"})
		return
	}

	messageOk := fmt.Sprintf("New profile with nickName:%s has been created\n", profile.NickName)

	c.JSON(http.StatusOK, gin.H{"Message": messageOk})
}

// Хэндлер для получения своего профиля
func (p *ProfileHandler) GetMyProfileHandler(c *gin.Context) {}

// Хэндлер для оновления своего профиля
func (p *ProfileHandler) UpdateMyProfileHandler(c *gin.Context) {}

// Хэндлер для удаления своего профиля
func (p *ProfileHandler) DeleteMyProfileHandler(c *gin.Context) {}
