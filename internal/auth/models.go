package auth

type HealthResponse struct {
	Status string `json:"status"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginResponce struct {
	Token string `json:"token"`
}

type RegisterResponce struct {
	Token string `json:"token"`
}

// Структура для входящего запроса
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
