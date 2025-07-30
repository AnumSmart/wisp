package match

import "simple_gin_server/pkg/db"

// Интерфейс для слоя productsSRepository для использования другими источниками
type MatchRepoInterface interface{}

type MatchRepository struct {
	Database db.PgRepoInterface
}

// Конструктор репозитория
func NewMatchRepository(dataBase db.PgRepoInterface) *MatchRepository {
	return &MatchRepository{
		Database: dataBase,
	}
}
