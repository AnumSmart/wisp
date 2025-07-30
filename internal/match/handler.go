package match

import (
	"simple_gin_server/configs"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	service ServiceInterface
	config  *configs.Config
}

func NewMatchHandler(service ServiceInterface, config *configs.Config) *MatchHandler {
	return &MatchHandler{
		service: service,
		config:  config,
	}
}

// Хэндлер для вывода имен всех продуктов
func (p *MatchHandler) SearchMatchesHandler(c *gin.Context) {}

// Хэндлер регистрации действия пользователя (лайк/скип/жалоба)
func (p *MatchHandler) RegisterActionHandler(c *gin.Context) {}

// Хэндлер получения списка совпадений, где статус = accepted
func (p *MatchHandler) GetAcceptedMatchesHandler(c *gin.Context) {}

// Хэндлер удаления мовпадения по ID
func (p *MatchHandler) DeleteMetchByIdHandler(c *gin.Context) {}
