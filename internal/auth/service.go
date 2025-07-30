package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"simple_gin_server/internal/users"
	"simple_gin_server/pkg/db"
	"simple_gin_server/pkg/jwt_stuff"

	"time"

	"golang.org/x/crypto/bcrypt"
)

// Интерфейс для слоя authService для использования другими источниками
type ServiceInterface interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) error
	GetUserList(ctx context.Context) ([]string, error)
	AddRefreshTokenToDb(ctx context.Context, email, refreshToken string) error
	InvalidateRefreshToken(ctx context.Context, refreshToken string) error
	ExistsInBlackList(ctx context.Context, key string) (bool, error)
	GetUserByClaims(ctx context.Context, claims jwt_stuff.CustomClaims) (*users.User, error)
	GetUserByEmail(ctx context.Context, email string) (*users.User, error)
}

type AuthService struct {
	repo      users.UserRepoInterface
	redisRepo db.ReddisRepoInterface
}

// Конструктор слоя сервис
func NewAuthService(repo users.UserRepoInterface, redisRepo db.ReddisRepoInterface) *AuthService {
	return &AuthService{
		repo:      repo,
		redisRepo: redisRepo,
	}
}

// Добавление нового пользователя в базу с хэшированным паролем
func (s *AuthService) Register(ctx context.Context, email, password string) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	isInBase, _ := s.repo.CheckIfInBaseByEmail(ctx, email)

	if isInBase {
		return errors.New("user with such Email is in base")
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка при хешировании пароля")
	}

	err = s.repo.AddUser(ctx, email, string(hashedPassword), "user", true)
	if err != nil {
		return errors.New("failed to add new user to the DB")
	}
	return nil
}

// Логи юзера по email и pasword, при успешном логировании - в ответе будет access и refresh jwt токены
func (s *AuthService) Login(ctx context.Context, email, password string) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	// Проверяем существует ли пользователь с данным email уже в базе
	existedUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existedUser == nil {
		log.Printf("error during search in the DB: %v", existedUser)
		return errors.New(users.ErrWrongCredentials)
	}

	//сравниваем хэши паролей, тот, что в базе и тот, что логинится
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.HashPass), []byte(password))
	if err != nil {
		return errors.New(users.ErrWrongCredentials)
	}
	return nil
}

// Получаем список email зарегестрированных пользователей
func (s *AuthService) GetUserList(ctx context.Context) ([]string, error) {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	res, err := s.repo.GetEmailLIst(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Добавляем refresh токен в БД
func (s *AuthService) AddRefreshTokenToDb(ctx context.Context, email, refreshToken string) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	err := s.repo.AddRefreshToken(ctx, email, refreshToken)
	if err != nil {
		return err
	}
	return nil
}

// Инвалидируем refresh токен в БД
// ---------------------------------------????????????????????????????????????????????????
func (s *AuthService) InvalidateRefreshToken(ctx context.Context, refreshToken string) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	// Получаем jti из токена
	claims, err := jwt_stuff.ParseTokenWithoutVerification(refreshToken)
	log.Printf("[service.go]---[InvalidateRefreshToken()]---[ParseTokenWithoutVerification()]---err: %v", err)
	if err != nil {
		return err
	}

	// Дополнительная проверка что это именно refresh token
	if claims.TokenType != "refresh" {
		return errors.New("not a refresh token")
	}

	// Вычисляем оставшееся время жизни токена
	ttl := time.Until(claims.ExpiresAt.Time) // Верный способ для jwt.NumericDate
	log.Printf("[service.go]---[InvalidateRefreshToken()]---TTL(refresh token): %v", ttl)

	// Сохраняем в Redis
	key := fmt.Sprintf("refresh_token:%s", claims.ID)
	if err := s.redisRepo.Set(ctx, key, "invalid", ttl); err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}

	// 2. Очистка в PostgreSQL
	if err := s.repo.ClearRefreshToken(ctx, claims.Email); err != nil {
		// Важно: даже если очистка в БД не удалась, токен уже инвалидирован в Redis
		log.Printf("warning: failed to clear refresh token in DB for user %s: %v", claims.Email, err)
	}
	return nil
}

// проверка токена на присутствие в черном списке
func (s *AuthService) ExistsInBlackList(ctx context.Context, key string) (bool, error) {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return false, err
	}

	return s.redisRepo.Exists(ctx, key)
}

// Достать юзера используя calims от refresh токена
func (s *AuthService) GetUserByClaims(ctx context.Context, claims jwt_stuff.CustomClaims) (*users.User, error) {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	email := claims.Email

	return s.repo.FindByEmail(ctx, email)
}

// Достать userId, используя Email
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Ищем юзера по Email
	existedUser, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return existedUser, nil
}
