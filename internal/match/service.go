package match

// Интерфейс для слоя productsService для использования другими источниками
type ServiceInterface interface{}

type MatchService struct {
	repo MatchRepoInterface
}

// Конструктор слоя сервис
func NewMatchService(repo MatchRepoInterface) *MatchService {
	return &MatchService{
		repo: repo,
	}
}
